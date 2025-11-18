package container

import (
	"wn/config"
	"wn/internal/endpoint/controller/http"
	v1 "wn/internal/endpoint/controller/http/api/v1"
	"wn/pkg/applogger"
	"wn/pkg/database/dragonfly"
	"wn/pkg/database/postgres"
	"wn/pkg/httpserver"
	"wn/pkg/response"
	"wn/pkg/restclient"
	"wn/pkg/trx"
)

type Container struct {
	cfg                *config.Config
	cacheClient        *dragonfly.Client
	pool               *postgres.Pool
	poolContextManager *postgres.ContextManager
	logger             applogger.Logger
	builder            *response.Builder
	trxManager         trx.TransactionManager
	httpServer         *httpserver.Server
	httpDispatcher     *v1.Dispatcher
	httpKernel         *http.Kernel
	restClient         restclient.RestClient

	repositories *repositories
	applications *applications
	services     *services
	workers      *workers
	caches       *cache
}

func New(cfg *config.Config) *Container {
	return &Container{cfg: cfg}
}

func (c *Container) getConfig() *config.Config {
	return c.cfg
}

func (c *Container) Start() error {
	//migrate
	if err := c.Migrate(); err != nil {
		return err
	}

	//servers
	c.getHTTPServer().Start()

	//workers
	if err := c.getWorkers().start(); err != nil {
		return err
	}
	c.getServices().getSocketService().RegisterHandler(
		"UPDATE_DRAFT_REQUEST", c.getServices().getNoteService().HandleCreateDraft,
	)
	return nil
}

func (c *Container) Stop() error {
	if err := c.getHTTPServer().Shutdown(); err != nil {
		return err
	}
	return nil
}
