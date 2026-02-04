package services

import (
	"time"

	"github.com/mfvitale/curto/repository"
)

// AnalyticsService provides methods for tracking and retrieving URL analytics
type AnalyticsService struct {
	analyticsRepo repository.AnalyticsRepository
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(analyticsRepo repository.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{analyticsRepo: analyticsRepo}
}

// RecordClick records a click event for a short URL
func (s *AnalyticsService) RecordClick(shortURL, ip, userAgent, referer string) error {
	event := repository.ClickEvent{
		Timestamp: time.Now().UTC(),
		IP:        ip,
		UserAgent: userAgent,
		Referer:   referer,
	}

	return s.analyticsRepo.RecordClick(shortURL, event)
}

// GetAnalytics returns aggregated analytics for a short URL
func (s *AnalyticsService) GetAnalytics(shortURL string) (*repository.Analytics, error) {
	totalClicks, err := s.analyticsRepo.GetTotalClicks(shortURL)
	if err != nil {
		return nil, err
	}

	return &repository.Analytics{
		ShortURL:    shortURL,
		TotalClicks: totalClicks,
	}, nil
}

// GetAnalyticsWithHistory returns analytics with click history for a short URL
func (s *AnalyticsService) GetAnalyticsWithHistory(shortURL string, limit int64) (*repository.Analytics, error) {
	totalClicks, err := s.analyticsRepo.GetTotalClicks(shortURL)
	if err != nil {
		return nil, err
	}

	clicks, err := s.analyticsRepo.GetClickHistory(shortURL, limit)
	if err != nil {
		return nil, err
	}

	return &repository.Analytics{
		ShortURL:    shortURL,
		TotalClicks: totalClicks,
		Clicks:      clicks,
	}, nil
}
