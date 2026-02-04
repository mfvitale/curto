package services

import (
	"testing"
	"time"

	"github.com/mfvitale/curto/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type AnalyticsRepositoryMock struct {
	mock.Mock
}

func (m *AnalyticsRepositoryMock) RecordClick(shortURL string, event repository.ClickEvent) error {
	args := m.Called(shortURL, event)
	return args.Error(0)
}

func (m *AnalyticsRepositoryMock) GetTotalClicks(shortURL string) (int64, error) {
	args := m.Called(shortURL)
	return args.Get(0).(int64), args.Error(1)
}

func (m *AnalyticsRepositoryMock) GetClickHistory(shortURL string, limit int64) ([]repository.ClickEvent, error) {
	args := m.Called(shortURL, limit)
	return args.Get(0).([]repository.ClickEvent), args.Error(1)
}

func TestRecordClick(t *testing.T) {
	assert := assert.New(t)

	analyticsRepoMock := new(AnalyticsRepositoryMock)
	analyticsRepoMock.On("RecordClick", "abc123", mock.AnythingOfType("repository.ClickEvent")).Return(nil)

	analyticsService := NewAnalyticsService(analyticsRepoMock)

	err := analyticsService.RecordClick("abc123", "192.168.1.1", "Mozilla/5.0", "https://google.com")

	assert.Nil(err)
	analyticsRepoMock.AssertExpectations(t)
	analyticsRepoMock.AssertNumberOfCalls(t, "RecordClick", 1)
}

func TestGetAnalytics(t *testing.T) {
	assert := assert.New(t)

	analyticsRepoMock := new(AnalyticsRepositoryMock)
	analyticsRepoMock.On("GetTotalClicks", "abc123").Return(int64(42), nil)

	analyticsService := NewAnalyticsService(analyticsRepoMock)

	analytics, err := analyticsService.GetAnalytics("abc123")

	assert.Nil(err)
	assert.Equal("abc123", analytics.ShortURL)
	assert.Equal(int64(42), analytics.TotalClicks)
	analyticsRepoMock.AssertExpectations(t)
}

func TestGetAnalyticsWithHistory(t *testing.T) {
	assert := assert.New(t)

	clickEvents := []repository.ClickEvent{
		{
			Timestamp: time.Now().UTC(),
			IP:        "192.168.1.1",
			UserAgent: "Mozilla/5.0",
			Referer:   "https://google.com",
		},
		{
			Timestamp: time.Now().UTC(),
			IP:        "192.168.1.2",
			UserAgent: "Chrome/100",
			Referer:   "https://twitter.com",
		},
	}

	analyticsRepoMock := new(AnalyticsRepositoryMock)
	analyticsRepoMock.On("GetTotalClicks", "abc123").Return(int64(2), nil)
	analyticsRepoMock.On("GetClickHistory", "abc123", int64(10)).Return(clickEvents, nil)

	analyticsService := NewAnalyticsService(analyticsRepoMock)

	analytics, err := analyticsService.GetAnalyticsWithHistory("abc123", 10)

	assert.Nil(err)
	assert.Equal("abc123", analytics.ShortURL)
	assert.Equal(int64(2), analytics.TotalClicks)
	assert.Len(analytics.Clicks, 2)
	assert.Equal("192.168.1.1", analytics.Clicks[0].IP)
	analyticsRepoMock.AssertExpectations(t)
}

func TestGetAnalyticsZeroClicks(t *testing.T) {
	assert := assert.New(t)

	analyticsRepoMock := new(AnalyticsRepositoryMock)
	analyticsRepoMock.On("GetTotalClicks", "newurl").Return(int64(0), nil)

	analyticsService := NewAnalyticsService(analyticsRepoMock)

	analytics, err := analyticsService.GetAnalytics("newurl")

	assert.Nil(err)
	assert.Equal("newurl", analytics.ShortURL)
	assert.Equal(int64(0), analytics.TotalClicks)
	analyticsRepoMock.AssertExpectations(t)
}
