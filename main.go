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
	// Environment variable set
	SetEnv() //helpers/generally.go
	// DB bağlantısı kurulup tablo kontrolleri yapılıyor.
	DBCreate() //helpers/generally.go

	// Port ayarlaması yapılıyor.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server (:%s) starting...", port)

	// Handle'lar belirleniyor.
	r := mux.NewRouter().StrictSlash(true) // "/api/users/register/" bu şekilde de çalışır.
	r.Handle("/api/users/register", IsAuthorized(UserRegisterHandler)).Methods("POST")
	r.Handle("/api/users/signin", IsAuthorized(UserSignInHandler)).Methods("POST")
	r.Handle("/api/users/update/{id}", IsAuthorized(UserUpdateHandler)).Methods("PUT")
	r.Handle("/api/users/updatePassword/{id}", IsAuthorized(UserUpdatePasswordHandler)).Methods("PUT")
	r.HandleFunc("/api/users/activation/{id}", ActivationHandler).Methods("GET")

	r.Handle("/api/users/{userID}/address", IsAuthorized(UserAddressRegisterHandler)).Methods("POST")
	r.Handle("/api/users/{userID}/address/{adrsID}", IsAuthorized(UserAddressUpdateHandler)).Methods("PUT")
	r.Handle("/api/users/{userID}/address", IsAuthorized(UserAddressGetHandler)).Methods("GET")
	r.Handle("/api/users/{userID}/address/{adrsID}", IsAuthorized(UserAddressDeleteHandler)).Methods("DELETE")

	r.Handle("/api/users/{userID}/requests", IsAuthorized(UserRequestRegisterHandler)).Methods("POST")
	r.Handle("/api/users/{userID}/requests", IsAuthorized(UserRequestGetHandler)).Methods("GET")
	r.Handle("/api/users/{userID}/requests/{reqID}", IsAuthorized(UserRequestDeleteHandler)).Methods("DELETE")

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	server.ListenAndServe()
	log.Println("Server ending...")
}
