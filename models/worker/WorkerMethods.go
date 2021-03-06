package worker

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	. "oyeco-api/db"

	. "oyeco-api/helpers"

	. "oyeco-api/helpers/error"

	. "oyeco-api/models/info"
)

func (worker *Worker) SetRecordTime(recordTime time.Time) {
	worker.RecordTime = recordTime
}

func (worker *Worker) SetStatus(status int) {
	worker.Status = status
}

// String değerlerden en az biri boş mu kontrolü
func (worker *Worker) IsEmptyStringValues() error {
	if IsEmpty(worker.FirstName) || IsEmpty(worker.LastName) || IsEmpty(worker.PhoneNumber) || IsEmpty(worker.Email) || IsEmpty(worker.Password) || IsEmpty(worker.Gender) || worker.BirthDay.IsZero() {
		return errors.New("beklenen degerler gonderilmemis, degerleri kontol edin") // beklenen değerler gönderilmedi, değerlerinizi kontrol edin
	}
	return nil
}

func (worker *Worker) IsEmptyUpdate() error {
	if IsEmpty(worker.FirstName) || IsEmpty(worker.LastName) || IsEmpty(worker.PhoneNumber) || IsEmpty(worker.Email) || IsEmpty(worker.Password) || worker.BirthDay.IsZero() {
		return errors.New("beklenen degerler gonderilmemis, degerleri kontol edin") // beklenen değerler gönderilmedi, değerlerinizi kontrol edin
	}
	return nil
}

func (worker *Worker) ManagerSignIn() (int, []byte) {
	if IsEmpty(worker.Email) || IsEmpty(worker.Password) { // Gönderilen veriler boş mu?
		if value, data := JsonError(errors.New("error"), 400, "beklenen değerler gonderilmemis, degerleri kontrol edin"); value == true {
			return 400, data
		}
	}
	// Database Bağlantısı
	a := new(Db)
	db, errdb := a.Connect()
	if value, data := JsonError(errdb, 500, "veritabani baglanti hatasi"); value == true { // Database bağlantı hatası
		return 500, data
	}
	defer db.Close()
	var temp string // şifrelenmiş password'ü tutacak.
	sqlStatement := `select * from workers where status = $1 and email = $2`
	err := db.QueryRow(sqlStatement, 1, worker.Email).Scan(&worker.WID, &worker.FirstName, &worker.LastName, &worker.PhoneNumber, &worker.Email, &temp, &worker.Gender, &worker.BirthDay, &worker.RecordTime, &worker.Status)
	if value, data := JsonError(err, 404, "yonetici kaydı bulunamadi, islem basarisiz"); value == true {
		return 404, data
	}

	decoded, _ := hex.DecodeString(temp) // Şifrelenmiş id numarasını çözerek gerçek şifre hale çevirme
	pass, _ := Decrypt([]byte(decoded))
	if worker.Password != string(pass) {
		if value, data := JsonError(errors.New("error"), 400, "sifre hatali"); value == true {
			return 400, data
		}
	}
	worker.Password = temp

	info := new(Info)
	info.InfoConstructer(true, "giriş başarılı")
	infoPage := map[string]interface{}{"info": info, "content": worker} // Response sayfası oluşturuldu ve değerleri işlendi.
	data, _ := json.Marshal(infoPage)                                   // InfoPage nesnesi json'a parse ediliyor.

	return 200, data
}

func (worker *Worker) FieldWorkerSignIn() (int, []byte) {
	if IsEmpty(worker.Email) || IsEmpty(worker.Password) { // Gönderilen veriler boş mu?
		if value, data := JsonError(errors.New("error"), 400, "beklenen degerler gonderilmemis, degerleri kontrol edin"); value == true {
			return 400, data
		}
	}
	// Database Bağlantısı
	a := new(Db)
	db, errdb := a.Connect()
	if value, data := JsonError(errdb, 500, "veritabani baglanti hatasi"); value == true { // Database bağlantı hatası
		return 500, data
	}
	defer db.Close()
	var temp string // şifrelenmiş password'ü tutacak.
	sqlStatement := `select * from workers where status = $1 and email = $2`
	err := db.QueryRow(sqlStatement, 0, worker.Email).Scan(&worker.WID, &worker.FirstName, &worker.LastName, &worker.PhoneNumber, &worker.Email, &temp, &worker.Gender, &worker.BirthDay, &worker.RecordTime, &worker.Status)
	if value, data := JsonError(err, 404, "saha calisani kaydi bulunamadi, islem basarisiz"); value == true {
		return 404, data
	}

	decoded, _ := hex.DecodeString(temp) // Şifrelenmiş id numarasını çözerek gerçek şifre hale çevirme
	pass, _ := Decrypt([]byte(decoded))
	if worker.Password != string(pass) {
		if value, data := JsonError(errors.New("error"), 400, "sifre hatali"); value == true {
			return 400, data
		}
	}
	worker.Password = temp

	info := new(Info)
	info.InfoConstructer(true, "giris basarili")
	infoPage := map[string]interface{}{"info": info, "content": worker} // Response sayfası oluşturuldu ve değerleri işlendi.
	data, _ := json.Marshal(infoPage)                                   // InfoPage nesnesi json'a parse ediliyor.

	return 200, data
}

