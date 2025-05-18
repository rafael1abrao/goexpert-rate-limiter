// limiter_service_test.go
package limiter

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStrategy simula a camada de persistência (Redis, etc)
type MockStrategy struct {
	mock.Mock
}

func (m *MockStrategy) IsAllowed(ctx context.Context, key string, limit int, blockDurationSeconds int) (bool, error) {
	args := m.Called(ctx, key, limit, blockDurationSeconds)
	return args.Bool(0), args.Error(1)
}

func setupLimiter(rateIP, rateToken, block int, mock RateLimiterStrategy) *LimiterService {
	return &LimiterService{
		strategy:             mock,
		rateLimitIP:          rateIP,
		rateLimitToken:       rateToken,
		blockDurationSeconds: block,
	}
}

func TestRateLimiter_AllowByIP(t *testing.T) {
	ctx := context.Background()
	mockStrat := new(MockStrategy)
	limiterSvc := setupLimiter(5, 10, 300, mockStrat)

	mockStrat.On("IsAllowed", ctx, "ip:127.0.0.1", 5, 300).Return(true, nil)

	allowed, err := limiterSvc.IsRequestAllowed(ctx, "127.0.0.1", "")
	assert.NoError(t, err)
	assert.True(t, allowed)
	mockStrat.AssertExpectations(t)
}

func TestRateLimiter_BlockByToken(t *testing.T) {
	ctx := context.Background()
	mockStrat := new(MockStrategy)
	limiterSvc := setupLimiter(5, 10, 300, mockStrat)

	token := "abc123"
	mockStrat.On("IsAllowed", ctx, "token:"+token, 10, 300).Return(false, nil)

	allowed, err := limiterSvc.IsRequestAllowed(ctx, "192.168.0.1", token)
	assert.NoError(t, err)
	assert.False(t, allowed)
	mockStrat.AssertExpectations(t)
}

func TestRateLimiter_ErrorFromStrategy(t *testing.T) {
	ctx := context.Background()
	mockStrat := new(MockStrategy)
	limiterSvc := setupLimiter(5, 10, 300, mockStrat)

	mockStrat.On("IsAllowed", ctx, "ip:192.168.0.2", 5, 300).Return(false, errors.New("redis down"))

	allowed, err := limiterSvc.IsRequestAllowed(ctx, "192.168.0.2", "")
	assert.Error(t, err)
	assert.False(t, allowed)
	mockStrat.AssertExpectations(t)
}
