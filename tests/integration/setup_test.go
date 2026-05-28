//go:build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"wn/internal/domain/entity"
	"wn/pkg/constants"
	"wn/pkg/database/dragonfly"
	"wn/pkg/database/postgres"
	"wn/pkg/migrator"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

var env *Env

type Env struct {
	Ctx     context.Context
	Pool    *pgxpool.Pool
	DB      postgres.Connection
	CM      *postgres.ContextManager
	Trx     *postgres.PGTransactionManager
	Redis   *dragonfly.Client
	cleanup func()
}

type dbConn struct {
	pool *pgxpool.Pool
	cm   *postgres.ContextManager
}

func (c *dbConn) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	if tx := c.cm.ExtractTx(ctx); tx != nil {
		return tx.Exec(ctx, sql, arguments...)
	}
	return c.pool.Exec(ctx, sql, arguments...)
}

func (c *dbConn) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if tx := c.cm.ExtractTx(ctx); tx != nil {
		return tx.Query(ctx, sql, args...)
	}
	return c.pool.Query(ctx, sql, args...)
}

func (c *dbConn) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if tx := c.cm.ExtractTx(ctx); tx != nil {
		return tx.QueryRow(ctx, sql, args...)
	}
	return c.pool.QueryRow(ctx, sql, args...)
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	root, err := moduleRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "module root: %v\n", err)
		os.Exit(1)
	}

	pgContainer, pgDSN, err := startPostgres(ctx, root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "postgres: %v\n", err)
		os.Exit(1)
	}

	if err := waitForPostgres(ctx, pgDSN); err != nil {
		fmt.Fprintf(os.Stderr, "postgres wait: %v\n", err)
		os.Exit(1)
	}

	if err := migrator.Up(migrator.Config{
		MigrationsDirPath: filepath.Join(root, "migrations"),
		Dsn:               pgDSN,
		MigrationsTable:   "migrations",
		Schema:            "walrus",
		DBName:            "postgres",
	}); err != nil {
		fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
		os.Exit(1)
	}

	pool, db, cm, err := newPool(ctx, pgDSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pool: %v\n", err)
		os.Exit(1)
	}

	redisContainer, redisClient, err := startRedis(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "redis: %v\n", err)
		os.Exit(1)
	}

	env = &Env{
		Ctx:   ctx,
		Pool:  pool,
		DB:    db,
		CM:    cm,
		Trx:   postgres.NewPGTransactionManager(&postgres.Pool{Pool: pool}, cm),
		Redis: redisClient,
		cleanup: func() {
			pool.Close()
			redisClient.Close()
			_ = pgContainer.Terminate(ctx)
			_ = redisContainer.Terminate(ctx)
		},
	}

	code := m.Run()
	env.cleanup()
	os.Exit(code)
}

func waitForPostgres(ctx context.Context, dsn string) error {
	var lastErr error
	for range 30 {
		pool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			lastErr = err
			time.Sleep(500 * time.Millisecond)
			continue
		}
		lastErr = pool.Ping(ctx)
		pool.Close()
		if lastErr == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("postgres not ready: %w", lastErr)
}

func newPool(ctx context.Context, dsn string) (*pgxpool.Pool, postgres.Connection, *postgres.ContextManager, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, nil, nil, err
	}
	cfg.MaxConns = 5
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, nil, nil, err
	}
	cm := postgres.NewContextManager()
	return pool, &dbConn{pool: pool, cm: cm}, cm, nil
}

func startPostgres(ctx context.Context, root string) (*tcpostgres.PostgresContainer, string, error) {
	initScript := filepath.Join(root, "deploy", "postgres", "init.sql")
	container, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("postgres"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("postgres"),
		tcpostgres.WithInitScripts(initScript),
	)
	if err != nil {
		return nil, "", err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, "", err
	}
	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, "", err
	}

	dsn := fmt.Sprintf(
		"user=postgres password=postgres host=%s port=%s dbname=postgres search_path=walrus sslmode=disable",
		host, port.Port(),
	)
	return container, dsn, nil
}

func startRedis(ctx context.Context) (*tcredis.RedisContainer, *dragonfly.Client, error) {
	container, err := tcredis.Run(ctx, "redis:7-alpine")
	if err != nil {
		return nil, nil, err
	}
	endpoint, err := container.Endpoint(ctx, "")
	if err != nil {
		return nil, nil, err
	}
	client, err := dragonfly.New(endpoint, "", "")
	if err != nil {
		return nil, nil, err
	}
	return container, client, nil
}

func moduleRoot() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("runtime.Caller failed")
	}
	dir := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}

func resetDB(t *testing.T) {
	t.Helper()
	_, err := env.Pool.Exec(env.Ctx, `
		TRUNCATE TABLE
			refresh_tokens,
			permissions,
			links,
			positions,
			notes,
			layouts,
			users
		RESTART IDENTITY CASCADE
	`)
	if err != nil {
		t.Fatalf("resetDB: %v", err)
	}
}

func newTestUser() *entity.User {
	id := uuid.New()
	return &entity.User{
		Id:             id,
		Username:       "user_" + id.String()[:8],
		Email:          id.String() + "@integration.test",
		Password:       "hashed-password",
		Role:           constants.ClientRole,
		ImgUrl:         "base.png",
		ConfirmedEmail: false,
		CreatedAt:      time.Now().UTC(),
	}
}
