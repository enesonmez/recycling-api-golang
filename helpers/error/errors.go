package error

import (
	"encoding/json"
	"log"

	. "oyeco-api/models/info"
)

// Genel hata kontrol fonksiyonu
func CheckEror(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

// Hata meydana  geldiğinde json hata sayfasını oluşturur
func JsonError(err error, status int, message string) (bool, []byte) { // true dönerse handlers return edecek
	if err != nil {
		info := new(Info)
		info.InfoConstructer(false, message)             // Info nesnesine değerler setleniyor.
		infoPage := map[string]interface{}{"info": info} // Hata sayfası oluşturuldu.
		data, err := json.Marshal(infoPage)              // Map parse edildi.
		CheckEror(err)

		return true, data
	}
	return false, []byte("")
}
