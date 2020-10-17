package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/buffup/GolangTechTask/cmd/server/internal/handlers"
)

func main() {
	port := "8080"
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		port = portEnv
	}

	// Replace with DB impl
	store := handlers.Storage()

	routes := handlers.Routes(store)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), routes); err != nil {
		panic(err)
	}
}
