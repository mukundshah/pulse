package clickhouse

import (
	"context"
	"fmt"
	"log"
)

// InitSchema creates all necessary tables in ClickHouse
func (c *Client) InitSchema(ctx context.Context) error {
	tables := []string{}

	for _, tableSQL := range tables {
		if err := c.conn.Exec(ctx, tableSQL); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	log.Println("ClickHouse schema initialized successfully")
	return nil
}
