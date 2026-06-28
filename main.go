package main

import (
	"UserService/handlers"
	"UserService/internal/config"
	"UserService/internal/repository"
	"UserService/internal/service"
	db2 "UserService/pkg/db"
	"UserService/pkg/logger"
	"UserService/pkg/memory"
	"fmt"
	"log"
	"net/http"
)

func main() {
	lg, err := logger.New(true)
	if err != nil {
		log.Fatal("failed to create logger", err)
	}

	cfg, err := config.New("config/config.env")
	if err != nil {
		log.Fatal("config.New", err)
	}

	db, err := db2.New(db2.Options{
		Host:     cfg.DBHOST,
		Port:     cfg.DBPORT,
		User:     cfg.DBUSER,
		Password: cfg.DBPASSWORD,
		DBName:   cfg.DBName,
	})
	if err != nil {
		log.Fatal("failed to connect to db %w", err)
		return
	}
	defer db.Close()

	myCache := memory.NewMemoryCache()

	repoUser := repository.NewUserRepo(db)
	repoOrder := repository.NewOrderRepo(db)
	repoAdmin := repository.NewAdminRepo(db)

	serviceUser := service.NewUserService(repoUser, repoOrder, myCache)
	serviceOrder := service.NewOrderService(repoOrder)
	serviceAdmin := service.NewAdminService(repoAdmin)

	handlerUser := handlers.NewHandler(*lg, serviceUser, serviceAdmin, serviceOrder)

	mw := handlers.NewMiddleware(repoUser)

	router := handlers.New(handlerUser, mw)

	fmt.Println("service is started")

	err = http.ListenAndServe(cfg.HttpPort, router)
	if err != nil {
		log.Fatal(err)
	}
}
