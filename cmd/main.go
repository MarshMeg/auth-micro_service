package main

import (
	server "github.com/DikosAs/GoAuthApi.git/src"
	"github.com/DikosAs/GoAuthApi.git/src/handler"
	"github.com/DikosAs/GoAuthApi.git/src/repository"
	"github.com/DikosAs/GoAuthApi.git/src/service"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"log"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatal(err.Error())
	}

	fileDB, err := repository.NewMysqlClient(viper.GetString("db_link"))
	if err != nil {
		log.Fatal(err.Error())
	}
	rediska := repository.NewRedisClient(&redis.Options{
		Addr: viper.GetString("redis_ka"),
	})

	var repository_ *repository.Repository = repository.NewRepository(fileDB, rediska)
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
