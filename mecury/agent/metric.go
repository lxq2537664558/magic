package agent

import (
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

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

	// Point returns a influxdb client.Point object
	Point() *client.Point
}
