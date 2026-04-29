package container

import (
	"wn/internal/infrastructure/producer/layout"
)

func (c *Container) getProducers() *producers {
	if c.producers == nil {
		c.producers = &producers{c: c}
	}
	return c.producers
}

type producers struct {
	c *Container

	layoutProducer *layout.Producer
}

func (p *producers) getLayoutProducer() *layout.Producer {
	if p.layoutProducer == nil {
		p.layoutProducer = layout.NewProducer(
			p.c.getLogger(),
			p.c.getServices().getPermissionsService(),
			p.c.getServices().getSocketService(),
		)
	}
	return p.layoutProducer
}
