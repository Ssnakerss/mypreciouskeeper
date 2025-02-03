package models

import "time"

//Assets types

const (
	AssetTypeCredentials = "CRED"
	AssetTypeCard        = "CARD"
	AssetTypeMemo        = "MEMO"
	AssetTypeFile        = "FILE"
)

//Asset general struct for work with storage and transport
type Asset struct {
	ID     int64
	UserID int64

	Type      string
	Sticker   string
	Body      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

//Credential struct to store credentials precious
type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

//Memo struct to store text precious
type Memo struct {
	Text string `json:"text"`
}

//Card struct to store card precious
type Card struct {
	Number   string `json:"number"`
	Name     string `json:"name"`
	ExpMonth int    `json:"expmonth"`
	ExpYear  int    `json:"expyear"`
	CVV      string `json:"cvv"`
}
