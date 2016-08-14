package agent

import "time"

type Config struct {
	Common *CommonConfig

	Agent *AgentConfig

	Inputs  []*InputConfig
	Outputs []*OutputConfig
}

var Conf *Config

func (c *Config) Load() {

}

func Reload(r chan struct{}) {
	time.Sleep(5 * time.Second)
	r <- struct{}{}
}

type CommonConfig struct {
}

type AgentConfig struct {
	// Interval at which to gather information
	Interval time.Duration

	// RoundInterval rounds collection interval to 'interval'.
	//     ie, if Interval=10s then always collect on :00, :10, :20, etc.
	RoundInterval bool

	// By default, precision will be set to the same timestamp order as the
	// collection interval, with the maximum being 1s.
	//   ie, when interval = "10s", precision will be "1s"
	//       when interval = "250ms", precision will be "1ms"
	// Precision will NOT be used for service inputs. It is up to each individual
	// service input to set the timestamp at the appropriate precision.
	Precision time.Duration

	// CollectionJitter is used to jitter the collection by a random amount.
	// Each plugin will sleep for a random time within jitter before collecting.
	// This can be used to avoid many plugins querying things like sysfs at the
	// same time, which can have a measurable effect on the system.
	CollectionJitter time.Duration

	// FlushInterval is the Interval at which to flush data
	FlushInterval time.Duration

	// FlushJitter Jitters the flush interval by a random amount.
	// This is primarily to avoid large write spikes for users running a large
	// number of telegraf instances.
	// ie, a jitter of 5s and interval 10s means flushes will happen every 10-15s
	FlushJitter time.Duration

	// MetricBatchSize is the maximum number of metrics that is wrote to an
	// output plugin in one call.
	MetricBatchSize int

	// MetricBufferLimit is the max number of metrics that each output plugin
	// will cache. The buffer is cleared when a successful write occurs. When
	// full, the oldest metrics will be overwritten. This number should be a
	// multiple of MetricBatchSize. Due to current implementation, this could
	// not be less than 2 times MetricBatchSize.
	MetricBufferLimit int

	// FlushBufferWhenFull tells Telegraf to flush the metric buffer whenever
	// it fills up, regardless of FlushInterval. Setting this option to true
	// does _not_ deactivate FlushInterval.
	FlushBufferWhenFull bool

	// TODO(cam): Remove UTC and parameter, they are no longer
	// valid for the agent config. Leaving them here for now for backwards-
	// compatability
	UTC bool `toml:"utc"`

	Hostname string

	OmitHostname bool

	// default tags
	Tags map[string]string
}

type InputConfig struct {
	Name   string
	Prefix string
	Suffix string

	Input Inputer

	Tags     map[string]string
	Filter   InputFilter
	Interval time.Duration
}

type OutputConfig struct {
	Name string

	Output Outputer

	Metrics *Buffer
}
