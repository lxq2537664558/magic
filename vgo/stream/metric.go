package stream

import "github.com/influxdata/influxdb/client"

// TransferData transfer data(inpute transfer data)
type Metric struct {
	pt *client.Point
}
