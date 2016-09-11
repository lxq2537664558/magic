package system

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/net"

	"github.com/corego/vgo/mecury/agent"
)

type IfStat struct {
	ps PS

	LastIfStats     map[string]net.IOCountersStat
	LastCollectTime time.Time

	skipChecks bool
	Interfaces []string
}

func (_ *IfStat) Description() string {
	return "Read metrics about network interface usage"
}

var ifstatSampleConfig = `
  ## By default, telegraf gathers stats from any up interface (excluding loopback)
  ## Setting interfaces will tell it to gather these explicit interfaces,
  ## regardless of status.
  ##
  # interfaces = ["eth0"]
`

func (_ *IfStat) SampleConfig() string {
	return netSampleConfig
}

func (s *IfStat) Gather(acc agent.Accumulator) error {
	netio, err := s.ps.NetIO()
	if err != nil {
		return fmt.Errorf("error getting netif info: %s", err)
	}

	now := time.Now()

	// first time,just record the stats
	if s.LastIfStats == nil {
		s.LastIfStats = make(map[string]net.IOCountersStat)
		for _, v := range netio {
			s.LastIfStats[v.Name] = v
		}
		s.LastCollectTime = now
	}

	for _, io := range netio {
		if len(s.Interfaces) != 0 {
			var found bool

			for _, name := range s.Interfaces {
				if name == io.Name {
					found = true
					break
				}
			}

			if !found {
				continue
			}
		}

		tags := map[string]string{
			"interface": io.Name,
		}

		lio := s.LastIfStats[io.Name]
		duration := now.Sub(s.LastCollectTime).Seconds()

		fields := map[string]interface{}{
			"out_bytes": float64(io.BytesSent-lio.BytesSent) / duration,
			"in_bytes":  float64(io.BytesRecv-lio.BytesRecv) / duration,

			"out_packets": float64(io.PacketsSent-lio.PacketsSent) / duration,
			"in_packets":  float64(io.PacketsRecv-lio.PacketsRecv) / duration,

			"out_errors": float64(io.Errout-lio.Errout) / duration,
			"in_errors":  float64(io.Errin-lio.Errin) / duration,

			"out_dropped": float64(io.Dropout-lio.Dropout) / duration,
			"in_dropped":  float64(io.Dropin-lio.Dropin) / duration,
		}
		acc.AddFields("ifstat", fields, tags)

	}

	for _, v := range netio {
		s.LastIfStats[v.Name] = v
	}
	s.LastCollectTime = now

	return nil
}

func init() {
	agent.AddInput("ifstat", &IfStat{ps: &systemPS{}})
}
