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
	r.Handle("/api/users/updateBlock/{id}", IsAuthorized(UserUpdateBlockHandler)).Methods("PUT")
	r.HandleFunc("/api/users/activation/{id}", ActivationHandler).Methods("GET")
	r.Handle("/api/users", IsAuthorized(UserAllGetHandler)).Methods("GET")
	r.Handle("/api/users/{userID}", IsAuthorized(UserAllGetHandler)).Methods("GET")

	r.Handle("/api/users/{userID}/address", IsAuthorized(UserAddressRegisterHandler)).Methods("POST")
	r.Handle("/api/users/{userID}/address/{adrsID}", IsAuthorized(UserAddressUpdateHandler)).Methods("PUT")
	r.Handle("/api/users/{userID}/address", IsAuthorized(UserAddressGetHandler)).Methods("GET")
	r.Handle("/api/users/{userID}/address/{adrsID}", IsAuthorized(UserAddressDeleteHandler)).Methods("DELETE")

	r.Handle("/api/users/{userID}/requests", IsAuthorized(UserRequestRegisterHandler)).Methods("POST")
	r.Handle("/api/users/{userID}/requests", IsAuthorized(UserRequestGetHandler)).Methods("GET")
	r.Handle("/api/users/{userID}/requests/{reqID}", IsAuthorized(UserRequestDeleteHandler)).Methods("DELETE")
	r.Handle("/api/requests", IsAuthorized(UserRequestAllGetHandler)).Methods("GET")

	r.Handle("/api/manageworkers/signin", IsAuthorized(ManageWorkerSignInHandler)).Methods("POST")
	r.Handle("/api/fieldworkers/signin", IsAuthorized(FieldWorkerSignInHandler)).Methods("POST")

	r.Handle("/api/fieldworkers/register", IsAuthorized(FieldWorkerRegisterHandler)).Methods("POST")
	r.Handle("/api/fieldworkers", IsAuthorized(FieldWorkerAllGetHandler)).Methods("GET")
	r.Handle("/api/fieldworkers/{wID}", IsAuthorized(FieldWorkerDeleteHandler)).Methods("DELETE")
	r.Handle("/api/fieldworkers/{wID}", IsAuthorized(FieldWorkerUpdateHandler)).Methods("PUT")

	r.Handle("/api/routes/register", IsAuthorized(RouteRegisterHandler)).Methods("POST")

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	server.ListenAndServe()
	log.Println("Server ending...")
}
