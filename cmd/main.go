package main

import (
	"context"
	"github.com/NawafSwe/media-scout-service/cmd/config"
	"github.com/NawafSwe/media-scout-service/cmd/mediascout"
	"github.com/NawafSwe/media-scout-service/pkg/db"
	"log"
)

func main() {
	cfg, err := config.NewConfig(".", ".env")
	if err != nil {
		log.Fatalf("err loading config, err: %v", err)
	}
	dbConn, err := db.NewDBConn(cfg.DB)
	if err != nil {
		log.Fatalf("err creating db conn, err: %v", err)
	}
	if err := mediascout.RunHTTPServer(context.Background(), dbConn, cfg); err != nil {
		log.Fatalf("failed to run http server: %v", err)
	}
}
