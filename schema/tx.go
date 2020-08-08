package schema

import "time"

// Transaction defines the schema for transaction information
type Transaction struct {
	ID        int32     `json:"id" gorm:"pk"`
	Height    int64     `json:"height" gorm:"not null"`
	TxHash    string    `json:"tx_hash" gorm:"not null;unique"`
	Code      uint32    `json:"code"  gorm:",notnull"`
	Messages  string    `json:"messages" gorm:"type:json;not null"`
	Memo      string    `json:"memo"`
	Fee       string    `json:"fee"`
	GasWanted int64     `json:"gas_wanted" gorm:"default:0"`
	GasUsed   int64     `json:"gas_used" gorm:"default:0"`
	Timestamp time.Time `json:"timestamp" gorm:"default:now()"`
}
