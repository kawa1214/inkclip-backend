package main

import (
	"database/sql"
	"log"

	"github.com/bookmark-manager/bookmark-manager/api"
	"github.com/bookmark-manager/bookmark-manager/config"
	db "github.com/bookmark-manager/bookmark-manager/db/sqlc"
	_ "github.com/lib/pq"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
