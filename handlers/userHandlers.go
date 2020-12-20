package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	. "../db"
	. "../helpers/error"
	. "../models/info"
	. "../models/page"
	. "../models/user"
)

// HTTP POST - /api/users/register
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	user.SetRecordTime(time.Now()) // User değşkenlerine olması gereken değerler etleniyor.
	user.SetIsVerifyEmail(false)
	user.SetIsBlock(false)

	err := json.NewDecoder(r.Body).Decode(&user)                                         // Body içeriği User modeli ile eşleştirliyor.
	if value := JsonError(err, w, 500, "beklenmedik json parse hatası"); value == true { // Hata kontrolü
		return
	}

	a := new(Db) // Database Bağlantısı
	db, errdb := a.Connect()
	if value := JsonError(errdb, w, 500, "beklenmeyen bir hata (not connected db)"); value == true {
		return
	}

	if err := user.IsEmptyStringValues(); err != nil { // Gönderilen değerler boş mu kontrolü
		if value := JsonError(err, w, 412, "gerekli değerler gönderilmedi"); value == true {
			return
		}
	}

	id, errdb := user.InsertRow(db)                                                                       // DB'ye Kullanıcı Kaydı
	if value := JsonError(errdb, w, 412, "email veya telefon numarası sistemde kayıtlı"); value == true { // precondition faild
		return
	}
	a.Close(db) // DB bağlantısı kapatıldı.
	user.SetuID(id)

	var infoPage InfoPage // Response sayfası oluşturuldu ve değerleri işlendi.
	info := new(Info)
	info.InfoConstructer(true, 201, "veri kaydı başarılı")
	infoPage.InfoPageConstructer(info)

	data, err := json.Marshal(infoPage) // InfoPage nesnesi json'a parse ediliyor.
	if value := JsonError(err, w, 500, "beklenmedik json parse hatası"); value == true {
		return
	}

	w.Header().Set("Content-Type", "application/json") // Header değerleri ve response body'si
	w.WriteHeader(http.StatusCreated)
	w.Write(data)

}
