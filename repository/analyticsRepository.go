package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// ClickEvent represents a single click on a shortened URL
type ClickEvent struct {
	Timestamp time.Time `json:"timestamp"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Referer   string    `json:"referer"`
	Country   string    `json:"country,omitempty"`
	City      string    `json:"city,omitempty"`
}

// Analytics represents aggregated analytics for a shortened URL
type Analytics struct {
	ShortURL    string       `json:"short_url"`
	TotalClicks int64        `json:"total_clicks"`
	Clicks      []ClickEvent `json:"clicks,omitempty"`
}

// AnalyticsRepository defines the interface for analytics data storage
type AnalyticsRepository interface {
	RecordClick(shortURL string, event ClickEvent) error
	GetTotalClicks(shortURL string) (int64, error)
	GetClickHistory(shortURL string, limit int64) ([]ClickEvent, error)
}

type redisAnalyticsRepository struct {
	client *redis.Client
}

// NewRedisAnalyticsRepository creates a new Redis-based analytics repository
func NewRedisAnalyticsRepository(db *redis.Client) AnalyticsRepository {
	return &redisAnalyticsRepository{db}
}

func (repo *redisAnalyticsRepository) clickCountKey(shortURL string) string {
	return fmt.Sprintf("analytics:clicks:%s", shortURL)
}

func (repo *redisAnalyticsRepository) clickHistoryKey(shortURL string) string {
	return fmt.Sprintf("analytics:history:%s", shortURL)
}

// RecordClick stores a click event for a short URL
func (repo *redisAnalyticsRepository) RecordClick(shortURL string, event ClickEvent) error {
	ctx := context.Background()

	// Increment click counter
	err := repo.client.Incr(ctx, repo.clickCountKey(shortURL)).Err()
	if err != nil {
		return err
	}

	// Store click event in a sorted set with timestamp as score
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	score := float64(event.Timestamp.UnixNano())
	err = repo.client.ZAdd(ctx, repo.clickHistoryKey(shortURL), &redis.Z{
		Score:  score,
		Member: string(eventJSON),
	}).Err()

	return err
}

// GetTotalClicks returns the total number of clicks for a short URL
func (repo *redisAnalyticsRepository) GetTotalClicks(shortURL string) (int64, error) {
	ctx := context.Background()

	result, err := repo.client.Get(ctx, repo.clickCountKey(shortURL)).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	count, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetClickHistory returns the most recent click events for a short URL
func (repo *redisAnalyticsRepository) GetClickHistory(shortURL string, limit int64) ([]ClickEvent, error) {
	ctx := context.Background()

	// Get events in reverse order (most recent first)
	results, err := repo.client.ZRevRange(ctx, repo.clickHistoryKey(shortURL), 0, limit-1).Result()
	if err != nil {
		return nil, err
	}

	events := make([]ClickEvent, 0, len(results))
	for _, result := range results {
		var event ClickEvent
		if err := json.Unmarshal([]byte(result), &event); err != nil {
			continue // Skip malformed events
		}
		events = append(events, event)
	}

	return events, nil
}
