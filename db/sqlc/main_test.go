package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/bookmark-manager/bookmark-manager/config"
	_ "github.com/lib/pq"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	config, err := config.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQueries = New(conn)

	code := m.Run()
	os.Exit(code)
}
