package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	DB = initDB("tasks.db")
	defer DB.Close()

	RegisterRoutes(router)

	router.Run(":3000")
}
