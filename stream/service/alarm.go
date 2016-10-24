package service

import "time"

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
		groups := streamer.hostsTogroups.Get(hostname)
		if groups == nil {
			continue
		}
		for gid, _ := range groups {
			// VLogger.Info("Compute", zap.String("@show——gid", gid), zap.Int("@len", len(groups)))
			if group := streamer.groups.GetGroup(gid); group != nil {
				originalGroup := group
				for {
					if !group.ComputAlarm(am, metric, m.Interval, originalGroup) {
						if group.Parent != nil {
							group = group.Parent
						} else {
							// VLogger.Info("Compute", zap.String("@show2	", group.ID))
							break
						}
					} else {
						// VLogger.Info("Compute", zap.String("@show", group.ID))
						break
					}
				}
			}
		}
	}
	return nil
}

func (am *Alarmer) compute(alert *Alert, metric *MetricData, fieldValue float64, Interval int, originalGroup *Group) (int, bool) {
	if alert.AlertSt.Type == 1 {
		if alert.AlertDy == nil {
			len := alert.AlertSt.Duration / int32(Interval)
			if len == 0 {
				len = 1
			}
			alert.AlertDy = &AlertDynatic{
				StartTime: time.Now().Unix(),
				NowIndex:  0,
				Len:       len,
				RingArray: make([]float64, len),
			}
		}
		// alert.computAverage(fieldValue, originalGroup)
		alert.computAverage(fieldValue, originalGroup)
	} else {
		alert.computGauge(fieldValue, originalGroup)
	}
	return 0, false
}
