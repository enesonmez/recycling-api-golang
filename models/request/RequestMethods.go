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
	request.ReqID = ids

	// Başarılı response için bilgi sayfası oluşturuluyor
	info := new(Info)
	info.InfoConstructer(true, "veri kaydı başarılı")
	infoPage := map[string]interface{}{"info": info, "content": request} // Response sayfası oluşturuldu ve değerleri işlendi.

	data, err := json.Marshal(infoPage) // InfoPage nesnesi json'a parse ediliyor.
	if value, data := JsonError(err, 500, "beklenmeyen json parse hatası"); value == true {
		return 500, data
	}

	return 201, data // Başarılı response return yapılıyor.
}

func (request *Request) Get(userID string) (int, []byte) { // (int, []byte) => (statusCode, responseData)
	var rqst []Request
	temp, _ := strconv.Atoi(userID)
	request.UserID = temp
	// Database Bağlantısı
	a := new(Db)
	db, errdb := a.Connect()
	if value, data := JsonError(errdb, 500, "veritabanı bağlantı hatası"); value == true { // Database bağlantı hatası
		return 500, data
	}
	defer db.Close()

	// Kullanıcı adresleri getirilmesi için
	rows, _ := db.Query("Select * from requests where userID = $1", request.UserID)

	temp = 0
	for rows.Next() {
		temp = 1
		err := rows.Scan(&request.ReqID, &request.UserID, &request.AddressID, &request.RequestCreateTime, &request.State)
		if value, data := JsonError(err, 404, "veri hatası"); value == true {
			return 404, data
		}
		rqst = append(rqst, *request)
	}
	if temp == 0 {
		if value, data := JsonError(errors.New("error"), 404, "böyle bir kullanıcı bulunmuyor veya hiçbir istek yok"); value == true {
			return 404, data
		}
	}

	info := new(Info)
	info.InfoConstructer(true, "kullanıcı istekleri")
	infoPage := map[string]interface{}{"info": info, "content": rqst} // Response sayfası oluşturuldu ve değerleri işlendi.
	data, _ := json.Marshal(infoPage)                                 // InfoPage nesnesi json'a parse ediliyor.

	return 200, data
}

func (request *Request) Delete(userID, reqID string) (int, []byte) {
	temp, _ := strconv.Atoi(userID)
	request.UserID = temp
	temp, _ = strconv.Atoi(reqID)
	request.ReqID = temp
	// Database Bağlantısı
	a := new(Db)
	db, errdb := a.Connect()
	if value, data := JsonError(errdb, 500, "veritabanı bağlantı hatası"); value == true { // Database bağlantı hatası
		return 500, data
	}
	defer db.Close()
	sqlStatement := `DELETE FROM requests WHERE reqID = $1 and userID = $2`
	res, _ := db.Exec(sqlStatement, request.ReqID, request.UserID)

	count, _ := res.RowsAffected()
	if count == 0 {
		if value, data := JsonError(errors.New("silme"), 404, "silme işlemi başarısız"); value == true {
			return 404, data
		}
	}

	info := new(Info)
	info.InfoConstructer(true, "silme işlemi başarılı")
	infoPage := map[string]interface{}{"info": info} // Response sayfası oluşturuldu ve değerleri işlendi.
	data, _ := json.Marshal(infoPage)                // InfoPage nesnesi json'a parse ediliyor.

	return 200, data
}
