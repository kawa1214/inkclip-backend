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
var testDB *sql.DB

// const location = "Etc/UTC"

// func initLocale() {
// 	loc, err := time.LoadLocation(location)
// 	if err != nil {
// 		loc = time.FixedZone(location, 9*60*60)
// 	}
// 	time.Local = loc
// }

func TestMain(m *testing.M) {
	config, err := config.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQueries = New(testDB)

	code := m.Run()
	os.Exit(code)
}
