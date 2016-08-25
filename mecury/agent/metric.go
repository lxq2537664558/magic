package agent

import (
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

// Metric ...
type Metric interface {
	// Name returns the measurement name of the metric
	Name() string

	// Name returns the tags associated with the metric
	Tags() map[string]string

	// Time return the timestamp for the metric
	Time() time.Time

	// UnixNano returns the unix nano time of the metric
	UnixNano() int64

	// Fields returns the fields for the metric
	Fields() map[string]interface{}

	// String returns a line-protocol string of the metric
	String() string

	// PrecisionString returns a line-protocol string of the metric, at precision
	PrecisionString(string) string

	// Point returns a influxdb client.Point object1
	Point() *client.Point
}

// metric is a wrapper of the influxdb client.Point struct
type MetricWrapper struct {
	Pt *client.Point `json:"pt"`
}

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

// NewMetric returns a metric with the given timestamp. If a timestamp is not
// given, then data is sent to the database without a timestamp, in which case
// the server will assign local time upon reception. NOTE: it is recommended to
// send data with a timestamp.
func NewMetric(
	name string,
	tags map[string]string,
	fields map[string]interface{},
	t time.Time,
) (Metric, error) {
	pt, err := client.NewPoint(name, tags, fields, t)
	if err != nil {
		return nil, err
	}
	return &MetricWrapper{
		Pt: pt,
	}, nil
}

func (m *MetricWrapper) Name() string {
	return m.Pt.Name()
}

func (m *MetricWrapper) Tags() map[string]string {
	return m.Pt.Tags()
}

func (m *MetricWrapper) Time() time.Time {
	return m.Pt.Time()
}

func (m *MetricWrapper) UnixNano() int64 {
	return m.Pt.UnixNano()
}

func (m *MetricWrapper) Fields() map[string]interface{} {
	return m.Pt.Fields()
}

func (m *MetricWrapper) String() string {
	return m.Pt.String()
}

func (m *MetricWrapper) PrecisionString(precison string) string {
	return m.Pt.PrecisionString(precison)
}

func (m *MetricWrapper) Point() *client.Point {
	return m.Pt
}
