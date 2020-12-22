package helpers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	. "oyeco-api/db"
)

// String değer boş mu?
func IsEmpty(value string) bool {
	if len(value) == 0 {
		return true
	}
	return false
}

// Response yollamak için dinamik fonksiyon
func Respond(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json") // Header değerleri ve response body'si
	w.WriteHeader(status)
	w.Write(data)
}

// DB bağlantısı kurulup tablo kontrolleri yapılıyor.
func DBCreate() {
	a := new(Db)
	db, _ := a.Connect()
	err := a.CreateTables(db)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	db.Close()
}

var mySigningKey = "captainjacksparrowsayshi"

// api authorized function
func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("api-key") != "" {

			if r.Header.Get("api-key") == mySigningKey {
				endpoint(w, r)
			}
		} else {

			fmt.Fprintf(w, "Not Authorized")
		}
	})
}
