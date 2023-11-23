package database

import (
	"os"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	customer "gin-dbo/model/customer"
	login "gin-dbo/model/login"
	order "gin-dbo/model/order"
)

var Db *gorm.DB

func ConnectSQL(log *logrus.Logger) (*gorm.DB, error) {
	mysqlDialector := os.Getenv("MYSQL_DIALECTOR")
	Db, err := gorm.Open(mysql.Open(mysqlDialector), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err = Db.AutoMigrate(&login.User{}, &customer.Customer{}, &order.Order{}); err != nil {
		return nil, err
	}

	return Db, err
}
