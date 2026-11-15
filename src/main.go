package main

import (
	"wallet-transfer-assignment/config"
	database "wallet-transfer-assignment/db"
	store "wallet-transfer-assignment/db/generated"
	"wallet-transfer-assignment/internal/handlers"
	"wallet-transfer-assignment/internal/handlers/wallethlr"
	"wallet-transfer-assignment/server"

	"github.com/rs/zerolog/log"
)

func main() {
	app := server.GetFiberApplication()

	config.ConfigureZeroLog()
	cfg, err := config.GetServerConfig("./cfg.json")
	if err != nil {
		log.Fatal().Err(err).Msg("Server configuration retrieval fail.")
	}
	pgdbServ := database.PGDBService(cfg.PostgresConfig)
	db, err := pgdbServ.GetDBInstance()
	if err != nil {
		log.Fatal().Err(err).Msg("DB instantiation failed")
	}

	//logger := log.Logger
	// v1 := app.Group("/api/v1")
	dbstore := store.NewStore(db)
	walletHandler := wallethlr.NewWalletHandler(wall, &log.Logger)
	routerMap := handlers.RouterMap{
		"/wallet/":   walletHandler,
		"/transfer/": transferHandler,
	}

	server := server.GetServer(cfg.APIServerConfig, app, routerMap)

	waitforShutdownInterrupt := server.StartServer()
	<-waitforShutdownInterrupt
}
