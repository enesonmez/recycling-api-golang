package worker

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
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
		return errors.New("beklenen değerler gönderilmemiş, değerleri kontol edin") // beklenen değerler gönderilmedi, değerlerinizi kontrol edin
	}
	return nil
}

func (worker *Worker) ManagerSignIn() (int, []byte) {
	if IsEmpty(worker.Email) || IsEmpty(worker.Password) { // Gönderilen veriler boş mu?
		if value, data := JsonError(errors.New("error"), 400, "beklenen değerler gönderilmemiş, değerleri kontrol edin"); value == true {
			return 400, data
		}
	}
	// Database Bağlantısı
	a := new(Db)
	db, errdb := a.Connect()
	if value, data := JsonError(errdb, 500, "veritabanı bağlantı hatası"); value == true { // Database bağlantı hatası
		return 500, data
	}
	defer db.Close()
	var temp string // şifrelenmiş password'ü tutacak.
	sqlStatement := `select * from workers where status = $1 and email = $2`
	err := db.QueryRow(sqlStatement, 1, worker.Email).Scan(&worker.WID, &worker.FirstName, &worker.LastName, &worker.PhoneNumber, &worker.Email, &temp, &worker.Gender, &worker.BirthDay, &worker.RecordTime, &worker.Status)
	if value, data := JsonError(err, 404, "yönetici kaydı bulunamadı, işlem başarısız"); value == true {
		return 404, data
	}

	decoded, _ := hex.DecodeString(temp) // Şifrelenmiş id numarasını çözerek gerçek şifre hale çevirme
	pass, _ := Decrypt([]byte(decoded))
	if worker.Password != string(pass) {
		if value, data := JsonError(errors.New("error"), 400, "şifre hatalı"); value == true {
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
	if value, data := JsonError(errdb, 500, "veritabanı bağlantı hatası"); value == true { // Database bağlantı hatası
		return 500, data
	}

	sqlStatement := `
		INSERT INTO workers (firstName, lastName, phoneNumber, email, password, gender, birthDay, recordTime, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING wID`
	id := 0
	encryptPass := fmt.Sprintf("%x", Encrypt([]byte(worker.Password)))
	err := db.QueryRow(sqlStatement, worker.FirstName, worker.LastName, worker.PhoneNumber, worker.Email, encryptPass, worker.Gender, worker.BirthDay, worker.RecordTime, worker.Status).Scan(&id)
	if value, data := JsonError(err, 412, "email veya telefon sistemde kayıtlı ya da veri tipi boyutları fazla"); value == true {
		return 500, data
	}
	defer db.Close() // DB bağlantısı kapatıldı.

	// Başarılı response için bilgi sayfası oluşturuluyor
	info := new(Info)
	info.InfoConstructer(true, "veri kaydı başarılı")
	infoPage := map[string]interface{}{"info": info} // Response sayfası oluşturuldu ve değerleri işlendi.

	data, err := json.Marshal(infoPage) // InfoPage nesnesi json'a parse ediliyor.
	if value, data := JsonError(err, 500, "beklenmedik json parse hatası"); value == true {
		return 500, data
	}

	return 201, data // Başarılı response return yapılıyor.
}

func (worker *Worker) AllGet() (int, []byte) {
	var wrkr []Worker
	// Database Bağlantısı
	a := new(Db)
	db, errdb := a.Connect()
	if value, data := JsonError(errdb, 500, "veritabanı bağlantı hatası"); value == true { // Database bağlantı hatası
		return 500, data
	}
	defer db.Close()

	// Kullanıcı adresleri getirilmesi için
	rows, _ := db.Query("Select * from workers where status = 0")
	temp := 0
	for rows.Next() {
		temp = 1
		err := rows.Scan(&worker.WID, &worker.FirstName, &worker.LastName, &worker.PhoneNumber, &worker.Email, &worker.Password, &worker.Gender, &worker.BirthDay, &worker.RecordTime, &worker.Status)
		if value, data := JsonError(err, 404, "veri hatası"); value == true {
			return 404, data
		}
		wrkr = append(wrkr, *worker)
	}
	if temp == 0 {
		if value, data := JsonError(errors.New("error"), 404, "sistemde kayıtlı saha çalışanı bulunmamaktadır"); value == true {
			return 404, data
		}
	}

	info := new(Info)
	info.InfoConstructer(true, "saha çalışanı listesi")
	infoPage := map[string]interface{}{"info": info, "content": wrkr} // Response sayfası oluşturuldu ve değerleri işlendi.
	data, _ := json.Marshal(infoPage)                                 // InfoPage nesnesi json'a parse ediliyor.

	return 200, data
}
