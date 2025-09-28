package file

import (
	"context"
	"wn/pkg/applogger"
	"wn/pkg/constants"
	"wn/pkg/util"
)

type fileService interface {
	GenerateStatics(ctx context.Context) error
}

type Cron struct {
	logger      applogger.Logger
	fileService fileService
}

func NewCron(logger applogger.Logger, fileService fileService) *Cron {
	return &Cron{
		logger:      logger,
		fileService: fileService,
	}
}

func (c *Cron) GenerateStatics() {

	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.ApiNameCtx, "GenerateStatics")
	ctx = context.WithValue(ctx, constants.RequestIdCtx, util.NewUUID().String())
	err := c.fileService.GenerateStatics(ctx)
	if err != nil {
		c.logger.WithCtx(ctx).Warnf("GenerateStatics: %s", err.Error())
	}
}
