package models

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Tx struct {
	ID        int32               `json:"id"`
	Height    int64               `json:"height"`
	TxHash    string              `json:"tx_hash"`
	Code      uint32              `json:"code"`
	Messages  sdk.ABCIMessageLogs `json:"messages"`
	Memo      string              `json:"memo"`
	Fee       string              `json:"fee"`
	GasWanted int64               `json:"gas_wanted"`
	GasUsed   int64               `json:"gas_used"`
	Timestamp time.Time           `json:"timestamp"`
}
