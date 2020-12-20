package main

import (
	"log"
	"net/http"
	"os"

	. "../db"
	. "../handlers"
	"github.com/gorilla/mux"
)

func main() {
	a := new(Db)
	a.Assign()
	db, _ := a.Connect()
	err := a.CreateTables(db)
	if err != nil {
		log.Fatal(err)
	}

	//b.UserConsturcter("Enes", "Sönmez", "05350607409", "son@hotmail.com", "123456789", "male", time.Time{}, time.Now(), false, false)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server starting...")

	r := mux.NewRouter().StrictSlash(true) // "/api/users/register/" bu şekilde de çalışır.
	r.HandleFunc("/api/users/register", RegisterHandler).Methods("POST")

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	server.ListenAndServe()
	log.Println("Server ending...")
}
