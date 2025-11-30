package clickhouse

import (
	"context"
	"fmt"
	"log"
	"pulse/internal/config"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Client struct {
	conn driver.Conn
}

func NewClient(cfg *config.Config) (*Client, error) {
	opts, err := clickhouse.ParseDSN(cfg.ClickHouseDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ClickHouse DSN: %w", err)
	}

	conn, err := clickhouse.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	ctx := context.Background()
	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	log.Println("Connected to ClickHouse")
	return &Client{conn: conn}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Conn() driver.Conn {
	return c.conn
}

// HealthCheck pings ClickHouse to verify connectivity
func (c *Client) HealthCheck(ctx context.Context) error {
	return c.conn.Ping(ctx)
}
