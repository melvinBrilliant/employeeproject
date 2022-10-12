package main

import (
	"net/http"
	"os"

	"com.melvin.employee/src/router"
)

func main() {
	e := router.New()
	os.Setenv("APP_NAME", "EmployeeApp")
	confAppName := os.Getenv("APP_NAME")

	os.Setenv("SERVER_PORT", "8080")
	confServerPort := os.Getenv("SERVER_PORT")

	server := new(http.Server)
	server.Addr = ":" + confServerPort
	e.Logger.Print("Starting", confAppName)
	e.Logger.Fatal(e.StartServer(server))
}