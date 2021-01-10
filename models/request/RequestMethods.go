package request

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	. "oyeco-api/db"

	. "oyeco-api/helpers/error"

	. "oyeco-api/models/info"
)

func (request *Request) SetRequestCreateTime(createTime time.Time) {
	request.RequestCreateTime = createTime
}

func (request *Request) SetState(state int) {
	request.State = state
}

func (request *Request) Create(userID string) (int, []byte) { // (int, []byte) => (statusCode, responseData)
	temp, _ := strconv.Atoi(userID)
	request.UserID = temp

	// Database Bağlantısı
	a := new(Db)
	db, errdb := a.Connect()
	if value, data := JsonError(errdb, 500, "database connection error"); value == true { // Database bağlantı hatası
		return 500, data
	}

	// Aynı kullanıcı adına ve adresine sahip bekleniyor ya da onaylandı durumunda bir kayıt var mı?
	rows, _ := db.Query("Select * from requests where (userID = $1 and addressID = $2) and (state = $3 or state = $4)", request.UserID, request.AddressID, 1, 2)
	for rows.Next() {
		if value, data := JsonError(errors.New("error"), 404, "bu adrese ait bekleyen ya da onaylanan bir adres bulunmaktadır"); value == true {
			return 404, data
		}
	}

	sqlStatement := `
		INSERT INTO requests (userID, addressID, requestCreateTime, state)
		VALUES ($1, $2, $3, $4)
		RETURNING reqID`
	ids := 0
	err := db.QueryRow(sqlStatement, request.UserID, request.AddressID, request.RequestCreateTime, request.State).Scan(&ids)
	if value, data := JsonError(err, 412, "böyle bir kullanıcı veya adres bulunmuyor"); value == true {
		return 500, data
	}
	defer db.Close() // DB bağlantısı kapatıldı.

	// Başarılı response için bilgi sayfası oluşturuluyor
	info := new(Info)
	info.InfoConstructer(true, "veri kaydı başarılı")
	infoPage := map[string]interface{}{"info": info} // Response sayfası oluşturuldu ve değerleri işlendi.

	data, err := json.Marshal(infoPage) // InfoPage nesnesi json'a parse ediliyor.
	if value, data := JsonError(err, 500, "beklenmeyen json parse hatası"); value == true {
		return 500, data
	}

	return 201, data // Başarılı response return yapılıyor.
}
