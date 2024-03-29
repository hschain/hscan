package db

import (
	"hscan/schema"
)

func (db *Database) QueryBlockCount() (int64, error) {
	var count int64
	if err := db.Model(&schema.Block{}).Count(&count).Error; err != nil {
		if IsRecordNotFoundError(err) {
			return 0, nil
		}
		return -1, err
	}

	return count, nil
}

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

func (db *Database) QueryTxBlockCount(SupplementAddress string) (int64, error) {
	var count int64

	if err := db.Model(&schema.Transaction{}).Where("(Sender <> ? and Recipient <> ?)", SupplementAddress, SupplementAddress).Count(&count).Error; err != nil {
		if IsRecordNotFoundError(err) {
			return 0, nil
		}
		return -1, err
	}

	return count, nil
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
	var count int64
	if err := db.Model(&schema.Transaction{}).Where(" Sender = ? or Recipient = ?", address, address).Count(&count).Error; err != nil {
		if IsRecordNotFoundError(err) {
			return 0, nil
		}
		return -1, err
	}
	return count, nil
}
