package address

type Address struct {
	AID         int    `json:"aid"`
	FullAddress string `json:"fullAddress"`
	District    string `json:"district"`
	City        string `json:"city"`
	Postcode    string `json:"postcode"`
	UserID      int    `json:"userID"`
}
