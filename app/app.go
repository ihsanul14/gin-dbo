package app

import (
	"gin-dbo/framework/database"
	"log"
	"os"

	"gin-dbo/framework/logger"

	"github.com/subosito/gotenv"

	controller "gin-dbo/controller"
	customerController "gin-dbo/controller/customer"
	loginController "gin-dbo/controller/login"
	orderController "gin-dbo/controller/order"

	_ "gin-dbo/docs"
)

func Run() {
	err := gotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	var baseLogger = logger.Logger()
	dbConn, err := database.ConnectSQL(baseLogger)
	if err != nil {
		baseLogger.Fatal(err)
	}

	customerRepository := customerController.NewRepository(dbConn)
	customerUsecase := customerController.NewUsecase(customerRepository)

	loginRepository := loginController.NewRepository(dbConn)
	loginUsecase := loginController.NewUsecase(loginRepository, customerRepository)

	orderRepository := orderController.NewRepository(dbConn)
	orderUsecase := orderController.NewUsecase(orderRepository, customerRepository)

	httpRouter := &controller.Controller{
		Login:    loginUsecase,
		Customer: customerUsecase,
		Order:    orderUsecase,
	}

	router := controller.Router(httpRouter, baseLogger)
	if err = router.Run(":" + os.Getenv("PORT")); err != nil {
		baseLogger.Fatal(err)
	}
}
