package error

import (
	"encoding/json"
	"log"
	"net/http"

	. "../../models/info"
	. "../../models/page"
)

func CheckEror(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

func JsonError(err error, w http.ResponseWriter, status int, message string) bool { // true dönerse handlers return edecek
	if err != nil {
		info := new(Info)
		info.InfoConstructer(false, status, message) // Info nesnesine değerler setleniyor.
		var infoPage InfoPage
		infoPage.InfoPageConstructer(info)
		data, err := json.Marshal(infoPage)
		CheckEror(err)

		w.Header().Set("Content-Type", "application/json") // Client'a respos olarak hata mesajı döndürülür.
		w.WriteHeader(status)
		http.Error(w, string(data), status)
		return true
	}
	return false
}
