package system

import (
	"fmt"

	"github.com/corego/vgo/mecury/agent"
)

type MemStats struct {
	ps PS
}

func (_ *MemStats) Description() string {
	return "Read metrics about memory usage"
}

func (_ *MemStats) SampleConfig() string { return "" }

func (s *MemStats) Gather(acc agent.Accumulator) error {
	vm, err := s.ps.VMStat()
	if err != nil {
		return fmt.Errorf("error getting virtual memory info: %s", err)
	}

	fields := map[string]interface{}{
		"total":             vm.Total,
		"available":         vm.Available,
		"used":              vm.Used,
		"free":              vm.Free,
		"cached":            vm.Cached,
		"buffered":          vm.Buffers,
		"active":            vm.Active,
		"inactive":          vm.Inactive,
		"used_percent":      100 * float64(vm.Used) / float64(vm.Total),
		"available_percent": 100 * float64(vm.Available) / float64(vm.Total),
	}
	acc.AddFields("mem", fields, nil)

	return nil
}

func init() {
	agent.AddInput("mem", &MemStats{ps: &systemPS{}})
}
