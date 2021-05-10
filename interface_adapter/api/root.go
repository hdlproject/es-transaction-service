package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func ActivateAPI(port string) {
	router := gin.Default()
	apiRouter := router.Group("/api")

	RegisterTransactionAPI(apiRouter)

	router.Run(fmt.Sprintf(":%s", port))
}
