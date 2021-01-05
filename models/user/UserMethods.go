package user

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	. "oyeco-api/db"

	. "oyeco-api/helpers"

	. "oyeco-api/helpers/error"

	. "oyeco-api/models/info"
)

// Get - Set Metotları
func (user *User) GetuID() int {
	return user.ID
}

func (user *User) SetuID(id int) {
	user.ID = id
}

func (user *User) GetRecordTime() time.Time {
	return user.RecordTime
}

func (user *User) SetRecordTime(recordTime time.Time) {
	user.RecordTime = recordTime
}

func (user *User) GetIsVerifyEmail() bool {
	return user.IsVerifyEmail
}

func (user *User) SetIsVerifyEmail(isVerifyEmail bool) {
	user.IsVerifyEmail = isVerifyEmail
}

func (user *User) GetIsBlock() bool {
	return user.IsBlock
}

func (user *User) SetIsBlock(isBlock bool) {
	user.IsBlock = isBlock
}

// Fake User Constructer
func (user *User) UserConsturcter(firstName, lastName, phoneNumber, email, password, gender string, birthDay, recordTime time.Time, isVerifyEmail, isBlock bool) {
	user.FirstName = firstName
	user.LastName = lastName
	user.PhoneNumber = phoneNumber
	user.Email = email
	user.Password = password
	user.Gender = gender
	user.BirthDay = birthDay
	user.RecordTime = recordTime
	user.IsVerifyEmail = isVerifyEmail
	user.IsBlock = isBlock
}

// String değerlerden en az biri boş mu kontrolü
func (user *User) IsEmptyStringValues() error {
	if IsEmpty(user.FirstName) || IsEmpty(user.LastName) || IsEmpty(user.PhoneNumber) || IsEmpty(user.Email) || IsEmpty(user.Password) || IsEmpty(user.Gender) {
		return errors.New("beklenen değerler gönderilmemiş, değerleri kontol edin") // beklenen değerler gönderilmedi, değerlerinizi kontrol edin
	}
	return nil
}

