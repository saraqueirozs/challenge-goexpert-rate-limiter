package usecases

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRedisRepository struct {
	mock.Mock
}

func (m *MockRedisRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockRedisRepository) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisRepository) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockRedisRepository) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisRepository) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestLimiterUseCase_ValidRateLimiter(t *testing.T) {
	tests := []struct {
		name          string
		parameter     string
		limit         int
		existsReturn  bool
		existsError   error
		getReturn     string
		getError      error
		setError      error
		expectedError error
	}{
		{
			name:          "Blocked parameter",
			parameter:     "user1",
			limit:         5,
			existsReturn:  true,
			existsError:   nil,
			expectedError: errors.New(errorMessage),
		},
		{
			name:          "Limit reached",
			parameter:     "user2",
			limit:         5,
			existsReturn:  false,
			existsError:   nil,
			getReturn:     "5",
			getError:      nil,
			expectedError: errors.New(errorMessage),
		},
		{
			name:          "Limit not reached",
			parameter:     "user3",
			limit:         5,
			existsReturn:  false,
			existsError:   nil,
			getReturn:     "3",
			getError:      nil,
			setError:      nil,
			expectedError: nil,
		},
		{
			name:          "Error on Set",
			parameter:     "user4",
			limit:         5,
			existsReturn:  false,
			existsError:   nil,
			getReturn:     "4",
			getError:      nil,
			setError:      errors.New("set error"),
			expectedError: errors.New("set error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRedisRepository)
			uc := NewLimiterUseCase(mockRepo)

			blockKey := fmt.Sprintf("%s:block", tt.parameter)
			mockRepo.On("Exists", mock.Anything, blockKey).Return(tt.existsReturn, tt.existsError)

			if !tt.existsReturn {
				mockRepo.On("Get", mock.Anything, tt.parameter).Return(tt.getReturn, tt.getError)
			}

			if tt.getError == nil && !tt.existsReturn {
				if tt.getReturn != "" {
					quantidade, _ := strconv.Atoi(tt.getReturn)
					if quantidade >= tt.limit {
						mockRepo.On("Set", mock.Anything, blockKey, true, time.Minute).Return(tt.setError)
					} else {
						mockRepo.On("Set", mock.Anything, tt.parameter, quantidade+1, time.Second).Return(tt.setError)
					}
				}
			}

			err := uc.ValidRateLimiter(tt.parameter, tt.limit)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestLimiterUseCase_RemoveBlock(t *testing.T) {
	mockRepo := new(MockRedisRepository)
	uc := NewLimiterUseCase(mockRepo)

	parameter := "user1"
	blockKey := fmt.Sprintf("%s:block", parameter)

	mockRepo.On("Delete", mock.Anything, blockKey).Return(nil)

	uc.RemoveBlock(parameter)

	mockRepo.AssertCalled(t, "Delete", mock.Anything, blockKey)
}
