package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Banana-Boat/terryminal/gateway-service/internal/db"
	"github.com/Banana-Boat/terryminal/gateway-service/internal/http"
	"github.com/Banana-Boat/terryminal/gateway-service/internal/util"
	_ "github.com/go-sql-driver/mysql"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// 美化 zerolog
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	/* 加载配置 */
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Error().Err(err).Msg("cannot load config")
		return
	}

	/* 运行 DB migration */
	if err = runDBMigrate(config); err != nil {
		log.Error().Err(err).Msg("failed to migrate db")
		return
	}
	log.Info().Msg("db migrated successfully")

	/* 获取数据库连接 */
	store, err := getDBStore(config)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to db")
		return
	}
	log.Info().Msg("db connected successfully")

	/* 运行 http 服务 */
	runHttpServer(config, store)
}

func runDBMigrate(config util.Config) error {
	m, err := migrate.New(
		config.MigrationFileUrl,
		fmt.Sprintf("%s://%s:%s@tcp(%s:%s)/%s", config.DBDriver, config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName),
	)
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func getDBStore(config util.Config) (*db.Store, error) {
	conn, err := sql.Open(
		config.DBDriver,
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName),
	)
	if err != nil {
		return nil, err
	}
	store := db.NewStore(conn)
	return store, nil
}

func runHttpServer(config util.Config, store *db.Store) {
	server, err := http.NewServer(config, store)
	if err != nil {
		log.Error().Err(err).Msg("cannot create server")
		return
	}

	if err := server.Start(); err != nil {
		log.Error().Err(err).Msg("failed to start server")
	}
}
