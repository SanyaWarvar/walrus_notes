package container

import (
	"log"
	"wn/internal/endpoint/controller/http"
	"wn/pkg/applogger"
	"wn/pkg/database/dragonfly"
	"wn/pkg/database/postgres"
	"wn/pkg/httpserver"
	"wn/pkg/migrator"
	"wn/pkg/response"
	"wn/pkg/restclient"
	"wn/pkg/trx"
)

func (c *Container) Migrate() error {
	dsn := postgres.GetCompleteDsn(
		c.cfg.Postgres.Username,
		c.cfg.Postgres.Password,
		c.cfg.Postgres.Host,
		c.cfg.Postgres.Port,
		c.cfg.Postgres.DBName,
		c.cfg.Postgres.Schema,
		c.cfg.Postgres.SSLMode,
	)
	conf := migrator.Config{
		MigrationsDirPath: "migrations",
		Dsn:               dsn,
		MigrationsTable:   "migrations",
		Schema:            c.cfg.Postgres.Schema,
		DBName:            c.cfg.Postgres.DBName,
	}
	return migrator.Up(conf)
}

func (c *Container) getContextManager() *postgres.ContextManager {
	if c.poolContextManager == nil {
		c.poolContextManager = postgres.NewContextManager()
	}
	return c.poolContextManager
}

func (c *Container) getTransactionManager() trx.TransactionManager {
	if c.trxManager == nil {
		c.trxManager = postgres.NewPGTransactionManager(c.getDBPool(), c.getContextManager())
		c.pool = c.getDBPool().WithContextManagerManager(c.getContextManager())
	}
	return c.trxManager
}

func (c *Container) getKernel() *http.Kernel {
	if c.httpKernel == nil {
		c.httpKernel = http.NewKernel(
			c.getConfig().Internal.LogInputParamOnErr,

			c.getLogger(),
			c.getResponseBuilder(),
			c.getHTTPDispatcher(),
		)
	}
	return c.httpKernel
}

func (c *Container) getCacheClient() *dragonfly.Client {
	if c.cacheClient == nil {
		client, err := dragonfly.New(
			c.getConfig().Cache.Url,
			c.getConfig().Cache.Username,
			c.getConfig().Cache.Password,
		)
		if err != nil {
			log.Fatalf("getCacheClient: %v", err)
		}
		c.cacheClient = client
	}
	return c.cacheClient
}

func (c *Container) getHTTPServer() *httpserver.Server {
	if c.httpServer == nil {
		c.httpServer = httpserver.New(
			c.getKernel().Init(),
			httpserver.Port(c.getConfig().HTTP.Port),
			httpserver.ReadTimeout(c.getConfig().HTTP.ReadTimeout),
			httpserver.WriteTimeout(c.getConfig().HTTP.WriteTimeout),
		)
	}
	return c.httpServer
}

func (c *Container) getDBPool() *postgres.Pool {
	if c.pool == nil {
		pool, err := postgres.New(postgres.GetCompleteDsn(
			c.getConfig().Postgres.Username,
			c.getConfig().Postgres.Password,
			c.getConfig().Postgres.Host,
			c.getConfig().Postgres.Port,
			c.getConfig().Postgres.DBName,
			c.getConfig().Postgres.Schema,
			c.getConfig().Postgres.SSLMode),
			postgres.MaxPoolSize(c.getConfig().Postgres.PoolMax),
			postgres.MinPoolSize(c.getConfig().Postgres.PoolMin),
			postgres.ConnMaxLifetime(c.getConfig().Postgres.ConnectionMaxLifeTime),
			postgres.ConnMaxIdletime(c.getConfig().Postgres.ConnectionMaxIdleTime),
			postgres.HealthCheckPeriod(c.getConfig().Postgres.HealthCheckPeriod),
		)
		if err != nil {
			log.Fatalf("getDBPool: %v", err)
		}
		c.pool = pool
	}
	return c.pool
}

func (c *Container) getLogger() applogger.Logger {
	if c.logger == nil {
		logger, err := applogger.NewLogger(c.cfg.Log.Level)
		if err != nil {
			log.Fatalf("NewLogger: %v", err)
		}
		c.logger = logger
	}
	return c.logger
}

func (c *Container) getResponseBuilder() *response.Builder {
	if c.builder == nil {
		c.builder = response.NewResponseBuilder(c.getConfig().Response.ExportError)
	}
	return c.builder
}

func (c *Container) getRestClient() restclient.RestClient {
	if c.restClient == nil {
		c.restClient = restclient.NewRestClient(c.getLogger(), c.getConfig().Log.RequestLogEnabled, c.getConfig().Log.RequestLogWithBody)
	}
	return c.restClient
}
