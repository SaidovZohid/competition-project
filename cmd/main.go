package main

import (
	"fmt"

	"github.com/SaidovZohid/competition-project/api"
	"github.com/SaidovZohid/competition-project/config"
	"github.com/SaidovZohid/competition-project/pkg/logger"
	"github.com/SaidovZohid/competition-project/pkg/token"
	"github.com/SaidovZohid/competition-project/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load(".")

	logger.Init()
	log := logger.GetLogger()

	psqlUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)

	psqlConn, err := sqlx.Connect("postgres", psqlUrl)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	strg := storage.NewStoragePg(psqlConn)
	inMemory := storage.NewInMemoryStorage(rdb)

	tMaker, err := token.NewJWTMaker(cfg.AuthSecretKey)
	if err != nil {
		log.WithError(err).Fatal("error while making token JWT maker")
	}
	api := api.New(&api.RouterOptions{
		Cfg:        &cfg,
		Storage:    strg,
		InMemory:   inMemory,
		TokenMaker: tMaker,
		Logger:     &log,
	})

	if err := api.Run(cfg.HttpPort); err != nil {
		log.WithError(err).Fatal("error while running server")
	}
}
