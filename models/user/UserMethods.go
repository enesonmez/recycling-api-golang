package user

import (
	"database/sql"
	"errors"
	"time"

	. "../../helpers"
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

// User Tablosuna Kayıt İşlemini Gerçekleştiriyor.
func (user *User) InsertRow(db *sql.DB) (int, error) {
	sqlStatement := `
		INSERT INTO users (firstName, lastName, phoneNumber, email, password, gender, birthDay, recordTime, isVerifyEmail, isBlock)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING uID`
	id := 0
	tx, _ := db.Begin() // callback yapmak için transaction başlatılıyor.
	err := db.QueryRow(sqlStatement, user.FirstName, user.LastName, user.PhoneNumber, user.Email, user.Password, user.Gender, user.BirthDay, user.RecordTime, user.IsVerifyEmail, user.IsBlock).Scan(&id)
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	_ = tx.Commit()
	return id, nil
}

// String değerlerden en az biri boş mu kontrolü
func (user *User) IsEmptyStringValues() error {
	if IsEmpty(user.FirstName) || IsEmpty(user.LastName) || IsEmpty(user.PhoneNumber) || IsEmpty(user.Email) || IsEmpty(user.Password) || IsEmpty(user.Gender) {
		return errors.New("beklenilen değerler yollanmamış, değerlerinizi kontrol edin")
	}
	return nil
}
