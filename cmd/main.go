package main

import (
	"fmt"
	"log"

	"github.com/SaidovZohid/competition-project/config"
	"github.com/SaidovZohid/competition-project/storage"
	"github.com/redis/go-redis/v9"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load(".")

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

	_ = storage.NewStoragePg(psqlConn)
	_ = storage.NewInMemoryStorage(rdb)
}
