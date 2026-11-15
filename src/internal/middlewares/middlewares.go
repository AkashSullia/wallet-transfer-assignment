package middlewares

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func SetupMiddlewares(app *fiber.App) {
	app.Get("/", func(c fiber.Ctx) error {
		c.Accepts("application/json")
		c.AcceptsEncodings("compress", "br")
		return nil
	})
	app.Use(
		cors.New(),
		compress.New(compress.Config{
			Next: func(c fiber.Ctx) bool {
				return c.Get("no-compress") == "yes"
			},
		}),
		recoverer.New(recoverer.Config{
			EnableStackTrace: true,
		}),
		requestid.New(),
		logger.New(logger.Config{
			Format: "${time} ${pid} ${locals:requestid} ${status} - ${method} ${path} [${latency}] \n${reqHeaders} ${body}\n",
		}))
}
