package ARSProjekat

import (
	"ars_projekat/repositories"
	"ars_projekat/services"
)

func main() {

	repo := repositories.NewConfigInMemoryRepository()

	service := services.NewConfigurationService(repo)

	service.Hello()
}
