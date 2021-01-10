package address

import (
	"encoding/json"
	"errors"
	"strconv"

	. "oyeco-api/db"

	. "oyeco-api/helpers"

	. "oyeco-api/helpers/error"

	. "oyeco-api/models/info"
)

func (address *Address) IsEmptyStringValues() error {
	if IsEmpty(address.FullAddress) || IsEmpty(address.District) || IsEmpty(address.City) || IsEmpty(address.Postcode) {
		return errors.New("expected values not sent, check your values") // beklenen değerler gönderilmedi, değerlerinizi kontrol edin
	}
	return nil
}

// Adres kaydı yapar
func (address *Address) Create(userID string) (int, []byte) { // (int, []byte) => (statusCode, responseData)
	temp, _ := strconv.Atoi(userID)
	address.UserID = temp
	// Gönderilen değerler boş mu kontrolü
	if err := address.IsEmptyStringValues(); err != nil {
		if value, data := JsonError(err, 412, err.Error()); value == true {
			return 412, data
		}
	}
	// Database Bağlantısı
	a := new(Db)
	db, errdb := a.Connect()
	if value, data := JsonError(errdb, 500, "database connection error"); value == true { // Database bağlantı hatası
		return 500, data
	}

	sqlStatement := `
		INSERT INTO address (fullAddress, district, city, postcode, userID)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING aID`
	ids := 0
	err := db.QueryRow(sqlStatement, address.FullAddress, address.District, address.City, address.Postcode, address.UserID).Scan(&ids)
	if value, data := JsonError(err, 412, "böyle bir kullanıcı bulunmuyor"); value == true {
		return 500, data
	}
	defer db.Close() // DB bağlantısı kapatıldı.

	// Başarılı response için bilgi sayfası oluşturuluyor
	info := new(Info)
	info.InfoConstructer(true, "veri kaydı başarılı")
	infoPage := map[string]interface{}{"info": info} // Response sayfası oluşturuldu ve değerleri işlendi.

	data, err := json.Marshal(infoPage) // InfoPage nesnesi json'a parse ediliyor.
	if value, data := JsonError(err, 500, "unexpected json parse error"); value == true {
		return 500, data
	}

	return 201, data // Başarılı response return yapılıyor.
}

// Adres güncellemei yapar
func (address *Address) Update(userID, adrsID string) (int, []byte) {
	temp, _ := strconv.Atoi(userID)
	address.UserID = temp
	temp, _ = strconv.Atoi(adrsID)
	address.AID = temp
	// Gönderilen değerler boş mu kontrolü
	if err := address.IsEmptyStringValues(); err != nil {
		if value, data := JsonError(err, 412, err.Error()); value == true {
			return 412, data
		}
	}
	// Database Bağlantısı
	a := new(Db)
	db, errdb := a.Connect()
	if value, data := JsonError(errdb, 500, "veritabanı bağlantı hatası"); value == true { // Database bağlantı hatası
		return 500, data
	}
	defer db.Close()
	sqlStatement := `update address set fullAddress=$1, district=$2, city=$3, postcode=$4 where userID = $5 and aID = $6 RETURNING aID`
	err := db.QueryRow(sqlStatement, address.FullAddress, address.District, address.City, address.Postcode, address.UserID, address.AID).Scan(&temp)
	if value, data := JsonError(err, 404, "güncelleme işlemi başarısız"); value == true {
		return 404, data
	}

	info := new(Info)
	info.InfoConstructer(true, "güncelleme işlemi başarılı")
	infoPage := map[string]interface{}{"info": info, "content": address} // Response sayfası oluşturuldu ve değerleri işlendi.
	data, _ := json.Marshal(infoPage)                                    // InfoPage nesnesi json'a parse ediliyor.

	return 200, data
}

// Kullanıcıya ait adresleri döndürür
func (address *Address) Get(id string) (int, []byte) {
	var adrs []Address
	temp, _ := strconv.Atoi(id)
	address.UserID = temp
	// Database Bağlantısı
	a := new(Db)
	db, errdb := a.Connect()
	if value, data := JsonError(errdb, 500, "veritabanı bağlantı hatası"); value == true { // Database bağlantı hatası
		return 500, data
	}
	defer db.Close()

	// Kullanıcı adresleri getirilmesi için
	rows, err := db.Query("Select * from address where userID = $1", address.UserID)
	if value, data := JsonError(err, 404, "kullanıcı adresleri getirilemedi (kulllanıcı adı yanlış olabilir)"); value == true {
		return 404, data
	}
	temp = 0
	for rows.Next() {
		temp = 1
		err = rows.Scan(&address.AID, &address.FullAddress, &address.District, &address.City, &address.Postcode, &address.UserID)
		if value, data := JsonError(err, 404, "veri hatası"); value == true {
			return 404, data
		}
		adrs = append(adrs, *address)
	}
	if temp == 0 {
		if value, data := JsonError(errors.New("böyle bir kullanıcı bulunmuyor"), 404, "böyle bir kullanıcı bulunmuyor"); value == true {
			return 404, data
		}
	}

	info := new(Info)
	info.InfoConstructer(true, "kullanıcı adresleri")
	infoPage := map[string]interface{}{"info": info, "content": adrs} // Response sayfası oluşturuldu ve değerleri işlendi.
	data, _ := json.Marshal(infoPage)                                 // InfoPage nesnesi json'a parse ediliyor.

	return 200, data
}

func (address *Address) Delete(userID, adrsID string) (int, []byte) {
	temp, _ := strconv.Atoi(userID)
	address.UserID = temp
	temp, _ = strconv.Atoi(adrsID)
	address.AID = temp
	// Database Bağlantısı
	a := new(Db)
	db, errdb := a.Connect()
	if value, data := JsonError(errdb, 500, "veritabanı bağlantı hatası"); value == true { // Database bağlantı hatası
		return 500, data
	}
	defer db.Close()
	sqlStatement := `DELETE FROM address WHERE aID = $1 and userID = $2`
	res, err := db.Exec(sqlStatement, address.AID, address.UserID)
	if value, data := JsonError(err, 404, "silme işlemi başarısız"); value == true {
		return 404, data
	}
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
