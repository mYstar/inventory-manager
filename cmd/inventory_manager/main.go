package main

import (
	"inventory_manager/internal/api"
	"inventory_manager/internal/db"

	"github.com/gin-gonic/gin"
)

func main() {
	inventory := api.NewInventory(db.NewMemoryPersistence())
	router := gin.Default()
	api.SetupRoutes(router, inventory)

	err := router.Run("0.0.0.0:80")
	if err != nil {
		panic(err)
	}
}
