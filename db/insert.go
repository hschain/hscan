package db

import (
	"hscan/schema"
)

func (db *Database) InsertScannedData(blocks []*schema.Block, txs []*schema.Transaction) error {

	for i := 0; i < len(txs); i++ {
		if err := db.Save(txs[i]).Error; err != nil {
			return err
		}
	}

	for i := 0; i < len(blocks); i++ {
		if err := db.Save(blocks[i]).Error; err != nil {
			return err
		}
	}

	return nil

}
