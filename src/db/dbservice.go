package database

import (
	"context"
	"fmt"
	"time"
	"wallet-transfer-assignment/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBService struct {
	Config     *config.ServerConfigDB
	dbinstance *pgxpool.Pool
}

func configureDB(config *config.ServerConfigDB) (*pgxpool.Config, error) {
	pgConfig, err := pgxpool.ParseConfig(config.DSN)
	if err != nil {
		return nil, err
	}
	// TODO: Remove hardcoding
	pgConfig.MaxConns = 10
	pgConfig.MinConns = 2
	pgConfig.MaxConnLifetime = time.Hour
	pgConfig.MaxConnIdleTime = 30 * time.Minute
	return pgConfig, nil
}

func (dbs *DBService) createNewDBInstance() (*pgxpool.Pool, error) {
	config, err := configureDB(dbs.Config)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return pool, err
}

func (dbs *DBService) GetDBInstance() (*pgxpool.Pool, error) {
	if dbs.dbinstance != nil {
		return dbs.dbinstance, nil
	}

	dbinst, err := dbs.createNewDBInstance()
	if err != nil {
		return nil, err
	}
	dbs.dbinstance = dbinst
	return dbs.dbinstance, nil
}

func (dbs *DBService) ClearDBInstance() (bool, error) {
	if dbs.dbinstance == nil {
		return true, nil
	}

	dbs.dbinstance.Close()
	return true, nil
}

func PGDBService(pgconfig *config.ServerConfigDB) *DBService {
	return &DBService{
		Config:     pgconfig,
		dbinstance: nil,
	}
}
