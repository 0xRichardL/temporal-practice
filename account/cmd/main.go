package main

import (
	"github.com/0xRichardL/temporal-practice/account/internal/models"
	"github.com/0xRichardL/temporal-practice/account/internal/rest"
	"github.com/0xRichardL/temporal-practice/account/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// DBs
	db, err := gorm.Open(postgres.Open(""), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Account{})
	// Services
	accountService := services.NewAccountService(db)
	// REST
	accountControler := rest.NewAccountController(accountService)
	accountControler.RegisterRoutes()
	r := gin.Default()
	r.Run()
}
