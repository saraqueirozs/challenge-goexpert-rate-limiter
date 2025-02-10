package usecases

import (
	"challenge-goexpert-rate-limiter/config/application/repository"
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

type LimiterUseCaseInterface interface {
	ValidRateLimiter(parameter string, limit int) error
	RemoveBlock(parameter string)
}

type limiterUseCase struct {
	redisRepository repository.RedisRepositoryInterface
}

func NewLimiterUseCase(repository repository.RedisRepositoryInterface) LimiterUseCaseInterface {
	return &limiterUseCase{
		redisRepository: repository,
	}
}

const (
	errorMessage string = "you have reached the maximum number of requests or actions allowed within a certain time frame"
)

func (l *limiterUseCase) ValidRateLimiter(parameter string, limit int) error {
	repository := l.redisRepository

	log.Printf("Parameter: %s", parameter)

	// Verifica se a chave de bloqueio existe
	blockKey := fmt.Sprintf("%s:block", parameter)
	blocked, _ := repository.Exists(context.Background(), blockKey)
	if blocked {
		return errors.New(errorMessage)
	}

	resp, _ := repository.Get(context.Background(), parameter)
	log.Printf("Quantidade atual: %s", resp)

	quantidade, _ := strconv.Atoi(resp)

	if quantidade >= limit {
		repository.Set(context.Background(), blockKey, true, time.Minute)
		return errors.New(errorMessage)
	}

	err := repository.Set(context.Background(), parameter, quantidade+1, time.Second)
	if err != nil {
		return err
	}

	return nil
}

func (l *limiterUseCase) RemoveBlock(parameter string) {
	l.redisRepository.Delete(context.Background(), fmt.Sprintf("%s:block", parameter))
}
