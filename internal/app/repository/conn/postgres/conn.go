package rcpostgres

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/AlexanderTal/order-service/internal/app/config/section"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Client struct {
	db  *gorm.DB
	cfg section.RepositoryPostgres
}

func NewClient(ctx context.Context, cfg section.RepositoryPostgres) (*Client, error) {
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.Username, cfg.Password),
		Host:   cfg.Address,
		Path:   cfg.Name,
	}
	q := dsn.Query()
	q.Add("sslmode", "disable")
	dsn.RawQuery = q.Encode()

	db, err := gorm.Open(postgres.Open(dsn.String()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open gorm connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(10)

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return &Client{
		db:  db,
		cfg: cfg,
	}, nil
}

func (c *Client) DB() *gorm.DB {
	return c.db
}

func (c *Client) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB for closing: %w", err)
	}
	return sqlDB.Close()
}

func (c *Client) GetDB(ctx context.Context) *gorm.DB {
	if tx := getTxFromCtx(ctx); tx != nil {
		return tx
	}
	return c.db
}

func (c *Client) InsideTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if getTxFromCtx(ctx) != nil {
		return fn(ctx)
	}

	return c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := ctxWithTx(ctx, tx)
		return fn(txCtx)
	})
}
