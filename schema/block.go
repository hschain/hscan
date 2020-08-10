package schema

import "time"

// Block defines the schema for block information
type Block struct {
	Height        int64          `json:"height" gorm:"pk;not null"`
	Proposer      string         `json:"proposer" gorm:"not null"`
	Moniker       string         `json:"moniker" gorm:"not null"`
	BlockHash     string         `json:"block_hash" gorm:"not null;unique"`
	ParentHash    string         `json:"parent_hash" gorm:"not null"`
	NumPrecommits int64          `json:"num_pre_commits" gorm:"not null"`
	NumTxs        int64          `json:"num_txs" gorm:"default:0"`
	TotalTxs      int64          `json:"total_txs" gorm:"default:0"`
	Timestamp     time.Time      `json:"timestamp" gorm:"default:now()"`
	Txs           []*Transaction `json:"txs" gorm:"-"`
}
