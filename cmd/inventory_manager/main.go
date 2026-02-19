package main

import (
	"inventory_manager/internal/api"
	"inventory_manager/internal/storage"
	"os"

	"github.com/gin-gonic/gin"
)

// main initializes the application and starts the HTTP server.
func main() {
	dbFile := os.Getenv("DB_FILE")
	database, err := storage.NewSqlitePersistence(dbFile)
	if err != nil {
		panic(err)
	}

	inventory := api.NewInventory(database)
	router := gin.Default()
	api.SetupRoutes(router, inventory)

	err = router.Run("0.0.0.0:8080")
	if err != nil {
		panic(err)
	}
}
