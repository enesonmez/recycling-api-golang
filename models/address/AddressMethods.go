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

func (address *Address) Create(id string) (int, []byte) { // (int, []byte) => (statusCode, responseData)
	temp, _ := strconv.Atoi(id)
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
