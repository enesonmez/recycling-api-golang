package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	. "oyeco-api/helpers"
	. "oyeco-api/helpers/error"

	//. "oyeco-api/models/address"
	//. "oyeco-api/models/request"
	. "oyeco-api/models/worker"

	"github.com/gorilla/mux"
)

// HTTP POST - /api/manageworkers/signin
func ManageWorkerSignInHandler(w http.ResponseWriter, r *http.Request) {
	var mngWorker Worker
	err := json.NewDecoder(r.Body).Decode(&mngWorker)                                                                        // Body içeriği User modeli ile eşleştirliyor.
	if value, data := JsonError(err, 500, "beklenmeyen json parse hatası (json verilerinizi kontrol edin)"); value == true { // Hata kontrolü
		Respond(w, 500, data)
		return
	}
	status, resp := mngWorker.ManagerSignIn() // resp değişkenine json verisi alınır.
	Respond(w, status, resp)                  // Respond fonksiyonu ile response yollanır.
}

// HTTP POST - /api/fieldworkers/signin
func FieldWorkerSignInHandler(w http.ResponseWriter, r *http.Request) {
	var fieldWorker Worker
	err := json.NewDecoder(r.Body).Decode(&fieldWorker)                                                                      // Body içeriği User modeli ile eşleştirliyor.
	if value, data := JsonError(err, 500, "beklenmeyen json parse hatası (json verilerinizi kontrol edin)"); value == true { // Hata kontrolü
		Respond(w, 500, data)
		return
	}
	status, resp := fieldWorker.FieldWorkerSignIn() // resp değişkenine json verisi alınır.
	Respond(w, status, resp)                        // Respond fonksiyonu ile response yollanır.
}

// HTTP POST - /api/fieldworkers/register
func FieldWorkerRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var field Worker
	field.SetRecordTime(time.Now()) // Worker değişkenlerine olması gereken değerler etleniyor.
	field.SetStatus(0)              // field worker

	err := json.NewDecoder(r.Body).Decode(&field)                                                                            // Body içeriği User modeli ile eşleştirliyor.
	if value, data := JsonError(err, 500, "beklenmeyen json parse hatası (json verilerinizi kontrol edin)"); value == true { // Hata kontrolü
		Respond(w, 500, data)
		return
	}

	status, resp := field.Create() // resp değişkenine json verisi alınır.
	Respond(w, status, resp)       // Respond fonksiyonu ile response yollanır.
}

// HTTP GET - /api/fieldworkers
func FieldWorkerAllGetHandler(w http.ResponseWriter, r *http.Request) {
	var mngWorker Worker
	status, resp := mngWorker.AllGet() // resp değişkenine json verisi alınır.
	Respond(w, status, resp)           // Respond fonksiyonu ile response yollanır.
}

// HTTP DELETE - /api/fieldworkers/{wID}
func FieldWorkerDeleteHandler(w http.ResponseWriter, r *http.Request) {
	var mngWorker Worker
	vars := mux.Vars(r)
	status, resp := mngWorker.Delete(vars["wID"]) // resp değişkenine json verisi alınır.
	Respond(w, status, resp)                      // Respond fonksiyonu ile response yollanır.
}

// HTTP PUT - /api/fieldworkers/{wID}
func FieldWorkerUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var mngWorker Worker
	vars := mux.Vars(r)
	err := json.NewDecoder(r.Body).Decode(&mngWorker)                                                                        // Body içeriği User modeli ile eşleştirliyor.
	if value, data := JsonError(err, 500, "beklenmeyen json parse hatası (json verilerinizi kontrol edin)"); value == true { // Hata kontrolü
		Respond(w, 500, data)
		return
	}
	status, resp := mngWorker.Update(vars["wID"]) // resp değişkenine json verisi alınır.
	Respond(w, status, resp)                      // Respond fonksiyonu ile response yollanır.
}
