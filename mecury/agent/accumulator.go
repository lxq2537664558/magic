package agent

import (
	"log"
	"math"
	"time"
)

type Accumulator interface {
	// Create a point with a value, decorating it with tags
	// NOTE: tags is expected to be owned by the caller, don't mutate
	// it after passing to Add.
	Add(measurement string,
		value interface{},
		tags map[string]string,
		t ...time.Time)

	AddFields(measurement string,
		fields map[string]interface{},
		tags map[string]string,
		t ...time.Time)
}

type Accumulate struct {
	metricC     chan Metric
	inputConfig *InputConfig
	precision   time.Duration
}

func NewAccumulate(
	inputConfig *InputConfig,
	metrics chan Metric,
) *Accumulate {
	acc := Accumulate{
		metricC:     metrics,
		inputConfig: inputConfig,
	}
	return &acc
}

// By default, precision will be set to the same timestamp order as the
// collection interval, with the maximum being 1s.
//   ie, when interval = "10s", precision will be "1s"
//       when interval = "250ms", precision will be "1ms"
// Precision will NOT be used for service inputs. It is up to each individual
// service input to set the timestamp at the appropriate precision.
func (ac *Accumulate) SetPrecision(interval time.Duration) {
	switch {
	case interval >= time.Second:
		ac.precision = time.Second
	case interval >= time.Millisecond:
		ac.precision = time.Millisecond
	case interval >= time.Microsecond:
		ac.precision = time.Microsecond
	default:
		ac.precision = time.Nanosecond
	}
}

func (ac *Accumulate) Add(
	measurement string,
	value interface{},
	tags map[string]string,
	t ...time.Time,
) {

	fields := make(map[string]interface{})
	fields["value"] = value

	ac.AddFields(measurement, fields, tags, t...)
}

func (ac *Accumulate) AddFields(
	measurement string,
	fields map[string]interface{},
	tags map[string]string,
	t ...time.Time,
) {
	if len(fields) == 0 || len(measurement) == 0 {
		return
	}

	// Apply measurement prefix and suffix if set
	if len(ac.inputConfig.Prefix) != 0 {
		measurement = ac.inputConfig.Prefix + measurement
	}
	if len(ac.inputConfig.Suffix) != 0 {
		measurement = measurement + ac.inputConfig.Suffix
	}

	if tags == nil {
		tags = make(map[string]string)
	}

	// Apply plugin-wide tags if set
	for k, v := range ac.inputConfig.Tags {
		if _, ok := tags[k]; !ok {
			tags[k] = v
		}
	}
	// Apply daemon-wide tags if set
	for k, v := range Conf.Tags {
		if _, ok := tags[k]; !ok {
			tags[k] = v
		}
	}

	if !ac.inputConfig.Filter.ShouldTagsPass(tags) {
		return
	}

	result := make(map[string]interface{})
	for k, v := range fields {
		// Filter out any filtered fields
		if !ac.inputConfig.Filter.ShouldFieldsPass(k) {
			continue
		}
		// Validate uint64 and float64 fields
		switch val := v.(type) {
		case uint64:
			// InfluxDB does not support writing uint64
			if val < uint64(9223372036854775808) {
				result[k] = int64(val)
			} else {
				result[k] = int64(9223372036854775807)
			}
			continue
		case float64:
			// NaNs are invalid values in influxdb, skip measurement
			if math.IsNaN(val) || math.IsInf(val, 0) {
				continue
			}
		}

		result[k] = v
	}

	fields = nil
	if len(result) == 0 {
		return
	}

	var timestamp time.Time
	if len(t) > 0 {
		timestamp = t[0]
	} else {
		timestamp = time.Now()
	}
	timestamp = timestamp.Round(ac.precision)

	m, err := NewMetric(measurement, tags, result, timestamp)
	if err != nil {
		log.Printf("Error adding point [%s]: %s\n", measurement, err.Error())
		return
	}

	ac.metricC <- m
}

func (ac *Accumulate) DisablePrecision() {
	ac.precision = time.Nanosecond
}
