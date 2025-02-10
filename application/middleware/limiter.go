package middleware

import (
	"challenge-goexpert-rate-limiter/config/application/usecases"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RateLimiterConfig struct {
	Token          string
	Requests       int
	LimiterUseCase usecases.LimiterUseCaseInterface
}

const (
	errorMessage string = "you have reached the maximum number of requests or actions allowed within a certain time frame"
)

func RateLimiterMiddleware(config RateLimiterConfig) fiber.Handler {
	start := time.Now()

	return func(c *fiber.Ctx) error {

		clientIP := c.IP()
		parameter := clientIP
		limit := 10

		headers := c.GetReqHeaders()

		if headers["Api_key"] != nil {
			token := headers["Api_key"]
			if !strings.EqualFold(token[0], "") && strings.EqualFold(token[0], config.Token) {
				log.Printf("Token informado: %s", token)
				parameter = token[0]
				limit = config.Requests
				log.Printf("Parameter: %s | Limit: %d", parameter, limit)
			}
		}

		err := config.LimiterUseCase.ValidRateLimiter(parameter, limit)
		if err != nil && err.Error() == errorMessage {

			go func(ip string) {
				time.Sleep(time.Minute)
				config.LimiterUseCase.RemoveBlock(ip)
			}(clientIP)

			log.Printf("IP %s com requisições excedidas:  %d.", parameter, limit)
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": errorMessage,
			})
		}

		log.Printf("[%s] %s | Status: %d | Request Time: %s",
			c.Method(),
			c.Path(),
			c.Response().StatusCode(),
			time.Since(start),
		)

		return c.Next()
	}
}
