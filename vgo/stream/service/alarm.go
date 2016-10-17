package service

import (
	"log"
	"time"
)

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
			if group := streamer.groups.GetGroup(gid); group != nil {
				// VLogger.Info("GetGroup", zap.String("@Gid", gid), zap.String("@field", field), zap.String("@alert", metric.Name+"."+field))
				for field, _ := range metric.Fields {
					if alert, ok := group.Alerts[metric.Name+"."+field]; ok {
						switch fieldValue := metric.Fields[field].(type) {
						case float64:
							am.compute(alert, metric, fieldValue, m.Interval)
						default:
							log.Printf("%T\n", fieldValue)
						}
					}
				}
			}
		}
	}
	return nil
}

func (am *Alarmer) compute(alert *Alert, metric *MetricData, fieldValue float64, Interval int) (int, bool) {
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
		alert.AlertDy.computAverage(fieldValue)
	} else {
		alert.AlertSt.computGauge(fieldValue)
	}
	return 0, false
}
