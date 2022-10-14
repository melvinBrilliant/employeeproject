package main

import (
	"log"

	"com.melvin.employee/config"
	"com.melvin.employee/router"
	"github.com/joho/godotenv"
)

func main() {
	var err error
	_echo := router.Router()

	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	if err = config.ConnectDB(); err != nil {
		log.Fatal(err)
	}

	_echo.Logger.Fatal(_echo.Start(":8099"))
}
