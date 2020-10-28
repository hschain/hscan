package schema

import (
	"time"
	//sdk "github.com/cosmos/cosmos-sdk/types"
)

/*type Message struct {
	MsgIndex int    `json:"msg_index"`
	Success  bool   `json:"success"`
	Log      string `json:"log"`
	Events   []struct {
		Type       string `json:"type"`
		Attributes []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"attributes"`
	} `json:"events"`
}*/

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
	Messages        []Message `json:"messages" gorm:"-"`
	Memo            string    `json:"memo"`
	Fee             string    `json:"fee"`
	GasWanted       int64     `json:"gas_wanted" gorm:"default:0"`
	GasUsed         int64     `json:"gas_used" gorm:"default:0"`
	Timestamp       time.Time `json:"timestamp" gorm:"default:now()"`
	Sender          string    `json:"sender" gorm:"not null"`
	Recipient       string    `json:"recipient" gorm:"not null"`
	Amount          string    `json:"amount" gorm:"not null"`
	Denom           string    `json:"denom" gorm:"not null"`
	SenderNotice    int       `json:"sender_notice" gorm:"not null"`
	RecipientNotice int       `json:"recipient_notice" gorm:"not null"`
}

// Transaction defines the schema for transaction information
type RavlTransaction struct {
	ID        int32     `json:"id" gorm:"pk"`
	Height    int64     `json:"height" gorm:"not null"`
	TxHash    string    `json:"tx_hash" gorm:"not null;unique"`
	Messages  []Message `json:"messages" gorm:"-"`
	Memo      string    `json:"memo"`
	Timestamp time.Time `json:"timestamp" gorm:"default:now()"`
}

type SignedTx struct {
	Tx   string `json:"tx" gorm:"pk"`
	Mode string `json:"mode" gorm:"not null"`
}

// NoneInfo defines the schema for nodeinfo information
type NodeInfo struct {
	ID       int32  `json:"id" gorm:"pk"`
	NodeName string `json:"name"  gorm:"not null"`
	NodeUrl  string `json:"url" gorm:"not null"`
}
