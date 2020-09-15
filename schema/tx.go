package schema

import (
	"time"
	//sdk "github.com/cosmos/cosmos-sdk/types"
)

type Message struct {
	MsgIndex uint16                       `json:"msg_index"`
	Success  bool                         `json:"success"`
	Log      string                       `json:"log"`
	Events   map[string]map[string]string `json:"events"`
}

// Transaction defines the schema for transaction information
type Transaction struct {
	ID          int32  `json:"id" gorm:"pk"`
	Height      int64  `json:"height" gorm:"not null"`
	TxHash      string `json:"tx_hash" gorm:"not null;unique"`
	Code        uint32 `json:"code"  gorm:",notnull"`
	RawMessages string `json:"-" gorm:"type:json;not null"`
	//Messages    sdk.ABCIMessageLogs `json:"messages" gorm:"-"`
	Messages  []*Message `json:"messages" gorm:"-"`
	Memo      string     `json:"memo"`
	Fee       string     `json:"fee"`
	GasWanted int64      `json:"gas_wanted" gorm:"default:0"`
	GasUsed   int64      `json:"gas_used" gorm:"default:0"`
	Timestamp time.Time  `json:"timestamp" gorm:"default:now()"`
	Sender    string     `json:"sender"`
	Recipient string     `json:"recipient"`
	Amount    string     `json:"amount"`
}

// Transaction defines the schema for transaction information
type RavlTransaction struct {
	ID        int32      `json:"id" gorm:"pk"`
	Height    int64      `json:"height" gorm:"not null"`
	TxHash    string     `json:"tx_hash" gorm:"not null;unique"`
	Messages  []*Message `json:"messages" gorm:"-"`
	Memo      string     `json:"memo"`
	Timestamp time.Time  `json:"timestamp" gorm:"default:now()"`
}

type SignedTx struct {
	Tx   string `json:"tx" gorm:"pk"`
	Mode string `json:"mode" gorm:"not null"`
}
