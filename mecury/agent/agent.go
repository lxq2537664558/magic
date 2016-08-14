package agent

type Agent struct {
}

func New() *Agent {
	// init config
	Conf = &Config{}
	Conf.Load()

	return &Agent{}
}

func (a *Agent) Init() {

}

func (a *Agent) Start(shutdown chan struct{}) {
	<-shutdown
}
