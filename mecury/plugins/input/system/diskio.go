package system

import (
	"fmt"

	"github.com/aiyun/openapm/mecury/agent"
)

type DiskIOStats struct {
	ps PS

	Devices          []string
	SkipSerialNumber bool
}

func (_ *DiskIOStats) Description() string {
	return "Read metrics about disk IO by device"
}

var diskIoSampleConfig = `
  ## By default, telegraf will gather stats for all devices including
  ## disk partitions.
  ## Setting devices will restrict the stats to the specified devices.
  # devices = ["sda", "sdb"]
  ## Uncomment the following line if you need disk serial numbers.
  # skip_serial_number = false
`

func (_ *DiskIOStats) SampleConfig() string {
	return diskIoSampleConfig
}

func (s *DiskIOStats) Gather(acc agent.Accumulator) error {
	diskio, err := s.ps.DiskIO()
	if err != nil {
		return fmt.Errorf("error getting disk io info: %s", err)
	}

	var restrictDevices bool
	devices := make(map[string]bool)
	if len(s.Devices) != 0 {
		restrictDevices = true
		for _, dev := range s.Devices {
			devices[dev] = true
		}
	}

	for _, io := range diskio {
		_, member := devices[io.Name]
		if restrictDevices && !member {
			continue
		}
		tags := map[string]string{}
		tags["name"] = io.Name
		if !s.SkipSerialNumber {
			if len(io.SerialNumber) != 0 {
				tags["serial"] = io.SerialNumber
			} else {
				tags["serial"] = "unknown"
			}
		}

		fields := map[string]interface{}{
			"reads":           io.ReadCount,
			"writes":          io.WriteCount,
			"read_bytes":      io.ReadBytes,
			"write_bytes":     io.WriteBytes,
			"read_time":       io.ReadTime,
			"write_time":      io.WriteTime,
			"io_time":         io.IoTime,
			"iop_in_progress": io.IopsInProgress,
		}
		acc.AddFields("diskio", fields, tags)
	}

	return nil
}

func init() {
	agent.AddInput("diskio", &DiskIOStats{ps: &systemPS{}, SkipSerialNumber: true})
}
