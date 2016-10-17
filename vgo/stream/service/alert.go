package service

import (
	"log"
	"time"

	"github.com/uber-go/zap"
)

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

func (as *AlertStatic) computGauge(reportData float64) {
	// 检查是否到达报警阀值
	// if reportData == as.WarnValue {

	// } else if reportData == as.CritValue {

	// }
}

type AlertDynatic struct {
	StartTime int64     // 插入时间
	Len       int32     // 总长度
	NowIndex  int       // 当前下标
	RingArray []float64 // 存储列表
}

func NewAlertDynatic() *AlertDynatic {
	return &AlertDynatic{}
}

func (ad *AlertDynatic) computAverage(reportData float64) {
	ad.RingArray[ad.NowIndex] = reportData
	if ad.NowIndex >= (int(ad.Len) - 1) {
		var total float64
		for _, data := range ad.RingArray {
			total += data
		}
		average := (total / float64(ad.Len))
		log.Println("计算平均值:  ", average)
		// 检查是否到达报警阀值
		// checkIsAlarm(operator int, setValue float64, getValue float64)
	}
	VLogger.Info("InsertAndComput", zap.Object("@AlertDynatic", ad))
	ad.NowIndex++
	// 回滚覆盖
	if ad.NowIndex >= int(ad.Len) {
		ad.NowIndex = 0
		ad.StartTime = time.Now().Unix()
	}
}

type Alert struct {
	AlertDy *AlertDynatic
	AlertSt *AlertStatic
}

func NewAlert() *Alert {
	alert := &Alert{}
	return alert
}

func checkIsAlarm(operator int, setValue float64, getValue float64) bool {
	if operator == 1 {
		if getValue > setValue {
			return true
		}
	} else if operator == 2 {
		if getValue == setValue {
			return true
		}
	} else if operator == 3 {
		if getValue < setValue {
			return true
		}
	}
	return false
}
