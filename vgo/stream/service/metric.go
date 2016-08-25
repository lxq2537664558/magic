package service



import "time"

// MetricData transfer data(inpute transfer data)
//easyjson:json
type Metrics struct {
	Data []*MetricData `json:"d"`
}

type MetricData struct {
	Name   string                 `json:"n"`
	Tags   map[string]string      `json:"ts"`
	Fields map[string]interface{} `json:"f"`
	Time   time.Time              `json:"t"`
}
