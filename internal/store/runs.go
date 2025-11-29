package store

import (
	"context"
	"pulse/internal/clickhouse"
	"pulse/internal/models"

	"github.com/google/uuid"
)

type RunsStore struct {
	ch *clickhouse.Client
}

func NewRunsStore(ch *clickhouse.Client) *RunsStore {
	return &RunsStore{ch: ch}
}

func (s *RunsStore) CreateCheckRun(ctx context.Context, run *models.CheckRun) error {
	return s.ch.RecordCheckRun(ctx, run)
}

func (s *RunsStore) GetCheckRuns(ctx context.Context, checkID uuid.UUID, limit int) ([]*models.CheckRun, error) {
	return s.ch.GetCheckRuns(ctx, checkID, limit)
}
