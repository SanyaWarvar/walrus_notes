package container

import (
	"fmt"
	"wn/internal/endpoint/worker/file"
	"wn/pkg/cron"
)

type workers struct {
	c  *Container
	cr *cron.Cron

	file *file.Cron
}

func (c *Container) getWorkers() *workers {
	if c.workers == nil {
		c.workers = &workers{
			cr: cron.NewCron(),
			c:  c,
		}
	}
	return c.workers
}

func (w *workers) getFileJob() *file.Cron {
	if w.file == nil {
		w.file = file.NewCron(w.c.getLogger(), w.c.getServices().getFileService())
	}
	return w.file
}

func (w *workers) start() error {
	if err := w.cr.AddFunc(w.c.getConfig().Cron.GenerateStatics, w.getFileJob().GenerateStatics); err != nil {
		return fmt.Errorf("GenerateStatics: %v", err)
	}
	w.cr.Start()
	return nil
}

func (w *workers) stop() {
	w.cr.Stop()
}