// User Tablosuna Kayıt İşlemini Gerçekleştiriyor.
func (user *User) Create() (int, []byte) { // (int, []byte) => (statusCode, responseData)
	// Gönderilen değerler boş mu kontrolü
	if err := user.IsEmptyStringValues(); err != nil {
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
		INSERT INTO users (firstName, lastName, phoneNumber, email, password, gender, birthDay, recordTime, isVerifyEmail, isBlock)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING uID`
	id := 0
	encryptPass := fmt.Sprintf("%x", Encrypt([]byte(user.Password)))
	err := db.QueryRow(sqlStatement, user.FirstName, user.LastName, user.PhoneNumber, user.Email, encryptPass, user.Gender, user.BirthDay, user.RecordTime, user.IsVerifyEmail, user.IsBlock).Scan(&id)
	if value, data := JsonError(err, 412, "email veya telefon sistemde kayıtlı ya da veri tipi boyutları fazla"); value == true {
		return 500, data
	}
	defer db.Close() // DB bağlantısı kapatıldı.

	// Mail yollama
	encryptID := fmt.Sprintf("%x", Encrypt([]byte(strconv.Itoa(id)))) // Id şifrelenerek ekleniyor.
	msg := fmt.Sprintf("İyi Günler %s  %s,  \nBu mail hesabınızı aktif edebilmeniz için atılmaktadır. Hesabınızı aktif etmek için aşağıdaki maile tıklayabilirsiniz. (Lütfen linke tıklanmıyorsa kopyalayıp tarayıcınızda açın)\n\n%s/api/users/activation/%s \n\nSevgilerle,\nUpcycling", user.FirstName, user.LastName, os.Getenv("BaseURL"), encryptID)
	errmail := SendMail(user.Email, "Hesap Aktivasyon Maili", msg)
	if value, data := JsonError(errmail, 500, "beklenmeyen email gönderme hatası"); value == true {
		deleteStatement := `DELETE FROM users WHERE uID = $1;` // eğer aktivasyon maili atılamazsa kayıt silinir.
		db.Exec(deleteStatement, id)
		return 500, data
	}

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

func (user *User) Activation(id string) (int, []byte) {
	// Database Bağlantısı
	a := new(Db)
	db, errdb := a.Connect()
	if value, data := JsonError(errdb, 500, "veritabanı bağlantı hatası"); value == true { // Database bağlantı hatası
		return 500, data
	}

	decoded, _ := hex.DecodeString(id) // Şifrelenmiş id numarasını çözerek int hale çevirme
	decrypt, errDecr := Decrypt([]byte(decoded))
	if value, data := JsonError(errDecr, 400, "hatalı url"); value == true {
		return 400, data
	}
	identity, _ := strconv.Atoi(string(decrypt))

	sqlStatement := `UPDATE users SET isVerifyEmail = $2 WHERE uID = $1` // Güncelleme işlemi
	tx, _ := db.Begin()                                                  // Rollback yapmak için transaction başlatılıyor.
	_, err := db.Exec(sqlStatement, identity, true)
	if value, data := JsonError(err, 404, "kullanıcı kaydı bukunamadı, işlem başarısız"); value == true {
		tx.Rollback() // Rollback yapılıyor.
		return 404, data
	}
	_ = tx.Commit()  // Transaction durduruldu.
	defer db.Close() // DB bağlantısı kapatıldı.

	// Başarılı response için bilgi sayfası oluşturuluyor
	info := new(Info)
	info.InfoConstructer(true, "aktivasyon işlemi başarılı")
	infoPage := map[string]interface{}{"info": info} // Response sayfası oluşturuldu ve değerleri işlendi.
	data, _ := json.Marshal(infoPage)                // InfoPage nesnesi json'a parse ediliyor.

	return 200, data
}

func (user *User) SignIn() (int, []byte) {
	if IsEmpty(user.Email) || IsEmpty(user.Password) { // Gönderilen veriler boş mu?
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
	sqlStatement := `select * from users where email=$1`
	err := db.QueryRow(sqlStatement, user.Email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.PhoneNumber, &user.Email, &temp, &user.Gender, &user.BirthDay, &user.RecordTime, &user.IsVerifyEmail, &user.IsBlock)
	if value, data := JsonError(err, 404, "kullanıcı kaydı bukunamadı, işlem başarısız"); value == true {
		return 404, data
	}

	if user.IsBlock {
		if value, data := JsonError(errors.New("error"), 400, "user blocked"); value == true {
			return 400, data
		}
	} else if user.IsVerifyEmail == false {
		if value, data := JsonError(errors.New("error"), 400, "no user activation"); value == true {
			return 400, data
		}
	}
	decoded, _ := hex.DecodeString(temp) // Şifrelenmiş id numarasını çözerek int hale çevirme
	pass, _ := Decrypt([]byte(decoded))
	if user.Password != string(pass) {
		if value, data := JsonError(errors.New("error"), 400, "password is wrong"); value == true {
			return 400, data
		}
	}
	user.Password = temp

	info := new(Info)
	info.InfoConstructer(true, "giriş başarılı")
	infoPage := map[string]interface{}{"info": info, "content": user} // Response sayfası oluşturuldu ve değerleri işlendi.
	data, _ := json.Marshal(infoPage)                                 // InfoPage nesnesi json'a parse ediliyor.

	return 200, data
}

func (user *User) Update(id string) (int, []byte) {
	temp, _ := strconv.Atoi(id)
	user.SetuID(temp)
	// Gönderilen değerler boş mu kontrolü
	if IsEmpty(user.FirstName) || IsEmpty(user.LastName) || IsEmpty(user.PhoneNumber) || IsEmpty(user.Gender) {
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
	sqlStatement := `update users set firstName=$1, lastName=$2, phoneNumber=$3, gender=$4, birthDay=$5 where uID = $6 RETURNING uID`
	err := db.QueryRow(sqlStatement, user.FirstName, user.LastName, user.PhoneNumber, user.Gender, user.BirthDay, user.ID).Scan(&temp)
	if value, data := JsonError(err, 404, "güncelleme işlemi başarısız"); value == true {
		return 404, data
	}

	info := new(Info)
	info.InfoConstructer(true, "güncelleme işlemi başarılı")
	infoPage := map[string]interface{}{"info": info, "content": user} // Response sayfası oluşturuldu ve değerleri işlendi.
	data, _ := json.Marshal(infoPage)                                 // InfoPage nesnesi json'a parse ediliyor.

	return 200, data
}

func (user *User) UpdatePassword(id string) (int, []byte) {
	temp, _ := strconv.Atoi(id)
	user.SetuID(temp)
	// Gönderilen değerler boş mu kontrolü
	if IsEmpty(user.Password) {
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
	encryptPass := fmt.Sprintf("%x", Encrypt([]byte(user.Password)))
	sqlStatement := `update users set password=$1 where uID = $2 RETURNING uID`
	err := db.QueryRow(sqlStatement, encryptPass, user.ID).Scan(&temp)
	if value, data := JsonError(err, 404, "güncelleme işlemi başarısız"); value == true {
		return 404, data
	}

	info := new(Info)
	info.InfoConstructer(true, "güncelleme işlemi başarılı")
	infoPage := map[string]interface{}{"info": info} // Response sayfası oluşturuldu ve değerleri işlendi.
	data, _ := json.Marshal(infoPage)                // InfoPage nesnesi json'a parse ediliyor.

	return 200, data
}
