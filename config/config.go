package config

import (
	"challenge-goexpert-rate-limiter/config/application/controllers"
	"challenge-goexpert-rate-limiter/config/application/middleware"
	"challenge-goexpert-rate-limiter/config/application/repository"
	"challenge-goexpert-rate-limiter/config/application/usecases"
	"fmt"

	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type Configure struct {
	apiKey string `mapstructure:"API_KEY"`
}

func Initialize() {

	// Configuração das variáveis de ambiente
	cfg, err := LoadConfig(".")
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	cfg.apiKey = viper.GetString("API_KEY")
	split := strings.Split(cfg.apiKey, ":")

	token := split[0]
	request, _ := strconv.Atoi(split[1])

	redisRepository := repository.NewRedisRepository()
	limiterUseCase := usecases.NewLimiterUseCase(redisRepository)

	// Configuração do middleware
	rateLimiterConfig := middleware.RateLimiterConfig{
		Token:          token,
		Requests:       request,
		LimiterUseCase: limiterUseCase,
	}

	// Middleware
	app.Use(middleware.RateLimiterMiddleware(rateLimiterConfig))

	setRoutes(app)

	// Iniciando o servidor
	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}
}

func LoadConfig(path string) (*Configure, error) {
	var cfg *Configure
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath("./")
	viper.SetConfigFile("config.env")
	viper.AutomaticEnv()

	fmt.Println("Loading config from path:", path)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}

func setRoutes(app *fiber.App) {

	rateLimiterController := controllers.NewRateLimiterController()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "API is running",
		})
	})

	app.Get("/", rateLimiterController.GetController)
}
