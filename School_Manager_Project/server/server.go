package main

import (
	"crypto/tls"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"restapi/internal/api/middlewares"
	"restapi/internal/api/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		return
	}

	port := os.Getenv("API_PORT")

	cert := "cert.pem"
	key := "key.pem"

	mux := router.Router()
	//jwtRouter := middlewares.JWTMiddleware(mux)
	jwtMiddleware := middlewares.MiddlewaresExcludePaths(middlewares.JWTMiddleware, "/execs/login")
	handler := jwtMiddleware(mux)
	tlfConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:      port,
		Handler:   handler,
		TLSConfig: tlfConfig,
	}

	fmt.Println("Server is running on port:", port)
	err2 := server.ListenAndServeTLS(cert, key)
	if err2 != nil {
		log.Fatal("Error starting the server:", err2)
	}
}
