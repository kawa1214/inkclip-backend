package main

import (
	"database/sql"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/inkclip/backend/api"
	"github.com/inkclip/backend/config"
	db "github.com/inkclip/backend/db/sqlc"
	"github.com/inkclip/backend/mail"

	_ "github.com/lib/pq"
)

// @securityDefinitions.apikey AccessToken
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

	mailClient := mail.NewMailClient(config)

	server, err := api.NewServer(config, store, mailClient)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
