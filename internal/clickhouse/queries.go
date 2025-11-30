package clickhouse

import (
	"context"
	"fmt"
	"pulse/internal/models"

	"github.com/google/uuid"
)

// GetCheckRuns retrieves check run history for a specific check
func (c *Client) GetCheckRuns(ctx context.Context, checkID uuid.UUID, limit int) ([]*models.CheckRun, error) {
	query := `
		SELECT id, check_id, status, latency_ms, status_code, error, run_at
		FROM check_runs
		WHERE check_id = ?
		ORDER BY run_at DESC
		LIMIT ?
	`

	rows, err := c.conn.Query(ctx, query, checkID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query check runs: %w", err)
	}
	defer rows.Close()

	var runs []*models.CheckRun
	for rows.Next() {
		var run models.CheckRun
		err := rows.Scan(
			&run.ID,
			&run.CheckID,
			&run.Status,
			&run.LatencyMs,
			&run.StatusCode,
			&run.Error,
			&run.RunAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		runs = append(runs, &run)
	}

	return runs, rows.Err()
}

// GetAverageLatency returns the average latency in milliseconds for all check runs
func (c *Client) GetAverageLatency(ctx context.Context) (float64, error) {
	query := `SELECT avg(latency_ms) FROM check_runs`

	var avgLatency float64
	err := c.conn.QueryRow(ctx, query).Scan(&avgLatency)
	if err != nil {
		return 0, fmt.Errorf("failed to query average latency: %w", err)
	}

	return avgLatency, nil
}
