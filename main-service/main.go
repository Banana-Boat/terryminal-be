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
	"github.com/Banana-Boat/terryminal/main-service/internal/worker"
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

	/* 运行 TaskProcessor */
	go runTaskProcessor(config)

	/* 运行 http 服务 */
	runHttpServer(config)
}

func runTaskProcessor(config util.Config) {

	taskProcessor := worker.NewTaskProcessor(config)

	log.Info().Msgf("task processor is running at %s:%s", config.RedisHost, config.RedisPort)
	if err := taskProcessor.Start(); err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}

func runHttpServer(config util.Config) {
	/* 执行DB migration */
	migration, err := migrate.New(
		config.MigrationFileUrl,
		fmt.Sprintf("%s://%s:%s@tcp(%s:%s)/%s", config.DBDriver, config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to migrate db")
		return
	}
	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error().Err(err).Msg("failed to migrate db")
		return
	}
	log.Info().Msg("db migrated successfully")

	/* 连接数据库 */
	conn, err := sql.Open(
		config.DBDriver,
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to db")
		return
	}
	store := db.NewStore(conn)
	log.Info().Msg("db connected successfully")

	/* 创建并启动Server */
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Error().Err(err).Msg("cannot create server")
		return
	}

	log.Info().Msgf("http server is running at %s:%s", config.MainServerHost, config.MainServerPort)
	if err := server.Start(); err != nil {
		log.Error().Err(err).Msg("failed to start server")
	}
}
