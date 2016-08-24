package stream

import "github.com/influxdata/influxdb/client/v2"

// TransferData transfer data(inpute transfer data)
type Metric struct {
	MetricData *client.Point
}
