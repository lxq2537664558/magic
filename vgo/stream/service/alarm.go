package service

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

func (am *Alarmer) Compute(m Metrics) error {

	// Compute
	// for _, v := range m.Data {
	// 	for k, _ := range v.Fields {
	// 		log.Println("Name", k, v.Name+"."+k)
	// 	}
	// }
	// Alarm
	// log.Println("Alarmer Compute message is", m)
	return nil
}

func (am *Alarmer) compute() {

}
