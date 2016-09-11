// +build linux

package system

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/corego/vgo/mecury/agent"
)

// /proc/stat file line prefixes to gather stats on:
var (
	interrupts       = []byte("intr")
	context_switches = []byte("ctxt")
	processes_forked = []byte("processes")
	disk_pages       = []byte("page")
	boot_time        = []byte("btime")
)

type Kernel struct {
	LastCTXT        int64
	LastIntr        int64
	LastPagein      int64
	LastPageout     int64
	LastCollectTime time.Time
	statFile        string
}

func (k *Kernel) Description() string {
	return "Get kernel statistics from /proc/stat"
}

func (k *Kernel) SampleConfig() string { return "" }

func (k *Kernel) Gather(acc agent.Accumulator) error {
	data, err := k.getProcStat()
	if err != nil {
		return err
	}

	fields := make(map[string]interface{})

	now := time.Now()

	dataFields := bytes.Fields(data)
	for i, field := range dataFields {
		switch {
		case bytes.Equal(field, interrupts):
			m, err := strconv.ParseInt(string(dataFields[i+1]), 10, 64)
			if err != nil {
				return err
			}

			if k.LastIntr != 0 {
				fields["interrupts"] = float64(m-k.LastIntr) / now.Sub(k.LastCollectTime).Seconds()
			}
			k.LastIntr = m

		case bytes.Equal(field, context_switches):
			m, err := strconv.ParseInt(string(dataFields[i+1]), 10, 64)
			if err != nil {
				return err
			}

			if k.LastCTXT != 0 {
				fields["context_switches"] = float64(m-k.LastCTXT) / now.Sub(k.LastCollectTime).Seconds()
			}
			k.LastCTXT = m

		case bytes.Equal(field, processes_forked):
			m, err := strconv.ParseInt(string(dataFields[i+1]), 10, 64)
			if err != nil {
				return err
			}
			fields["processes_forked"] = int64(m)
		case bytes.Equal(field, boot_time):
			m, err := strconv.ParseInt(string(dataFields[i+1]), 10, 64)
			if err != nil {
				return err
			}
			fields["boot_time"] = int64(m)
		case bytes.Equal(field, disk_pages):
			in, err := strconv.ParseInt(string(dataFields[i+1]), 10, 64)
			if err != nil {
				return err
			}
			out, err := strconv.ParseInt(string(dataFields[i+2]), 10, 64)
			if err != nil {
				return err
			}

			if k.LastPagein == 0 {
				fields["disk_pages_in"] = float64(in-k.LastPagein) / now.Sub(k.LastCollectTime).Seconds()
			}
			k.LastPagein = in

			if k.LastPageout == 0 {
				fields["disk_pages_out"] = float64(out-k.LastPageout) / now.Sub(k.LastCollectTime).Seconds()
			}
			k.LastPageout = out
		}
	}

	k.LastCollectTime = now
	acc.AddFields("kernel", fields, map[string]string{})

	return nil
}

func (k *Kernel) getProcStat() ([]byte, error) {
	if _, err := os.Stat(k.statFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("kernel: %s does not exist!", k.statFile)
	} else if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(k.statFile)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func init() {
	agent.AddInput("kernel", &Kernel{
		statFile: "/proc/stat",
	})
}
