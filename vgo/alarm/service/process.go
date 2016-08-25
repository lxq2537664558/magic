package service

import "github.com/nats-io/nats"

//easyjson:json
type AlertData struct {
	ID      string  `json:"id"`
	GroupID string  `json:"gid"`
	Value   float64 `json:"v"`
	Level   int     `json:"l"` //0: warn, 1 : critical
}

func process(m *nats.Msg) {
	a := &AlertData{}
	a.UnmarshalJSON(m.Data)

	gs.RLock()
	group := gs.groups[a.GroupID]
	alert := group.Alerts[a.ID]
	gs.RUnlock()

	if alert.NowCount[a.Level]+1 >= alert.Count[a.Level] {
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

}
