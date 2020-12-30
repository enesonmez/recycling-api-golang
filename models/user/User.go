package user

import (
	"time"
)

type User struct {
	ID            int       `json:"id"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	PhoneNumber   string    `json:"phoneNumber"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	Gender        string    `json:"gender"`
	BirthDay      time.Time `json:"birthDay"`
	RecordTime    time.Time `json:"-"`
	IsVerifyEmail bool      `json:"-"`
	IsBlock       bool      `json:"-"`
}
