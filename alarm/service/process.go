package service

import (
	"log"

	"time"

	"github.com/nats-io/nats"
)

//easyjson:json
type AlertData struct {
	ID       string  `json:"id"` // metric name + field
	GroupID  string  `json:"gid"`
	Value    float64 `json:"v"`
	Level    int     `json:"l"` //0: warn, 1 : critical
	HostName string  `json:"h"`
}

func process(m *nats.Msg) {
	a := &AlertData{}
	a.UnmarshalJSON(m.Data)

	gs.RLock()
	group := gs.groups[a.GroupID]
	alert := group.Alerts[a.ID]
	gs.RUnlock()

	// 判断当前时间是否超出允许的报警信息更新间隔
	now := time.Now()
	if now.Sub(alert.LastTime[a.Level]) > alert.Duration[a.Level] {
		// 清空之前的报警历史数据
		alert.NowCount[a.Level] = 1
		alert.LastTime[a.Level] = now
		return
	}

	if alert.NowCount[a.Level]+1 >= alert.Count[a.Level] {
		log.Println(alert.Count[a.Level])
		output := Conf.Outputs[alert.AlarmOutput[a.Level]]
		// 报警
		for _, u := range group.Users {
			data := &Alarm{
				Data: m.Data,
				User: u.Info[alert.AlarmOutput[a.Level]],
			}
			output.Write(data)
		}
		//清空当前count
		alert.NowCount[a.Level] = 0
	} else {
		// 不满足条件，计数＋1
		alert.NowCount[a.Level]++
	}

	// 更新报警数据的更新时间
	alert.LastTime[a.Level] = now
}
