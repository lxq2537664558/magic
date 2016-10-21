package system

import (
	"fmt"
	"time"

	"github.com/aiyun/openapm/mecury/agent"
	"github.com/shirou/gopsutil/disk"
)

type IOStats struct {
	ps PS

	LastDiskStats    map[string]disk.IOCountersStat
	LastCollectTime  time.Time
	Devices          []string
	SkipSerialNumber bool
}

func (_ *IOStats) Description() string {
	return "Read metrics about disk IO by device"
}

var ioStatsSampleConfig = `
  ## By default, telegraf will gather stats for all devices including
  ## disk partitions.
  ## Setting devices will restrict the stats to the specified devices.
  # devices = ["sda", "sdb"]
  ## Uncomment the following line if you need disk serial numbers.
  # skip_serial_number = false
`

func (_ *IOStats) SampleConfig() string {
	return diskIoSampleConfig
}

func (s *IOStats) Gather(acc agent.Accumulator) error {
	diskio, err := s.ps.DiskIO()
	if err != nil {
		return fmt.Errorf("error getting disk io info: %s", err)
	}

	now := time.Now()
	// first time,just record the stats
	if s.LastDiskStats == nil {
		s.LastDiskStats = diskio
		s.LastCollectTime = now
		return nil
	}

	var restrictDevices bool
	devices := make(map[string]bool)
	if len(s.Devices) != 0 {
		restrictDevices = true
		for _, dev := range s.Devices {
			devices[dev] = true
		}
	}

	for k, io := range diskio {
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

		lio := s.LastDiskStats[k]

		rio := io.ReadCount - lio.ReadCount
		wio := io.WriteCount - lio.WriteCount
		nio := rio + wio

		deltaRbytes := io.ReadBytes - lio.ReadBytes
		deltaWbytes := io.WriteBytes - lio.WriteBytes

		ruse := io.ReadTime - lio.ReadTime
		wuse := io.WriteTime - lio.WriteTime
		use := io.IoTime - lio.IoTime

		await := 0.0
		svctm := 0.0

		if nio != 0 {
			await = float64(ruse+wuse) / float64(nio)
			svctm = float64(use) / float64(nio)
		}

		duration := now.Sub(s.LastCollectTime).Nanoseconds() / 1000000

		ioutil := float64(use) * 100.0 / float64(duration)

		fields := map[string]interface{}{
			"read_bytes":  deltaRbytes,
			"write_bytes": deltaWbytes,
			"await":       await,
			"svctm":       svctm,
			"ioutil":      ioutil,
			"msec_read":   float64(ruse) / float64(now.Sub(s.LastCollectTime).Seconds()), // ms / s
			"msec_write":  float64(wuse) / float64(now.Sub(s.LastCollectTime).Seconds()),
		}

		acc.AddFields("iostats", fields, tags)
	}

	s.LastDiskStats = diskio
	s.LastCollectTime = now

	return nil
}

func init() {
	agent.AddInput("iostats", &IOStats{ps: &systemPS{}, SkipSerialNumber: true})
}
