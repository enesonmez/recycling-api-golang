package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	. "oyeco-api/helpers"
	. "oyeco-api/helpers/error"
	. "oyeco-api/models/address"
	. "oyeco-api/models/request"
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

// HTTP PUT - /api/users/updateBlock/{id}
func UserUpdateBlockHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	vars := mux.Vars(r)
	err := json.NewDecoder(r.Body).Decode(&user)
	if value, data := JsonError(err, 500, "beklenmeyen json parse hatası (json verilerinizi kontrol edin)"); value == true { // Hata kontrolü
		Respond(w, 500, data)
		return
	}
	status, resp := user.UpdateBlock(vars["id"]) // resp değişkenine json verisi alınır.
	Respond(w, status, resp)                     // Respond fonksiyonu ile response yollanır.
}

// HTTP GET - /api/users
func UserAllGetHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	vars := mux.Vars(r)

	status, resp := user.Get(vars["userID"]) // resp değişkenine json verisi alınır.
	Respond(w, status, resp)                 // Respond fonksiyonu ile response yollanır.
}

// HTTP POST - /api/users/{userID}/address
func UserAddressRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var address Address
	vars := mux.Vars(r)

	err := json.NewDecoder(r.Body).Decode(&address)                                                                          // Body içeriği User modeli ile eşleştirliyor.
	if value, data := JsonError(err, 500, "beklenmeyen json parse hatası (json verilerinizi kontrol edin)"); value == true { // Hata kontrolü
		Respond(w, 500, data)
		return
	}

	status, resp := address.Create(vars["userID"]) // resp değişkenine json verisi alınır.
	Respond(w, status, resp)                       // Respond fonksiyonu ile response yollanır.
}

// HTTP PUT - /api/users/{userID}/address/{adrsID}
func UserAddressUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var address Address
	vars := mux.Vars(r)

	err := json.NewDecoder(r.Body).Decode(&address)                                                                          // Body içeriği User modeli ile eşleştirliyor.
	if value, data := JsonError(err, 500, "beklenmeyen json parse hatası (json verilerinizi kontrol edin)"); value == true { // Hata kontrolü
		Respond(w, 500, data)
		return
	}

	status, resp := address.Update(vars["userID"], vars["adrsID"]) // resp değişkenine json verisi alınır.
	Respond(w, status, resp)                                       // Respond fonksiyonu ile response yollanır.
}

// HTTP GET - /api/users/{userID}/address
func UserAddressGetHandler(w http.ResponseWriter, r *http.Request) {
	var address Address
	vars := mux.Vars(r)

	status, resp := address.Get(vars["userID"]) // resp değişkenine json verisi alınır.
	Respond(w, status, resp)                    // Respond fonksiyonu ile response yollanır.
}

// HTTP DELETE - /api/users/{userID}/address/{adrsID}
func UserAddressDeleteHandler(w http.ResponseWriter, r *http.Request) {
	var address Address
	vars := mux.Vars(r)

	status, resp := address.Delete(vars["userID"], vars["adrsID"]) // resp değişkenine json verisi alınır.
	Respond(w, status, resp)                                       // Respond fonksiyonu ile response yollanır.
}

// HTTP POST - /api/users/{userID}/requests
func UserRequestRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var request Request
	vars := mux.Vars(r)
	request.SetRequestCreateTime(time.Now())
	request.SetState(1) // İstek bekleniyor.

	err := json.NewDecoder(r.Body).Decode(&request)                                                                          // Body içeriği User modeli ile eşleştirliyor.
	if value, data := JsonError(err, 500, "beklenmeyen json parse hatası (json verilerinizi kontrol edin)"); value == true { // Hata kontrolü
		Respond(w, 500, data)
		return
	}

	status, resp := request.Create(vars["userID"]) // resp değişkenine json verisi alınır.
	Respond(w, status, resp)                       // Respond fonksiyonu ile response yollanır.
}

// HTTP GET - /api/users/{userID}/requests
func UserRequestGetHandler(w http.ResponseWriter, r *http.Request) {
	var request Request
	vars := mux.Vars(r)

	status, resp := request.Get(vars["userID"]) // resp değişkenine json verisi alınır.
	Respond(w, status, resp)                    // Respond fonksiyonu ile response yollanır.
}

// HTTP DELETE - /api/users/{userID}/requests/{reqID}
func UserRequestDeleteHandler(w http.ResponseWriter, r *http.Request) {
	var request Request
	vars := mux.Vars(r)

	status, resp := request.Delete(vars["userID"], vars["reqID"]) // resp değişkenine json verisi alınır.
	Respond(w, status, resp)                                      // Respond fonksiyonu ile response yollanır.
}
