package postgresql

import (
	"context"
	"fmt"
	"hiv_mind/pkg/logger"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgreSQL struct {
	DB *pgxpool.Pool
}

type PoolConfig struct {
	MaxOpenConns    int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

func NewPostgreSQL(ctx context.Context, dsn string, poolCfg PoolConfig) (*PostgreSQL, error) {
	lg := logger.Get()
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		lg.Error("unable to open database connection: %w", zap.Error(err))
		return nil, err
	}

	config.MaxConns = poolCfg.MaxOpenConns
	config.MinConns = poolCfg.MinConns
	config.MaxConnLifetime = poolCfg.MaxConnLifetime
	config.MaxConnIdleTime = poolCfg.MaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		lg.Error("unable to create connection pool: %w", zap.Error(err))
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &PostgreSQL{DB: pool}, nil

}

func (db *PostgreSQL) PingContext(ctx context.Context) error {
	return db.DB.Ping(ctx)
}

func (db *PostgreSQL) Close() {
	if db.DB != nil {
		db.DB.Close()
	}
}
