package main

import (
	"github.com/personal-project/pitch-league/database"
	"github.com/personal-project/pitch-league/router"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	log.Println("bismillah")
	db := database.DB()

	app := fiber.New()
	router.Setup(app, db, router.Config{
		JWTSecret:              "secret",
		AccessTokenExpireTime:  15,
		RefreshTokenExpireTime: 24,
	})

	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("%v", err)
	}

	log.Println("hoşcagal")
}
