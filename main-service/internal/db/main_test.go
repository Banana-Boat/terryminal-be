package db

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/rs/zerolog/log"

	"github.com/Banana-Boat/go-micro-template/main-service/internal/util"
	_ "github.com/go-sql-driver/mysql"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
		return
	}

	migration, err := migrate.New(
		config.MigrationFileUrl,
		fmt.Sprintf("%s://%s:%s@tcp(%s:%s)/%s", config.DBDriver, config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create migration instance")
		return
	}
	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
		return
	}

	conn, err := sql.Open(
		config.DBDriver,
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("can't connect to db")
		return
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
