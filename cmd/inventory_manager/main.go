package main

import (
	"inventory_manager/internal/api"
	"inventory_manager/internal/db"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	dbFile := os.Getenv("DB_FILE")
	database, err := db.NewSqlitePersistence(dbFile)
	if err != nil {
		panic(err)
	}
	inventory := api.NewInventory(database)
	router := gin.Default()
	api.SetupRoutes(router, inventory)

	err = router.Run("0.0.0.0:80")
	if err != nil {
		panic(err)
	}
}
