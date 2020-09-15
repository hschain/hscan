package schema

import "time"

// Block defines the schema for block information
type Block struct {
	Height        int64              `json:"height" gorm:"pk"`
	Proposer      string             `json:"proposer"`
	Moniker       string             `json:"moniker"`
	BlockHash     string             `json:"block_hash" gorm:"unique"`
	ParentHash    string             `json:"parent_hash"`
	NumPrecommits int64              `json:"num_pre_commits"`
	NumTxs        int64              `json:"num_txs"`
	TotalTxs      int64              `json:"total_txs"`
	Timestamp     time.Time          `json:"timestamp" gorm:"default:now()"`
	Txs           []*RavlTransaction `json:"txs" gorm:"-"`
	Denom         string             `json:"denom"`
	Amount        string             `json:"amount"`
}
