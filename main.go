package main

import (
	"ars_projekat/handlers"
	"ars_projekat/model"
	"ars_projekat/repositories"
	"ars_projekat/services"
	"context"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	configInMemoryRepository := repositories.NewConfigInMemoryRepository()
	configService := services.NewConfigurationService(configInMemoryRepository)
	configHandler := handlers.NewConfigurationHandler(configService)

	configGroupInMemoryRepository := repositories.NewConfigGroupInMemoryRepository()
	configGroupService := services.NewConfigurationGroupService(configGroupInMemoryRepository)
	configGroupHandler := handlers.NewConfigurationGroupHandler(configGroupService, configService)

	testvr := model.Version{
		Major: 0,
		Minor: 0,
		Patch: 1,
	}

	testCfg1 := &model.Configuration{
		Name:       "TestKonfiguracija",
		Id:         203032,
		Version:    testvr,
		Parameters: make(map[string]string),
		Labels:     make(map[string]string),
	}

	testCfg2 := &model.Configuration{
		Name:       "TestKonfiguracija2",
		Id:         2030201323,
		Version:    testvr,
		Parameters: make(map[string]string),
		Labels:     make(map[string]string),
	}

	testCfg3 := &model.Configuration{
		Name: "TestKonfiguracija3",
		Id:   232312678,
		Version: model.Version{
			Major: 2,
			Minor: 0,
			Patch: 6,
		},
		Parameters: make(map[string]string),
		Labels:     make(map[string]string),
	}

	var group []model.Configuration
	group = append(group, *testCfg1)
	group = append(group, *testCfg2)
	group = append(group, *testCfg3)

	testGroup := model.ConfigurationGroup{
		Name: "TestGrupa",
		Id:   66564054,
		Version: model.Version{
			Major: 0,
			Minor: 0,
			Patch: 2,
		},
		Configurations: group,
	}
	configInMemoryRepository.Add(testCfg1)
	configInMemoryRepository.Add(testCfg2)
	configInMemoryRepository.Add(testCfg3)

	configGroupInMemoryRepository.Add(&testGroup)
	router := mux.NewRouter()
	// Config routes
	router.HandleFunc("/configs/{name}/{version}", configHandler.Get).Methods("GET")
	router.HandleFunc("/configs/", configHandler.Upsert).Methods("POST")
	router.HandleFunc("/configs/{name}/{version}", configHandler.Delete).Methods("DELETE")

	// Config group routes
	router.HandleFunc("/configs/groups/{name}/{version}", configGroupHandler.Get).Methods("GET")
	router.HandleFunc("/configs/groups/", configGroupHandler.Upsert).Methods("POST")
	router.HandleFunc("/configs/groups/{name}/{version}", configGroupHandler.Delete).Methods("DELETE")
	srv := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: router,
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		log.Println("Starting server..")
		if err := srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("Shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Stopped server")
}
