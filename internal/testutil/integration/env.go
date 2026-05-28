//go:build integration

package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"wn/pkg/database/dragonfly"
	"wn/pkg/database/postgres"
	"wn/pkg/migrator"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

// DBConn реализует postgres.Connection без prometheus-метрик из postgres.New.
type DBConn struct {
	Pool *pgxpool.Pool
	CM   *postgres.ContextManager
}

func (db *DBConn) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	if tx := db.CM.ExtractTx(ctx); tx != nil {
		return tx.Exec(ctx, sql, arguments...)
	}
	return db.Pool.Exec(ctx, sql, arguments...)
}

func (db *DBConn) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if tx := db.CM.ExtractTx(ctx); tx != nil {
		return tx.Query(ctx, sql, args...)
	}
	return db.Pool.Query(ctx, sql, args...)
}

func (db *DBConn) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if tx := db.CM.ExtractTx(ctx); tx != nil {
		return tx.QueryRow(ctx, sql, args...)
	}
	return db.Pool.QueryRow(ctx, sql, args...)
}

// Env — общее окружение для integration-тестов.
type Env struct {
	Conn    *DBConn
	Pool    *postgres.Pool
	Tx      *postgres.PGTransactionManager
	Redis   *dragonfly.Client
	RootDir string
	cleanup func()
}

var (
	sharedEnv  *Env
	setupOnce  sync.Once
	setupError error
)

// Setup поднимает контейнеры один раз на пакет и очищает таблицы перед каждым тестом.
func Setup(t *testing.T) *Env {
	t.Helper()

	setupOnce.Do(func() {
		sharedEnv, setupError = startEnv()
	})

	if setupError != nil {
		t.Fatalf("integration setup: %v", setupError)
	}

	if err := cleanTables(context.Background(), sharedEnv.Conn.Pool); err != nil {
		t.Fatalf("clean tables: %v", err)
	}

	return sharedEnv
}

func startEnv() (*Env, error) {
	root, err := moduleRoot()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("postgres"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("postgres"),
		tcpostgres.WithInitScripts(filepath.Join(root, "deploy/postgres/init.sql")),
	)
	if err != nil {
		return nil, fmt.Errorf("postgres container: %w", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable", "search_path=walrus")
	if err != nil {
		return nil, fmt.Errorf("postgres connection string: %w", err)
	}

	dsn := postgres.GetCompleteDsn("postgres", "postgres", "", "", "postgres", "walrus", "disable")
	// testcontainers возвращает host:port — используем connStr для pool, dsn для migrate через pgx
	_ = dsn

	if err := migrator.Up(migrator.Config{
		MigrationsDirPath: filepath.Join(root, "migrations"),
		Dsn:               connStr,
		MigrationsTable:   "migrations",
		Schema:            "walrus",
		DBName:            "postgres",
	}); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("pgxpool: %w", err)
	}

	cm := postgres.NewContextManager()
	conn := &DBConn{Pool: pool, CM: cm}

	pgPool := (&postgres.Pool{Pool: pool}).WithContextManagerManager(cm)
	txManager := postgres.NewPGTransactionManager(pgPool, cm)

	redisContainer, err := tcredis.Run(ctx, "redis:7-alpine")
	if err != nil {
		return nil, fmt.Errorf("redis container: %w", err)
	}

	redisAddr, err := redisContainer.Endpoint(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("redis endpoint: %w", err)
	}

	redisClient, err := dragonfly.New(redisAddr, "", "")
	if err != nil {
		return nil, fmt.Errorf("redis client: %w", err)
	}

	cleanup := func() {
		redisClient.Close()
		_ = redisContainer.Terminate(ctx)
		pool.Close()
		_ = pgContainer.Terminate(ctx)
	}

	return &Env{
		Conn:    conn,
		Pool:    pgPool,
		Tx:      txManager,
		Redis:   redisClient,
		RootDir: root,
		cleanup: cleanup,
	}, nil
}

func cleanTables(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		TRUNCATE TABLE
			links,
			positions,
			notes,
			permissions,
			refresh_tokens,
			layouts,
			users
		RESTART IDENTITY CASCADE
	`)
	return err
}

func moduleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
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

// TestUserID генерирует уникальный идентификатор для тестовых данных.
func TestUserID() string {
	return fmt.Sprintf("u%d", time.Now().UnixNano())
}
