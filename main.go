package ARSProjekat

import (
	"ars_projekat/handlers"
	"ars_projekat/repositories"
	"ars_projekat/services"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	repo := repositories.NewConfigInMemoryRepository()
	service := services.NewConfigurationService(repo)
	handler := handlers.NewConfigurationHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/configs/{name}/{version}", handler.Get).Methods("GET")
	http.ListenAndServe("0.0.0.0:8000", router)
}
