package main

import (
	"wallet-transfer-assignment/config"
	pgdb "wallet-transfer-assignment/db"
	db "wallet-transfer-assignment/db/generated"
	"wallet-transfer-assignment/internal/handlers"
	"wallet-transfer-assignment/internal/repo"
	"wallet-transfer-assignment/internal/usecase"
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
	pgdbServ := pgdb.PGDBService(cfg.PostgresConfig)
	pgdbInstance, err := pgdbServ.GetDBInstance()
	if err != nil {
		log.Fatal().Err(err).Msg("DB instantiation failed")
	}

	//logger := log.Logger
	//v1 := app.Group("/api/v1")
	dbstore := db.NewStore(pgdbInstance)
	walletRepo := repo.NewWalletRepository(dbstore)
	walletsvc := usecase.NewWalletService(walletRepo)
	transferRepo := repo.NewTransferRepository(dbstore)
	ledgerRepo := repo.NewLedgerEntryRepository(dbstore)
	transfersvc := usecase.NewTransferService(transferRepo, walletRepo, ledgerRepo, dbstore)
	transferHandler := handlers.NewTransferHandler(transfersvc)
	walletHandler := handlers.NewWalletHandler(walletsvc)
	healthHandler := &handlers.HealthCheck{}
	routerMap := handlers.RouterMap{
		"/health":    healthHandler,
		"/wallet/":   walletHandler,
		"/transfer/": transferHandler,
	}

	server := server.GetServer(cfg.APIServerConfig, app, routerMap)

	waitforShutdownInterrupt := server.StartServer()
	<-waitforShutdownInterrupt
}
