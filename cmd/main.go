package main

import (
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	server "go_back/src"
	"go_back/src/handler"
	"go_back/src/repository"
	"go_back/src/service"
	"log"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatal(err.Error())
	}

	fileDB, err := repository.NewMysqlDB(viper.GetString("db_link"))
	if err != nil {
		log.Fatal(err.Error())
	}

	var cacheDB *redis.Client = repository.NewRedisDB(&repository.RedisConfig{
		Host:     "localhost",
		Port:     0,
		Password: "",
		DBNum:    1,
	})

	var repository_ *repository.Repository = repository.NewRepository(fileDB, cacheDB)
	var services *service.Service = service.NewService(repository_)
	var handlers *handler.Handler = handler.NewHandler(services)

	srv := new(server.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatal(err.Error())
	}
}

func initConfig() error {
	viper.SetConfigFile("configs/conf.yml")
	return viper.ReadInConfig()
}
