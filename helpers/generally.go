package helpers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	. "oyeco-api/db"
	. "oyeco-api/models/config"
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

func SetEnv() {
	conf := new(Config) // Konfigürasyon dosyasındaki değişkenleri kullanmak için nesne kullanılmıştır.
	conf.ConfigRead()
	os.Setenv("DBHost", conf.DBHost)
	os.Setenv("DBPort", strconv.Itoa(conf.DBPort))
	os.Setenv("DBUsername", conf.DBUsername)
	os.Setenv("DBPassword", conf.DBPassword)
	os.Setenv("DBName", conf.DBName)
	os.Setenv("Email", conf.Email)
	os.Setenv("EmailPassword", conf.EmailPassword)
	os.Setenv("BaseURL", conf.BaseURL)
}
