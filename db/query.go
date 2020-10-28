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

func (db *Database) QueryLatestTxBlockHeight() (int64, error) {
	var txs schema.Transaction
	if err := db.Order("id Desc").Limit(1).First(&txs).Error; err != nil {
		if IsRecordNotFoundError(err) {
			return 0, nil
		}
		return -1, err
	}

	return int64(txs.ID), nil
}

func (db *Database) QueryAddressTxAcount(address string) (int64, error) {
	var txs []schema.Transaction

	if err := db.Order("id DESC").Where(" Sender = ? or Recipient = ?", address, address).Find(&txs).Error; err != nil {
		if IsRecordNotFoundError(err) {
			return 0, nil
		}
		return -1, err
	}
	return int64(len(txs)), nil
}
