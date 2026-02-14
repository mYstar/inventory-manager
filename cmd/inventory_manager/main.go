package main

import (
	"inventory_manager/internal/api"
)

func main() {
	router := api.SetupRoutes()

	router.Run("0.0.0.0:80")
}
