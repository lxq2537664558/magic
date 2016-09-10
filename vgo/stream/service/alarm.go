package service

import "github.com/uber-go/zap"

type Alarmer struct {
}

func NewAlarm() *Alarmer {
	alarmer := &Alarmer{}
	return alarmer
}

func (am *Alarmer) Init() {

}

func (am *Alarmer) Start() {

}

func (am *Alarmer) Close() error {
	return nil
}

// 三种监控
// 平均值
// 瞬时平均值
// 状态存活监控

func (am *Alarmer) Compute(m Metrics) error {

	// Compute
	for _, v := range m.Data {
		hostname, ok := v.Tags["host"]
		if !ok {
			VLogger.Error("MetricData unfind hostname")
			continue
		}
		VLogger.Debug("Alarmer Compute", zap.String("hostname", hostname))
		streamer.hosts.RLock()

		streamer.hosts.RUnlock()
	}
	// Alarm
	// log.Println("Alarmer Compute message is", m)
	return nil
}

func (am *Alarmer) compute() {

}
