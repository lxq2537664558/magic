package java_metrics_http

import (
	"encoding/json"
	"log"
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

	st := time.Now()
	for name, v := range rawM {
		fields := v.(map[string]interface{})
		h.acc.AddFields(name, fields, agent.Conf.Tags, time.Now())
	}

	tu := time.Now().Sub(st)
	log.Println("http time used: ", tu.Nanoseconds(), "ns")
}

func init() {
	agent.AddInput("java_metrics_http", &HttpListener{})
}
