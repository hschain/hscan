package db

import (
	"fmt"
	"hscan/schema"
	"time"
)

func (db *Database) InsertScannedData(blocks []*schema.Block, txs []*schema.Transaction) error {

	for i := 0; i < len(txs); i++ {
		if err := db.Save(txs[i]).Error; err != nil {
			return err
		}
	}

	for i := 0; i < len(blocks); i++ {
		if err := db.Create(blocks[i]).Error; err != nil {
			return err
		}
	}

	return nil

}

func (db *Database) InsertScannedAlassetsData(Alassets []schema.PersonAlassets) error {

	fmt.Println(len(Alassets))
	for i := 0; i < len(Alassets); i++ {
		var txs []schema.PersonAlassets
		err := db.Order("id Desc").Where(" address = ? and denom = ?", Alassets[i].Address, Alassets[i].Denom).Limit(1).First(&txs).Error

		if err != nil || len(txs) <= 0 {
			tx := Alassets[i]
			if err := db.Save(&tx).Error; err != nil {
				return err
			}

			continue
		}

		if err := db.Model(&schema.PersonAlassets{}).Where(" address = ? and denom = ?", Alassets[i].Address, Alassets[i].Denom).Update("amount", Alassets[i].Amount).Error; err != nil {
			return err
		}
	}
	return nil

}

func (db *Database) Insertnodes(infos schema.NodeInfo) error {

	var txs []schema.NodeInfo
	err := db.Order("id Desc").Where(" node_name = ?", infos.NodeName).Limit(1).First(&txs).Error

	if err != nil || len(txs) <= 0 {
		if err := db.Save(&infos).Error; err != nil {
			return err
		}

		return nil
	}

	if err := db.Model(&schema.NodeInfo{}).Where(" node_name = ? ", infos.NodeName).Update("node_url", infos.NodeUrl).Error; err != nil {
		return err
	}

	return nil

}
func (db *Database) InsertVersionControl(infos schema.VersionControl) error {

	var txs []schema.VersionControl
	err := db.Order("id Desc").Where(" app = ? and platform = ?", infos.App, infos.Platform).Limit(1).First(&txs).Error

	if err != nil || len(txs) <= 0 {
		if err := db.Save(&infos).Error; err != nil {
			return err
		}

		return nil
	}

	if err := db.Model(&schema.VersionControl{}).Where(" app = ? and platform = ?", infos.App, infos.Platform).Updates(
		map[string]interface{}{
			"url":             infos.Url,
			"version":         infos.Version,
			"synchronization": infos.Synchronization,
		}).Error; err != nil {
		return err
	}

	return nil

}

func (db *Database) InsertUserVersion(infos schema.UserVersion) error {

	var txs []schema.UserVersion
	err := db.Order("id Desc").Where("address = ? and app = ? and platform = ?", infos.Address, infos.App, infos.Platform).Limit(1).First(&txs).Error

	if err != nil || len(txs) <= 0 {
		if err := db.Save(&infos).Error; err != nil {
			return err
		}

		return nil
	}

	if err := db.Model(&schema.UserVersion{}).Where("address = ? and app = ? and platform = ?", infos.Address, infos.App, infos.Platform).Updates(
		map[string]interface{}{
			"version":   infos.Version,
			"timestamp": time.Now(),
		}).Error; err != nil {
		return err
	}

	return nil

}
