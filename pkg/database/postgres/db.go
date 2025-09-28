package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"time"
)

const (
	defaultMaxPoolSize  = 1
	defaultConnAttempts = 2
	defaultConnTimeout  = time.Second
	connMaxLifetime     = 3 * time.Second
)

// Pool -.
type Pool struct {
	maxPoolSize       int
	minPoolSize       int
	connMaxLifetime   time.Duration
	connMaxIdletime   time.Duration
	healthCheckPeriod time.Duration
	connAttempts      int
	connTimeout       time.Duration

	Pool       *pgxpool.Pool
	ctxManager *ContextManager
}

func GetCompleteUrl(userName, password, host, port, DBName, schema, sslMode, sslRootCert string) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?search_path=%s",
		userName,
		password,
		host,
		port,
		DBName,
		schema)
}

func GetCompleteDsn(userName, password, host, port, DBName, schema, sslMode string) string {
	str := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s search_path=%s sslmode=%s",
		userName,
		password,
		host,
		port,
		DBName,
		schema,
		sslMode)
	return str
}

func New(url string, opts ...Option) (*Pool, error) {
	pg := &Pool{
		maxPoolSize:     defaultMaxPoolSize,
		connAttempts:    defaultConnAttempts,
		connTimeout:     defaultConnTimeout,
		connMaxLifetime: connMaxLifetime,
	}

	// Custom options
	for _, opt := range opts {
		opt(pg)
	}
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}
	poolConfig.MaxConns = int32(pg.maxPoolSize)
	poolConfig.MinConns = int32(pg.minPoolSize)
	poolConfig.MaxConnIdleTime = pg.connMaxIdletime
	poolConfig.MaxConnLifetime = pg.connMaxLifetime
	poolConfig.HealthCheckPeriod = pg.healthCheckPeriod

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			collector := NewCollector(pg.Pool, map[string]string{"db_name": "my_db"})
			prometheus.MustRegister(collector)
			break
		}

		log.Printf("Pool is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}
	return pg, nil
}

func (p *Pool) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

func (p *Pool) WithContextManagerManager(ctxManager *ContextManager) *Pool {
	if p.ctxManager == nil {
		p.ctxManager = ctxManager
	}
	return p
}

func (p *Pool) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	if p.ctxManager != nil {
		if tx := p.ctxManager.ExtractTx(ctx); tx != nil {
			return tx.Exec(ctx, sql, arguments...)
		}
	}
	return p.Pool.Exec(ctx, sql, arguments...)
}

func (p *Pool) QueryRow(ctx context.Context, sql string, arguments ...any) pgx.Row {
	if p.ctxManager != nil {
		if tx := p.ctxManager.ExtractTx(ctx); tx != nil {
			return tx.QueryRow(ctx, sql, arguments...)
		}
	}
	return p.Pool.QueryRow(ctx, sql, arguments...)
}

func (p *Pool) Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error) {
	if p.ctxManager != nil {
		if tx := p.ctxManager.ExtractTx(ctx); tx != nil {
			return tx.Query(ctx, sql, arguments...)
		}
	}
	return p.Pool.Query(ctx, sql, arguments...)
}
