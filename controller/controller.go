package controller

import (
	customer "gin-dbo/controller/customer"
	login "gin-dbo/controller/login"
	order "gin-dbo/controller/order"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Controller struct {
	Login    login.Usecase
	Customer customer.Usecase
	Order    order.Usecase
}

func Router(usecase *Controller, logger *logrus.Logger) *gin.Engine {
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin, X-Requested-With, Content-Type, Accept, Authorization, Access-Control-Allow-Headers, Accept-Encoding, X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	login.Router(router, usecase.Login, logger)
	customer.Router(router, usecase.Customer, logger)
	order.Router(router, usecase.Order, logger)
	return router
}
