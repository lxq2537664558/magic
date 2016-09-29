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

	for _, metric := range m.Data {
		hostname, ok := metric.Tags["host"]
		if !ok {
			VLogger.Error("MetricData unfind hostname")
			continue
		}
		VLogger.Info("Compute", zap.String("@hostname", hostname))
		// streamer.groups.GetGroups(hostname)
		groups := streamer.hostsTogroups.Get(hostname)
		if groups == nil {
			// VLogger.Debug("unfind groups", zap.String("@hostname", hostname))
			continue
		}

		VLogger.Debug("Compute", zap.Object("@groups", groups))
		// for gid, _ := range groups {
		// 	if g := GetGroup(gid); g != nil {
		// 		for field, _ := range metric.Fields {
		// 			if alert, ok := g.Alerts[metric.Name+"."+field]; ok {
		// 				log.Println(alert)
		// 			}
		// 		}
		// 	}
		// }
	}
	return nil
}

func (am *Alarmer) compute() {

}
