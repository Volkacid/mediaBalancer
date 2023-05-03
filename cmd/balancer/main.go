package main

import (
	"context"
	"github.com/Volkacid/mediaBalancer/internal/config"
	"github.com/Volkacid/mediaBalancer/internal/handlers"
	"github.com/Volkacid/mediaBalancer/internal/service"
	"github.com/Volkacid/mediaBalancer/internal/storage"
	"log"
	"net/http"
)

func main() {
	servPool := service.ConnectToCloud() //Мок взаимодействия с AWS
	ctx := context.Background()

	rdb, err := storage.NewRedisStorage(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go service.ServerStatsService(ctx, servPool, rdb) //Мониторит серверы и удаляет неиспользуемые

	conf := config.GetConfig()

	log.Print("Starting server at ", conf.ServerAddress)
	http.Handle("/rooms/create", handlers.CreateNewRoom(servPool, rdb))
	log.Fatal(http.ListenAndServe(conf.ServerAddress, nil))
}
