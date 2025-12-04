package store

import (
	"pulse/internal/redis"

	"gorm.io/gorm"
)

type Store struct {
	db    *gorm.DB
	redis *redis.Client
}

func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

func NewWithCache(db *gorm.DB, redisClient *redis.Client) *Store {
	return &Store{db: db, redis: redisClient}
}
