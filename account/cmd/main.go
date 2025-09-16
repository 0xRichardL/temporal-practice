package main

import (
	"log"

	"github.com/0xRichardL/temporal-practice/account/internal/models"
	"github.com/0xRichardL/temporal-practice/account/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	/* DBs: */
	db, err := gorm.Open(postgres.Open(""), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Account{})
	/* Temporal */
	/* Services */
	services.NewAccountService(db)

	r := gin.Default()

	log.Println("Starting server on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalln("Unable to start server", err)
	}
}
