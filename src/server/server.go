package server

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"wallet-transfer-assignment/config"
	"wallet-transfer-assignment/internal/handlers"
	"wallet-transfer-assignment/internal/middlewares"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

type Server struct {
	Config   *config.ServerConfigAPIServer
	FiberApp *fiber.App
	Handlers handlers.RouterMap
}

func (s *Server) Setup() *Server {
	if s.FiberApp == nil {
		log.Fatal().Msg("Server setup incorrectly!")
	}
	middlewares.SetupMiddlewares(s.FiberApp)

	for path, apiFunction := range s.Handlers {
		apigroup := s.FiberApp.Group(path)
		apiFunction.RegisterRoutes(apigroup)
	}
	return s
}

func GetServer(Config *config.ServerConfigAPIServer, FiberApp *fiber.App, Handlers handlers.RouterMap) *Server {
	serv := &Server{
		Config,
		FiberApp,
		Handlers,
	}
	serv.Setup()
	return serv
}

func (s *Server) StartServer() <-chan os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := s.FiberApp.Listen(s.Config.Port, //+strconv.Itoa(s.Config.Port),
			fiber.ListenConfig{
				GracefulContext: nil, // TODO: Implement graceful shutdown.
				EnablePrefork:   false,
			})
		if err != nil {
			log.Fatal().Err(err)
		}
	}()
	return quit
}

func GetFiberApplication() *fiber.App {
	config := fiber.Config{
		AppName:         "Wallet Transferer",
		BodyLimit:       4 * 4 * 1024,
		ReadTimeout:     2 * time.Second,
		JSONEncoder:     sonic.Marshal,
		JSONDecoder:     sonic.Unmarshal,
		StructValidator: nil,
	}
	return fiber.New(config)
}
