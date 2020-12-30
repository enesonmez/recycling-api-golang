package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	. "oyeco-api/helpers"
	. "oyeco-api/helpers/error"
	. "oyeco-api/models/user"

	"github.com/gorilla/mux"
)

// HTTP POST - /api/users/register
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	user.SetRecordTime(time.Now()) // User değşkenlerine olması gereken değerler etleniyor.
	user.SetIsVerifyEmail(false)
	user.SetIsBlock(false)

	err := json.NewDecoder(r.Body).Decode(&user)                                                                               // Body içeriği User modeli ile eşleştirliyor.
	if value, data := JsonError(err, 500, "unexpected json parse error (check the json variable data types)"); value == true { // Hata kontrolü
		Respond(w, 500, data)
		return
	}

	status, resp := user.Create() // resp değişkenine json verisi alınır.
	Respond(w, status, resp)      // Respond fonksiyonu ile response yollanır.
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	vars := mux.Vars(r)
	status, resp := user.Activation(vars["id"]) // resp değişkenine json verisi alınır.
	Respond(w, status, resp)                    // Respond fonksiyonu ile
}
