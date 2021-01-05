package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	. "oyeco-api/helpers"
	. "oyeco-api/helpers/error"
	. "oyeco-api/models/address"
	. "oyeco-api/models/user"

	"github.com/gorilla/mux"
)

// HTTP POST - /api/users/register
func UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	user.SetRecordTime(time.Now()) // User değşkenlerine olması gereken değerler etleniyor.
	user.SetIsVerifyEmail(false)
	user.SetIsBlock(false)

	err := json.NewDecoder(r.Body).Decode(&user)                                                                             // Body içeriği User modeli ile eşleştirliyor.
	if value, data := JsonError(err, 500, "beklenmeyen json parse hatası (json verilerinizi kontrol edin)"); value == true { // Hata kontrolü
		Respond(w, 500, data)
		return
	}

	status, resp := user.Create() // resp değişkenine json verisi alınır.
	Respond(w, status, resp)      // Respond fonksiyonu ile response yollanır.
}

// HTTP GET - /api/users/activation/{id}
func ActivationHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	vars := mux.Vars(r)
	status, resp := user.Activation(vars["id"]) // resp değişkenine json verisi alınır.
	Respond(w, status, resp)                    // Respond fonksiyonu ile
}

// HTTP POST - /api/users/signin
func UserSignInHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)                                                                             // Body içeriği User modeli ile eşleştirliyor.
	if value, data := JsonError(err, 500, "beklenmeyen json parse hatası (json verilerinizi kontrol edin)"); value == true { // Hata kontrolü
		Respond(w, 500, data)
		return
	}
	status, resp := user.SignIn() // resp değişkenine json verisi alınır.
	Respond(w, status, resp)      // Respond fonksiyonu ile response yollanır.
}

// HTTP PUT - /api/users/update/{id}
func UserUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	vars := mux.Vars(r)
	err := json.NewDecoder(r.Body).Decode(&user)                                                                             // Body içeriği User modeli ile eşleştirliyor.
	if value, data := JsonError(err, 500, "beklenmeyen json parse hatası (json verilerinizi kontrol edin)"); value == true { // Hata kontrolü
		Respond(w, 500, data)
		return
	}
	status, resp := user.Update(vars["id"]) // resp değişkenine json verisi alınır.
	Respond(w, status, resp)                // Respond fonksiyonu ile response yollanır.
}

// HTTP PUT - /api/users/updatePassword/{id}
func UserUpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	vars := mux.Vars(r)
	err := json.NewDecoder(r.Body).Decode(&user)
	if value, data := JsonError(err, 500, "beklenmeyen json parse hatası (json verilerinizi kontrol edin)"); value == true { // Hata kontrolü
		Respond(w, 500, data)
		return
	}
	status, resp := user.UpdatePassword(vars["id"]) // resp değişkenine json verisi alınır.
	Respond(w, status, resp)                        // Respond fonksiyonu ile response yollanır.
}

// HTTP POST - /api/users/address/{id}
func UserAddressRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var address Address
	vars := mux.Vars(r)

	err := json.NewDecoder(r.Body).Decode(&address)                                                                          // Body içeriği User modeli ile eşleştirliyor.
	if value, data := JsonError(err, 500, "beklenmeyen json parse hatası (json verilerinizi kontrol edin)"); value == true { // Hata kontrolü
		Respond(w, 500, data)
		return
	}

	status, resp := address.Create(vars["id"]) // resp değişkenine json verisi alınır.
	Respond(w, status, resp)                   // Respond fonksiyonu ile response yollanır.
}

// HTTP PUT - /api/users/address/{id}
func UserAddressUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var address Address
	vars := mux.Vars(r)

	err := json.NewDecoder(r.Body).Decode(&address)                                                                          // Body içeriği User modeli ile eşleştirliyor.
	if value, data := JsonError(err, 500, "beklenmeyen json parse hatası (json verilerinizi kontrol edin)"); value == true { // Hata kontrolü
		Respond(w, 500, data)
		return
	}

	status, resp := address.Update(vars["id"]) // resp değişkenine json verisi alınır.
	Respond(w, status, resp)                   // Respond fonksiyonu ile response yollanır.
}

// HTTP GET - /api/users/address/{id}
func UserAddressGetHandler(w http.ResponseWriter, r *http.Request) {
	var address Address
	vars := mux.Vars(r)

	status, resp := address.Get(vars["id"]) // resp değişkenine json verisi alınır.
	Respond(w, status, resp)                // Respond fonksiyonu ile response yollanır.
}

// HTTP DELETE - /api/users/address/{id}
func UserAddressDeleteHandler(w http.ResponseWriter, r *http.Request) {
	var address Address
	vars := mux.Vars(r)

	err := json.NewDecoder(r.Body).Decode(&address)                                                                          // Body içeriği User modeli ile eşleştirliyor.
	if value, data := JsonError(err, 500, "beklenmeyen json parse hatası (json verilerinizi kontrol edin)"); value == true { // Hata kontrolü
		Respond(w, 500, data)
		return
	}

	status, resp := address.Delete(vars["id"]) // resp değişkenine json verisi alınır.
	Respond(w, status, resp)                   // Respond fonksiyonu ile response yollanır.
}
