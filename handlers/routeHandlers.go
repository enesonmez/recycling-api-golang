package handlers

import (
	"encoding/json"
	"net/http"

	. "oyeco-api/helpers"
	. "oyeco-api/helpers/error"

	//. "oyeco-api/models/address"
	//. "oyeco-api/models/request"
	. "oyeco-api/models/route"

	"github.com/gorilla/mux"
)

// HTTP POST - /api/routes/register
func RouteRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var route Route
	err := json.NewDecoder(r.Body).Decode(&route)                                                                            // Body içeriği User modeli ile eşleştirliyor.
	if value, data := JsonError(err, 500, "beklenmeyen json parse hatası (json verilerinizi kontrol edin)"); value == true { // Hata kontrolü
		Respond(w, 500, data)
		return
	}
	status, resp := route.Create() // resp değişkenine json verisi alınır.
	Respond(w, status, resp)       // Respond fonksiyonu ile response yollanır.
}

// HTTP POST - /api/fieldworkers/{fwID}/routes
func RouteGetHandler(w http.ResponseWriter, r *http.Request) {
	var route Route
	vars := mux.Vars(r)

	status, resp := route.Get(vars["fwID"]) // resp değişkenine json verisi alınır.
	Respond(w, status, resp)                // Respond fonksiyonu ile response yollanır.
}

// HTTP POST - /api/routes/{routeID}
func RouteAddressGetHandler(w http.ResponseWriter, r *http.Request) {
	var route Route
	vars := mux.Vars(r)

	status, resp := route.GetRouteAddress(vars["routeID"]) // resp değişkenine json verisi alınır.
	Respond(w, status, resp)                               // Respond fonksiyonu ile response yollanır.
}
