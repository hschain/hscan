package schema

type PersonAlassets struct {
	ID      int32  `json:"id" gorm:"pk"`
	Address string `json:"address" gorm:"not null"`
	Amount  int64  `json:"amount" gorm:"not null"`
	Denom   string `json:"denom" gorm:"not null"`
}
