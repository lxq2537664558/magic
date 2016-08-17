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

	st := time.Now()
	for name, v := range rawM {
		ms := strings.Split(name, ".")
		if len(ms) != 3 {
			ctx.WriteString("invalid metric name, you should pass name like this: app.func.timer.error name is:" + name)
			return
		}

		mType := ms[2]
		fields := v.(map[string]interface{})

		switch mType {
		case "gauge": //ok
			h.acc.AddFields(name, fields, agent.Conf.Tags, time.Now())

		case "count": //ok
			fields["value"] = fields["count"]
			delete(fields, "count")
			h.acc.AddFields(name, fields, agent.Conf.Tags, time.Now())

		case "timer": //ok
			fields["p50"] = fields["median"]
			delete(fields, "median")

			fields["p75"] = fields["75%"]
			delete(fields, "75%")

			fields["p95"] = fields["95%"]
			delete(fields, "95%")

			fields["p99"] = fields["99%"]
			delete(fields, "99%")

			fields["p999"] = fields["99.9%"]
			delete(fields, "99.9%")

			fields["m1"] = fields["1m.rate"]
			delete(fields, "1m.rate")

			fields["m5"] = fields["5m.rate"]
			delete(fields, "5m.rate")

			fields["m15"] = fields["15m.rate"]
			delete(fields, "15m.rate")

			fields["meanrate"] = fields["mean.rate"]
			delete(fields, "mean.rate")
			h.acc.AddFields(name, fields, agent.Conf.Tags, time.Now())

		case "meter": //ok
			fields["m1"] = fields["1m.rate"]
			delete(fields, "1m.rate")

			fields["m5"] = fields["5m.rate"]
			delete(fields, "5m.rate")

			fields["m15"] = fields["15m.rate"]
			delete(fields, "15m.rate")

			fields["mean"] = fields["mean.rate"]
			delete(fields, "mean.rate")

			h.acc.AddFields(name, fields, agent.Conf.Tags, time.Now())

		case "histogram":
			fields["p50"] = fields["median"]
			delete(fields, "median")

			fields["p75"] = fields["75%"]
			delete(fields, "75%")

			fields["p95"] = fields["95%"]
			delete(fields, "95%")

			fields["p99"] = fields["99%"]
			delete(fields, "99%")

			fields["p999"] = fields["99.9%"]
			delete(fields, "99.9%")

			h.acc.AddFields(name, fields, agent.Conf.Tags, time.Now())
		case "counter":
			ctx.WriteString("maybe you should pass app.func.count instead! error type: counter")
			return
		}
	}

	tu := time.Now().Sub(st)
	log.Println("http time used: ", tu.Nanoseconds(), "ns")
}

func init() {
	agent.AddInput("java_metrics_http", &HttpListener{})
}
