package clickhouse

import (
	"context"
	"fmt"
	"log"
)

// InitSchema creates all necessary tables in ClickHouse
func (c *Client) InitSchema(ctx context.Context) error {
	tables := []string{
		createCheckRunsTable,
	}

	for _, tableSQL := range tables {
		if err := c.conn.Exec(ctx, tableSQL); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	log.Println("ClickHouse schema initialized successfully")
	return nil
}

const createCheckRunsTable = `
CREATE TABLE IF NOT EXISTS check_runs (
    id UUID,
    check_id UUID,
    status String,
    latency_ms Int64,
    status_code Int32,
    error Nullable(String),
    run_at DateTime
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(run_at)
ORDER BY (check_id, run_at)
SETTINGS index_granularity = 8192;
`
