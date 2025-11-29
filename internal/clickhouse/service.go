package clickhouse

import (
	"context"
	"fmt"
	"pulse/internal/models"
)

// RecordCheckRun records a check run execution to ClickHouse
func (c *Client) RecordCheckRun(ctx context.Context, run *models.CheckRun) error {
	batch, err := c.conn.PrepareBatch(ctx, "INSERT INTO check_runs")
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}

	if err := batch.Append(
		run.ID,
		run.CheckID,
		run.Status,
		run.LatencyMs,
		run.StatusCode,
		run.Error,
		run.RunAt,
	); err != nil {
		return fmt.Errorf("failed to append to batch: %w", err)
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to send batch: %w", err)
	}

	return nil
}
