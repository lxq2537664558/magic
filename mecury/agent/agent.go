package agent

type Agent struct {
}

func New() *Agent {
	return &Agent{}
}

func (a *Agent) Init() {

}

func (a *Agent) Start(shutdown chan struct{}) {
	<-shutdown
}
