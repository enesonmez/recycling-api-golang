package worker

import (
	"time"
)

type Worker struct {
	WID         int       `json:"wID"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	PhoneNumber string    `json:"phoneNumber"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Gender      string    `json:"gender"`
	BirthDay    time.Time `json:"birthDay"`
	RecordTime  time.Time `json:"-"`
	Status      int       `json:"-"`
}
