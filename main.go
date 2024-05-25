// Post API
//
//	Title: ARS Projekat API
//
//	Schemes: http
//	Version: 0.0.1
//	BasePath: /
//
//	Produces:
//	  - application/json
//
// swagger:meta
package main

import (
	"ars_projekat/handlers"
	"ars_projekat/middleware"
	"ars_projekat/repositories"
	"ars_projekat/services"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	swaggerMiddleware "github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

func main() {

	logger := log.New(os.Stdout, "[config-api]", log.LstdFlags)

	store, err := repositories.New(logger)
	if err != nil {
		logger.Fatal(err)
	}

	configService := services.NewConfigurationService(*store)
	configHandler := handlers.NewConfigurationHandler(configService)

	configGroupService := services.NewConfigurationGroupService(*store)
	configGroupHandler := handlers.NewConfigurationGroupHandler(configGroupService)

	idempotencyService := services.NewIdempotencyService(*store)

	limiter := middleware.NewRateLimiter(time.Second, 3)
	idempotencyMiddleware := middleware.NewIdempotency(&idempotencyService)

	router := mux.NewRouter()
	router.Use(func(next http.Handler) http.Handler {
		return middleware.AdaptHandler(next, limiter)
	})
	router.Use(func(next http.Handler) http.Handler {
		return middleware.AdaptIdempotencyHandler(next, idempotencyMiddleware)
	})

	// Config routes
	router.HandleFunc("/configs/{name}/{version}", configHandler.Get).Methods("GET")
	router.HandleFunc("/configs/", configHandler.Upsert).Methods("POST")
	router.HandleFunc("/configs/{name}/{version}", configHandler.Delete).Methods("DELETE")

	// Config group routes
	router.HandleFunc("/groups/{name}/{version}/{labels: ?.*}", configGroupHandler.Get).Methods("GET")
	router.HandleFunc("/groups/", configGroupHandler.Upsert).Methods("POST")
	router.HandleFunc("/groups/{name}/{version}/{labels: ?.*}", configGroupHandler.Delete).Methods("DELETE")
	router.HandleFunc("/groups/{name}/{version}", configGroupHandler.AddConfig).Methods("PUT")

	// Serve the swagger.yaml file
	router.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "swagger.yaml")
	}).Methods("GET")

	// SwaggerUI
	optionsDevelopers := swaggerMiddleware.SwaggerUIOpts{SpecURL: "swagger.yaml"}
	developerDocumentationHandler := swaggerMiddleware.SwaggerUI(optionsDevelopers, nil)
	router.Handle("/docs", developerDocumentationHandler)

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

	/* testing rate limiter
	<-time.After(2 * time.Second)

	fmt.Println("Starting rate limiter test...")
	url := "http://0.0.0.0:8000/configs/TestKonfiguracija/0.0.1"
	numRequests := 10
	var wg sync.WaitGroup
	wg.Add(numRequests)
	for i := 0; i < numRequests; i++ {
		go sendRequest(url, &wg)
	}
	wg.Wait() */

	<-quit

	log.Println("Shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Stopped server")
}

/* testing rate limiter
func sendRequest(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Response Status:", resp.Status)
} */
