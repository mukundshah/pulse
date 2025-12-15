package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"pulse/internal/config"
	"pulse/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
	ctx    context.Context
}

func Connect(cfg *config.Config) (*Client, error) {
	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		// Try direct connection if URL parsing fails
		opt = &redis.Options{
			Addr: cfg.RedisURL,
		}
	}

	client := redis.NewClient(opt)
	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Connected to Redis")
	return &Client{
		client: client,
		ctx:    ctx,
	}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

// HealthCheck pings Redis to verify connectivity
func (c *Client) HealthCheck() error {
	return c.client.Ping(c.ctx).Err()
}

func (c *Client) EnqueueJob(checkID uuid.UUID) error {
	jobData, err := json.Marshal(checkID.String())
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	if err := c.client.LPush(c.ctx, "pulse:jobs", jobData).Err(); err != nil {
		return fmt.Errorf("failed to enqueue job: %w", err)
	}

	return nil
}

func (c *Client) DequeueJob(timeout time.Duration) (uuid.UUID, error) {
	result, err := c.client.BRPop(c.ctx, timeout, "pulse:jobs").Result()
	if err != nil {
		if err == redis.Nil {
			return uuid.Nil, fmt.Errorf("no job available")
		}
		return uuid.Nil, fmt.Errorf("failed to dequeue job: %w", err)
	}

	if len(result) < 2 {
		return uuid.Nil, fmt.Errorf("invalid job format")
	}

	var checkIDStr string
	if err := json.Unmarshal([]byte(result[1]), &checkIDStr); err != nil {
		return uuid.Nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}

	checkID, err := uuid.Parse(checkIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse check ID: %w", err)
	}

	return checkID, nil
}

// GetQueueDepth returns the number of jobs in the queue
func (c *Client) GetQueueDepth() (int64, error) {
	return c.client.LLen(c.ctx, "pulse:jobs").Result()
}

func (c *Client) GetCheck(checkID uuid.UUID) (*models.Check, error) {
	key := fmt.Sprintf("pulse:check:%s", checkID.String())
	val, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get check from cache: %w", err)
	}

	var check models.Check
	if err := json.Unmarshal([]byte(val), &check); err != nil {
		return nil, fmt.Errorf("failed to unmarshal check: %w", err)
	}

	return &check, nil
}

func (c *Client) SetCheck(check *models.Check, ttl time.Duration) error {
	key := fmt.Sprintf("pulse:check:%s", check.ID.String())
	val, err := json.Marshal(check)
	if err != nil {
		return fmt.Errorf("failed to marshal check: %w", err)
	}

	if err := c.client.Set(c.ctx, key, val, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set check in cache: %w", err)
	}

	return nil
}

func (c *Client) DeleteCheck(checkID uuid.UUID) error {
	key := fmt.Sprintf("pulse:check:%s", checkID.String())
	if err := c.client.Del(c.ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete check from cache: %w", err)
	}
	return nil
}

// SetSession stores a session indicator in Redis cache
func (c *Client) SetSession(jti string, ttl time.Duration) error {
	key := fmt.Sprintf("pulse:session:%s", jti)
	return c.client.Set(c.ctx, key, "active", ttl).Err()
}

// GetSession retrieves a session indicator from Redis cache
func (c *Client) GetSession(jti string) (string, error) {
	key := fmt.Sprintf("pulse:session:%s", jti)
	return c.client.Get(c.ctx, key).Result()
}

// DeleteSession removes a session indicator from Redis cache
func (c *Client) DeleteSession(jti string) error {
	key := fmt.Sprintf("pulse:session:%s", jti)
	return c.client.Del(c.ctx, key).Err()
}