func (worker *Worker) Create() (int, []byte) { // (int, []byte) => (statusCode, responseData)
	// Gönderilen değerler boş mu kontrolü
	if err := worker.IsEmptyStringValues(); err != nil {
		if value, data := JsonError(err, 412, err.Error()); value == true {
			return 412, data
		}
	}
	// Database Bağlantısı
	a := new(Db)
	db, errdb := a.Connect()
	if value, data := JsonError(errdb, 500, "veritabani baglanti hatasi"); value == true { // Database bağlantı hatası
		return 500, data
	}

	sqlStatement := `
		INSERT INTO workers (firstName, lastName, phoneNumber, email, password, gender, birthDay, recordTime, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING wID`
	id := 0
	encryptPass := fmt.Sprintf("%x", Encrypt([]byte(worker.Password)))
	err := db.QueryRow(sqlStatement, worker.FirstName, worker.LastName, worker.PhoneNumber, worker.Email, encryptPass, worker.Gender, worker.BirthDay, worker.RecordTime, worker.Status).Scan(&id)
	if value, data := JsonError(err, 412, "email veya telefon sistemde kayitli ya da veri tipi boyutlari fazla"); value == true {
		return 500, data
	}
	defer db.Close() // DB bağlantısı kapatıldı.

	// Başarılı response için bilgi sayfası oluşturuluyor
	info := new(Info)
	info.InfoConstructer(true, "veri kaydi basarili")
	infoPage := map[string]interface{}{"info": info} // Response sayfası oluşturuldu ve değerleri işlendi.

	data, err := json.Marshal(infoPage) // InfoPage nesnesi json'a parse ediliyor.
	if value, data := JsonError(err, 500, "beklenmedik json parse hatasi"); value == true {
		return 500, data
	}

	return 201, data // Başarılı response return yapılıyor.
}

func (worker *Worker) AllGet() (int, []byte) {
	var wrkr []Worker
	// Database Bağlantısı
	a := new(Db)
	db, errdb := a.Connect()
	if value, data := JsonError(errdb, 500, "veritabani baglanti hatasi"); value == true { // Database bağlantı hatası
		return 500, data
	}
	defer db.Close()

	// Kullanıcı adresleri getirilmesi için
	rows, _ := db.Query("Select * from workers where status = 0")
	temp := 0
	for rows.Next() {
		temp = 1
		err := rows.Scan(&worker.WID, &worker.FirstName, &worker.LastName, &worker.PhoneNumber, &worker.Email, &worker.Password, &worker.Gender, &worker.BirthDay, &worker.RecordTime, &worker.Status)
		if value, data := JsonError(err, 404, "veri hatasi"); value == true {
			return 404, data
		}
		wrkr = append(wrkr, *worker)
	}
	if temp == 0 {
		if value, data := JsonError(errors.New("error"), 404, "sistemde kayıtli saha calisani bulunmamaktadir"); value == true {
			return 404, data
		}
	}

	info := new(Info)
	info.InfoConstructer(true, "saha calisani listesi")
	infoPage := map[string]interface{}{"info": info, "content": wrkr} // Response sayfası oluşturuldu ve değerleri işlendi.
	data, _ := json.Marshal(infoPage)                                 // InfoPage nesnesi json'a parse ediliyor.

	return 200, data
}

func (worker *Worker) Delete(wID string) (int, []byte) {
	temp, _ := strconv.Atoi(wID)
	worker.WID = temp
	// Database Bağlantısı
	a := new(Db)
	db, errdb := a.Connect()
	if value, data := JsonError(errdb, 500, "veritabani baglanti hatasi"); value == true { // Database bağlantı hatası
		return 500, data
	}
	defer db.Close()
	sqlStatement := `DELETE FROM workers WHERE wID = $1 AND status = 0`
	res, err := db.Exec(sqlStatement, worker.WID)
	if value, data := JsonError(err, 404, "silme islemi basarisiz"); value == true {
		return 404, data
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		if value, data := JsonError(errors.New("silme"), 404, "silme islemi basarisiz"); value == true {
			return 404, data
		}
	}

	info := new(Info)
	info.InfoConstructer(true, "silme islemi basarili")
	infoPage := map[string]interface{}{"info": info} // Response sayfası oluşturuldu ve değerleri işlendi.
	data, _ := json.Marshal(infoPage)                // InfoPage nesnesi json'a parse ediliyor.

	return 200, data
}

func (worker *Worker) Update(wID string) (int, []byte) {
	temp, _ := strconv.Atoi(wID)
	worker.WID = temp
	// Gönderilen değerler boş mu kontrolü
	if err := worker.IsEmptyUpdate(); err != nil {
		if value, data := JsonError(err, 412, err.Error()); value == true {
			return 412, data
		}
	}
	// Database Bağlantısı
	a := new(Db)
	db, errdb := a.Connect()
	if value, data := JsonError(errdb, 500, "veritabani baglanti hatasi"); value == true { // Database bağlantı hatası
		return 500, data
	}
	defer db.Close()
	encryptPass := fmt.Sprintf("%x", Encrypt([]byte(worker.Password)))
	sqlStatement := `update workers set firstName=$1, lastName=$2, birthDay=$3, email=$4, password=$5, phoneNumber=$6 where wID = $7 and status = 0 RETURNING wID`
	err := db.QueryRow(sqlStatement, worker.FirstName, worker.LastName, worker.BirthDay, worker.Email, encryptPass, worker.PhoneNumber, worker.WID).Scan(&temp)
	if value, data := JsonError(err, 404, "güncelleme işlemi başarısız"); value == true {
		return 404, data
	}
	worker.Password = encryptPass

	info := new(Info)
	info.InfoConstructer(true, "guncelleme islemi basarili")
	infoPage := map[string]interface{}{"info": info, "content": worker} // Response sayfası oluşturuldu ve değerleri işlendi.
	data, _ := json.Marshal(infoPage)                                   // InfoPage nesnesi json'a parse ediliyor.

	return 200, data
}
