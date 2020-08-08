package db

import (
	"hscan/schema"
)

// QueryLatestBlockHeight queries latest block height in database
func (db *Database) QueryLatestBlockHeight() (int64, error) {
	var block schema.Block
	if err := db.Order("height Desc").Limit(1).First(&block).Error; err != nil {
		if IsRecordNotFoundError(err) {
			return 0, nil
		}
		return -1, err
	}

	return block.Height, nil
}
