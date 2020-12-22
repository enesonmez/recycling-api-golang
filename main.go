package main

import (
	"log"
	"net/http"
	"os"

	. "oyeco-api/handlers"

	. "oyeco-api/helpers"

	"github.com/gorilla/mux"
)

func main() {
	DBCreate() // DB bağlantısı kurulup tablo kontrolleri yapılıyor.

	// Port ayarlaması yapılıyor.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server (:%s) starting...", port)

	// Handle'lar belirleniyor.
	r := mux.NewRouter().StrictSlash(true) // "/api/users/register/" bu şekilde de çalışır.
	r.Handle("/api/users/register", IsAuthorized(RegisterHandler)).Methods("POST")

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	server.ListenAndServe()
	log.Println("Server ending...")
}
