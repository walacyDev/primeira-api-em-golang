package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()

	router.SetTrustedProxies(nil)

	RegisterRoutes(router)

	router.Run(":3000")
}
