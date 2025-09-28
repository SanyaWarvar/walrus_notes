package cron

import "github.com/robfig/cron/v3"

type Cron struct {
	cron *cron.Cron
}

func NewCron() *Cron {
	return &Cron{
		cron: cron.New(),
	}
}

func (c *Cron) AddFunc(spec string, cmd func()) error {
	_, err := c.cron.AddFunc(spec, cmd)
	return err
}

func (c *Cron) Start() {
	c.cron.Start()
}

func (c *Cron) Stop() {
	c.cron.Stop()
}
