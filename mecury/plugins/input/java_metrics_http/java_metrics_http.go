package java_metrics_http

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/corego/vgo/mecury/agent"
	"github.com/valyala/fasthttp"
)

type HttpListener struct {
	ServiceAddress string

	// Keep the accumulator in this struct
	acc agent.Accumulator

	listener *fasthttp.Server
}

func (u *HttpListener) SampleConfig() string {
	return ""
}

func (u *HttpListener) Description() string {
	return "Java metrics http listener"
}

func (u *HttpListener) Gather(_ agent.Accumulator) error {
	return nil
}

func (h *HttpListener) Start(acc agent.Accumulator) error {
	s := &fasthttp.Server{
		Handler: h.handle,
	}

	go func() {
		err := s.ListenAndServe(h.ServiceAddress)
		if err != nil {
			log.Fatalln("start java_metrics_http failed: ", err)
		}
	}()

	h.acc = acc

	log.Println("java metrics http  listening on: ", h.ServiceAddress)
	return nil
}

func (h *HttpListener) Stop() {

}

func (h *HttpListener) handle(ctx *fasthttp.RequestCtx) {
	m := ctx.FormValue("metrics")
	var rawM map[string]interface{}
	json.Unmarshal(m, &rawM)

	for n, v := range rawM {
		fields := v.(map[string]interface{})
		tags := make(map[string]string)
		for k, v := range agent.Conf.Tags {
			tags[k] = v
		}

		ns := strings.Split(n, ",")
		if len(ns) <= 0 {
			log.Println("invalid metric name: ", n)
			continue
		} else if len(ns) == 1 {
			h.acc.AddFields(ns[0], fields, tags, time.Now())
			continue
		}

		for _, v := range ns[1:] {
			tag := strings.Split(v, "=")
			if len(tag) == 2 {
				tags[tag[0]] = tag[1]
			}
		}

		h.acc.AddFields(ns[0], fields, tags, time.Now())
	}
}

func init() {
	agent.AddInput("java_metrics_http", &HttpListener{})
}
