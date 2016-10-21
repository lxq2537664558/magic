package raw

import (
	"time"

	"github.com/aiyun/openapm/mecury/agent"
)

type Parser struct {
	Name      string
	FieldName string
}

func (p *Parser) Compile() error {
	return nil
}

func (p *Parser) ParseLine(line string) (agent.Metric, error) {
	f := make(map[string]interface{})
	f[p.FieldName] = line
	m, err := agent.NewMetric(p.Name, agent.Conf.Tags, f, time.Now())
	return m, err
}
