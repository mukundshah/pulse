package store

import (
	"encoding/json"
	"sort"
	"time"

	"pulse/internal/models"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func (s *Store) CreateCheckRun(run *models.CheckRun) error {
	return s.db.Create(run).Error
}

func (s *Store) GetCheckRun(id uuid.UUID) (*models.CheckRun, error) {
	var run models.CheckRun
	if err := s.db.Preload("Check.Project").Preload("Region").First(&run, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &run, nil
}

func (s *Store) GetCheckRunsByCheck(checkID uuid.UUID, limit int, after, before *uuid.UUID, startTime, endTime *time.Time) ([]models.CheckRun, error) {
	var runs []models.CheckRun
	query := s.db.Preload("Region").Omit("Check").Where("check_id = ?", checkID)

	// Apply date range filter if provided
	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	// Handle cursor-based pagination
	// We sort by created_at DESC, id DESC (latest to oldest)
	if after != nil {
		// Get the cursor item to find its created_at
		var cursorRun models.CheckRun
		if err := s.db.First(&cursorRun, "id = ?", *after).Error; err == nil {
			// Get items after the cursor (older items in DESC order)
			// (created_at < cursor.created_at) OR (created_at = cursor.created_at AND id < cursor.id)
			query = query.Where(
				"(created_at < ?) OR (created_at = ? AND id < ?)",
				cursorRun.CreatedAt, cursorRun.CreatedAt, *after,
			)
		}
	} else if before != nil {
		// Get the cursor item to find its created_at
		var cursorRun models.CheckRun
		if err := s.db.First(&cursorRun, "id = ?", *before).Error; err == nil {
			// Get items before the cursor (newer items in DESC order)
			// (created_at > cursor.created_at) OR (created_at = cursor.created_at AND id > cursor.id)
			query = query.Where(
				"(created_at > ?) OR (created_at = ? AND id > ?)",
				cursorRun.CreatedAt, cursorRun.CreatedAt, *before,
			)
		}
	}

	// Always sort from latest to oldest
	query = query.Order("created_at DESC, id DESC")

	if limit > 0 {
		// Fetch one extra to determine if there's a next page
		query = query.Limit(limit + 1)
	}

	if err := query.Find(&runs).Error; err != nil {
		return nil, err
	}

	return runs, nil
}

// UptimeDataPoint represents a single data point in the uptime chart
type UptimeDataPoint struct {
	Timestamp        time.Time `json:"timestamp"`
	UptimePercentage float64   `json:"uptime_percentage"`
	TotalRuns        int       `json:"total_runs"`
	Passing          int       `json:"passing"`
	Degraded         int       `json:"degraded"`
	Failing          int       `json:"failing"`
}

// GetCheckUptimeData returns aggregated uptime data for a check over a specified time range
// startTime and endTime define the time range (inclusive)
// timeBucket determines the aggregation interval: "minute", "hour", or "day"
func (s *Store) GetCheckUptimeData(checkID uuid.UUID, startTime, endTime time.Time, timeBucket string) ([]UptimeDataPoint, error) {
	// Validate time bucket
	if timeBucket != "minute" && timeBucket != "hour" && timeBucket != "day" {
		timeBucket = "hour" // Default to hour
	}

	// Query to aggregate check runs by time bucket
	type Result struct {
		Timestamp time.Time `gorm:"column:time_bucket"`
		Status    string    `gorm:"column:status"`
		Count     int       `gorm:"column:count"`
	}

	var results []Result
	err := s.db.Raw(`
		SELECT
			DATE_TRUNC(?, created_at) as time_bucket,
			status,
			COUNT(*) as count
		FROM check_runs
		WHERE check_id = ?
			AND created_at >= ?
			AND created_at <= ?
			AND deleted_at IS NULL
		GROUP BY time_bucket, status
		ORDER BY time_bucket ASC
	`, timeBucket, checkID, startTime, endTime).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Group results by time bucket and calculate uptime
	bucketMap := make(map[time.Time]*UptimeDataPoint)

	for _, result := range results {
		bucket := result.Timestamp.Truncate(getTruncateDuration(timeBucket))
		point, exists := bucketMap[bucket]
		if !exists {
			point = &UptimeDataPoint{
				Timestamp: bucket,
				TotalRuns: 0,
				Passing:   0,
				Degraded:  0,
				Failing:   0,
			}
			bucketMap[bucket] = point
		}

		point.TotalRuns += result.Count
		switch models.CheckRunStatus(result.Status) {
		case models.CheckRunStatusPassing:
			point.Passing += result.Count
		case models.CheckRunStatusDegraded:
			point.Degraded += result.Count
		case models.CheckRunStatusFailing:
			point.Failing += result.Count
		}
	}

	// Convert map to slice and calculate uptime percentage
	dataPoints := make([]UptimeDataPoint, 0, len(bucketMap))
	for _, point := range bucketMap {
		if point.TotalRuns > 0 {
			// Uptime = (passing + degraded) / total * 100
			// Degraded is considered "up" but not optimal
			point.UptimePercentage = float64(point.Passing+point.Degraded) / float64(point.TotalRuns) * 100.0
		} else {
			point.UptimePercentage = 0.0
		}
		dataPoints = append(dataPoints, *point)
	}

	// Sort by timestamp
	sort.Slice(dataPoints, func(i, j int) bool {
		return dataPoints[i].Timestamp.Before(dataPoints[j].Timestamp)
	})

	return dataPoints, nil
}

// getTruncateDuration returns the duration to truncate time buckets
func getTruncateDuration(bucket string) time.Duration {
	switch bucket {
	case "minute":
		return time.Minute
	case "hour":
		return time.Hour
	case "day":
		return 24 * time.Hour
	default:
		return time.Hour
	}
}

// DetermineTimeBucket automatically determines the appropriate time bucket based on the time range
// Returns "minute", "hour", or "day"
func DetermineTimeBucket(startTime, endTime time.Time) string {
	duration := endTime.Sub(startTime)

	// Less than 3 hours: use minute buckets
	if duration < 3*time.Hour {
		return "minute"
	}
	// Less than 7 days: use hour buckets
	if duration < 7*24*time.Hour {
		return "hour"
	}
	// 7 days or more: use day buckets
	return "day"
}

// TimingDataPoint represents a single timing data point from a check run
type TimingDataPoint struct {
	RunID         uuid.UUID              `json:"run_id"`
	Timestamp     time.Time              `json:"timestamp"`
	NetworkTimings map[string]interface{} `json:"network_timings"`
}

// GetCheckTimingsData returns timing data for all check runs within a specified time range
// startTime and endTime define the time range (inclusive)
// Returns a list of timing data points, one per check run
func (s *Store) GetCheckTimingsData(checkID uuid.UUID, startTime, endTime time.Time) ([]TimingDataPoint, error) {
	var runs []struct {
		ID            uuid.UUID      `gorm:"column:id"`
		CreatedAt     time.Time      `gorm:"column:created_at"`
		NetworkTimings datatypes.JSON `gorm:"column:network_timings;type:jsonb"`
	}

	err := s.db.Table("check_runs").
		Select("id, created_at, network_timings").
		Where("check_id = ?", checkID).
		Where("created_at >= ?", startTime).
		Where("created_at <= ?", endTime).
		Where("deleted_at IS NULL").
		Where("network_timings IS NOT NULL").
		Order("created_at ASC").
		Find(&runs).Error

	if err != nil {
		return nil, err
	}

	dataPoints := make([]TimingDataPoint, 0, len(runs))
	for _, run := range runs {
		// Unmarshal the JSONB data
		var timings map[string]interface{}
		if len(run.NetworkTimings) > 0 {
			if err := json.Unmarshal(run.NetworkTimings, &timings); err == nil && len(timings) > 0 {
				dataPoints = append(dataPoints, TimingDataPoint{
					RunID:         run.ID,
					Timestamp:     run.CreatedAt,
					NetworkTimings: timings,
				})
			}
		}
	}

	return dataPoints, nil
}
