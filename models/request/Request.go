package request

import "time"

type Request struct {
	ReqID             int       `json:"reqID"`
	UserID            int       `json:"userID"`
	AddressID         int       `json:"addressID"`
	RequestCreateTime time.Time `json:"requestCreateTime"`
	State             int       `json:"state"`
}
