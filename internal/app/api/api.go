package api

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	recoverer "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	slogfiber "github.com/samber/slog-fiber"
)

type Config struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type Service interface {
}

type API struct {
	cfg Config
	svc Service

	r         *fiber.App
	validator *validator.Validate
}

func New(cfg Config, svc Service) *API {
	r := fiber.New()

	r.Use(slogfiber.New(slog.Default()))
	r.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	r.Use(requestid.New())
	r.Use(recoverer.New())

	return &API{
		cfg: cfg,
		svc: svc,
		r:   r,
		validator: validator.New(
			validator.WithRequiredStructEnabled(),
		),
	}
}

func (a *API) Run() error {
	return a.r.Listen(a.cfg.Host + ":" + a.cfg.Port)
}

func (a *API) Shutdown() error {
	return a.r.Shutdown()
}
