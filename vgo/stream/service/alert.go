package service

type AlertStatic struct {
	Type       int32 // 1 : average 2 : gauge 3: alive status
	Operator   int32 // 1:   >,   2:  = ,   3: <
	WarnValue  int32
	CritValue  int32
	WarnOutput string
	CritOutput string
	Duration   int32
	Template   string
}

func NewAlertStatic() *AlertStatic {
	return &AlertStatic{}
}

type AlertDynatic struct {
}

type Alert struct {
	AlertDy *AlertDynatic
	AlertSt *AlertStatic
}

func NewAlert() *Alert {
	alert := &Alert{
		AlertDy: &AlertDynatic{},
		AlertSt: &AlertStatic{},
	}
	return alert
}
