package handler

import (
	"go_kafka/docs"
	"go_kafka/pkg/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	docs.SwaggerInfo.Title = "Go Kafka API"
	docs.SwaggerInfo.Description = "Kafka ile çalışan hesap işlemleri API'si"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:2006"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http"}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/accounts")
	{
		api.POST("/", h.createAccount)
		api.GET("/", h.getAccountAll)
		api.GET("/:account_id", h.getAccountById)
		api.DELETE("/:account_id", h.deleteAccount)
	}

	trans := router.Group("/api")
	{
		trans.POST("/:account_id/deposit", h.deposit)
		trans.POST("/:account_id/withdraw", h.withdraw)
		trans.POST("/:account_id/transfer", h.transfer)
		trans.GET("/:account_id/transaction", h.getTransaction)
	}
	return router
}
