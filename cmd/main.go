package main

import (
	"github.com/DikosAs/GoAuthApi.git/src/handler"
	"github.com/DikosAs/GoAuthApi.git/src/router"
	"github.com/DikosAs/GoAuthApi.git/src/server"
	"github.com/DikosAs/GoAuthApi.git/src/storage"
	"github.com/DikosAs/GoAuthApi.git/src/storage/controllers"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"log"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatal(err.Error())
	}

	fileDB, err := controllers.NewMysqlClient(viper.GetString("db_link"))
	if err != nil {
		log.Fatal(err.Error())
	}
	appRedis := controllers.NewRedisClient(&redis.Options{
		Addr: viper.GetString("redis_ka"),
	})

	var appStorage *storage.Storage = storage.NewRepository(fileDB, appRedis)
	var appHandler *handler.Handler = handler.NewHandler(appStorage)
	var appRouter *router.Router = router.NewRouter(appHandler)

	srv := new(server.Server)
	if err := srv.Run(viper.GetString("app.port"), appRouter.InitRoutes(viper.GetString("app.mode"))); err != nil {
		log.Fatal(err.Error())
	}
}

func initConfig() error {
	viper.SetConfigFile("configs/conf.yml")
	return viper.ReadInConfig()
}
