package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Banana-Boat/terryminal/main-service/internal/api"
	"github.com/Banana-Boat/terryminal/main-service/internal/db"
	"github.com/Banana-Boat/terryminal/main-service/internal/util"
	_ "github.com/go-sql-driver/mysql"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// 美化zerolog
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	/* 加载配置 */
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Error().Err(err).Msg("cannot load config")
		return
	}
	log.Info().Msg("configuration loaded successfully")

	/* 运行 DB migration */
	if err = runDBMigrate(config); err != nil {
		return
	}

	/* 连接数据库 */
	conn, err := sql.Open(
		config.DBDriver,
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName),
	)
	if err != nil {
		log.Error().Err(err).Msg("cannot connect to db")
		return
	}
	store := db.NewStore(conn)
	log.Info().Msg("db connected successfully")

	/* 运行 gin 服务 */
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Error().Err(err).Msg("cannot create server")
		return
	}

	log.Info().Msgf("gin server started at %s:%s successfully", config.MainServerHost, config.MainServerPort)
	err = server.Start(fmt.Sprintf("%s:%s", config.MainServerHost, config.MainServerPort))
	if err != nil {
		log.Error().Err(err).Msg("cannot start server")
	}
}

func runDBMigrate(config util.Config) error {
	m, err := migrate.New(
		config.MigrationFileUrl,
		fmt.Sprintf("%s://%s:%s@tcp(%s:%s)/%s", config.DBDriver, config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName),
	)
	if err != nil {
		log.Error().Err(err).Msg("cannot create migration instance")
		return err
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error().Err(err).Msg("failed to run migrate up")
		return err
	}
	log.Info().Msg("db migrated successfully")
	return nil
}
