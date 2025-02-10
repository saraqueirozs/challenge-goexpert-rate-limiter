package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type RateLimiterController struct{}

func NewRateLimiterController() *RateLimiterController {
	return &RateLimiterController{}
}

func (c *RateLimiterController) GetController(ctx *fiber.Ctx) error {
	log.Printf("Processo encerrado com sucesso.")

	return ctx.Status(200).JSON("Requisição realizada com sucesso.")
}
