package service

import (
	"log"
	"time"

	"github.com/uber-go/zap"
)

const (
	// 不报警
	NoAlarm int = iota
	// 警告
	WarnAlarm
	// 严重警告
	CritAlarm
)

// AlertStatic alert 静态数据变量，存放报警策略信息
type AlertStatic struct {
	Type       int32 // 1 : average 2 : gauge 3: alive status
	Operator   int32 // 1:   >,   2:  = ,   3: <
	WarnValue  float64
	CritValue  float64
	WarnOutput string
	CritOutput string
	Duration   int32
	Template   string
}

func NewAlertStatic() *AlertStatic {
	return &AlertStatic{}
}

// AlertDynatic alert动态数据变量，动态存放上报的数据
type AlertDynatic struct {
	StartTime int64     // 插入时间
	Len       int32     // 总长度
	NowIndex  int       // 当前下标
	RingArray []float64 // 存储列表
}

func NewAlertDynatic() *AlertDynatic {
	return &AlertDynatic{}
}

type Alert struct {
	AlertDy *AlertDynatic
	AlertSt *AlertStatic
}

func NewAlert() *Alert {
	alert := &Alert{}
	return alert
}

func (a *Alert) computAverage(reportData float64, originalGroup *Group) {
	a.AlertDy.RingArray[a.AlertDy.NowIndex] = reportData
	if a.AlertDy.NowIndex >= (int(a.AlertDy.Len) - 1) {
		var total float64
		for _, data := range a.AlertDy.RingArray {
			total += data
		}
		average := (total / float64(a.AlertDy.Len))
		// 检查是否到达报警阀值
		alarmflg := checkIsAlarm(a.AlertSt.Operator, a.AlertSt.WarnValue, a.AlertSt.CritValue, average)
		if alarmflg != NoAlarm {
			if alarmflg == WarnAlarm {
				VLogger.Info("Alarm", zap.String("@WarnOutput", a.AlertSt.WarnOutput), zap.Float64("@average", average), zap.Float64("@WarnValue", a.AlertSt.WarnValue))
			} else if alarmflg == CritAlarm {
				VLogger.Info("Alarm", zap.String("@CritOutput", a.AlertSt.CritOutput), zap.Float64("@average", average), zap.Float64("@CritValue", a.AlertSt.CritValue))
			}
		}
		log.Println("是否需要报警:  ", alarmflg)
	}
	VLogger.Info("InsertAndComput", zap.Object("@AlertDynatic", a.AlertDy))
	a.AlertDy.NowIndex++
	// 回滚覆盖
	if a.AlertDy.NowIndex >= int(a.AlertDy.Len) {
		a.AlertDy.NowIndex = 0
		a.AlertDy.StartTime = time.Now().Unix()
	}
}

// computGauge 检查瞬时值
func (a *Alert) computGauge(reportData float64, originalGroup *Group) {
	if reportData == a.AlertSt.WarnValue {
		VLogger.Info("Alarm", zap.String("@WarnOutput", a.AlertSt.WarnOutput), zap.Float64("@reportData", reportData), zap.Float64("@WarnValue", a.AlertSt.WarnValue))
	} else if reportData == a.AlertSt.CritValue {
		VLogger.Info("Alarm", zap.String("@CritOutput", a.AlertSt.CritOutput), zap.Float64("@reportData", reportData), zap.Float64("@CritValue", a.AlertSt.CritValue))
	}
}

// checkIsAlarm return alarmflag and alarm type  0 no alarm, 1 warn alarm , 2 crit alarm
func checkIsAlarm(operator int32, warnValue float64, critValue float64, getValue float64) int {
	if operator == 1 {
		if getValue > warnValue && getValue < critValue {
			return WarnAlarm
		} else if getValue > critValue {
			return CritAlarm
		}
	} else if operator == 2 {
		if getValue == warnValue {
			return WarnAlarm
		} else if getValue == critValue {
			return CritAlarm
		}
	} else if operator == 3 {
		if getValue < warnValue {
			return WarnAlarm
		} else if getValue > warnValue && getValue < critValue {
			return CritAlarm
		}
	}
	return NoAlarm
}
