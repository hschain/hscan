package server

import (
	"encoding/json"
	"fmt"
	"hscan/schema"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	sdk "github.com/hschain/hschain/types"
)

func (s *Server) txresponse(c *gin.Context, total int64, txs []*schema.RavlTransaction) {

	if len(txs) <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"paging": map[string]interface{}{
				"total": total,
				"end":   0,
				"begin": 0,
			},
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"paging": map[string]interface{}{
			"total": total,
			"end":   txs[len(txs)-1].ID,
			"begin": txs[0].ID,
		},
		"data": txs,
	})
}

func (s *Server) formatRavlTransaction(txs []*schema.Transaction) []*schema.RavlTransaction {
	var Ravtxs []*schema.RavlTransaction
	Ravtxs = make([]*schema.RavlTransaction, 0)
	for j := 0; j < len(txs); j++ {

		tempTransaction := &schema.RavlTransaction{
			ID:        txs[j].ID,
			Height:    txs[j].Height,
			TxHash:    txs[j].TxHash,
			Messages:  txs[j].Messages,
			Memo:      txs[j].Memo,
			Timestamp: txs[j].Timestamp,
		}

		Ravtxs = append(Ravtxs, tempTransaction)
	}
	return Ravtxs
}

func (s *Server) format(txs []*schema.Transaction) {

	for i := range txs {
		var logs sdk.ABCIMessageLogs
		var messages []schema.Message
		s.cdc.UnmarshalJSON([]byte(txs[i].RawMessages), &logs)

		s.l.Printf("log is %+v", logs)

		for j := 0; j < len(logs); j++ {

			//convert
			msg := schema.Message{
				MsgIndex: logs[j].MsgIndex,
				Success:  logs[j].Success,
				Log:      logs[j].Log,
				Events:   make(map[string]map[string]string),
			}

			for k := 0; k < len(logs[j].Events); k++ {
				attrs := make(map[string]string)

				for l := 0; l < len(logs[j].Events[k].Attributes); l++ {
					if logs[j].Events[k].Attributes[l].Key == "amount" {
						if coin, err := sdk.ParseCoin(logs[j].Events[k].Attributes[l].Value); err == nil {
							attrs["amount"] = coin.Amount.String()
							attrs["denom"] = strings.ToUpper(coin.Denom)
							continue
						}

					}
					attrs[logs[j].Events[k].Attributes[l].Key] = logs[j].Events[k].Attributes[l].Value
				}

				msg.Events[logs[j].Events[k].Type] = attrs
			}

			messages = append(messages, msg)

		}

		txs[i].Messages = messages
	}

}

func (s *Server) gettxs(address, denom string, page, iLimit int64) ([]*schema.Transaction, error) {
	txs := make([]*schema.Transaction, 0)
	var sendTxs []*schema.Transaction
	var recipientTxs []*schema.Transaction
	if denom == "null" {

		if err := s.db.Order("id DESC").Where("(Sender = ? ) and (Sender <> ? and Recipient <> ?) ", address, s.Hschain.SupplementAddress, s.Hschain.SupplementAddress).Limit(iLimit).Find(&sendTxs).Error; err != nil {
			s.l.Printf("query blocks from db failed")
			return nil, err
		}
		if err := s.db.Order("id DESC").Where("(Recipient = ?) and (Sender <> ? and Recipient <> ?) ", address, s.Hschain.SupplementAddress, s.Hschain.SupplementAddress).Limit(iLimit).Find(&recipientTxs).Error; err != nil {
			s.l.Printf("query blocks from db failed")
			return nil, err
		}

		txs = append(txs, sendTxs...)
		txs = append(txs, recipientTxs...)
	} else {

		if err := s.db.Order("id DESC").Where("(Sender = ?) and (Sender <> ? and Recipient <> ?) and denom = ? and height <= ?", address, s.Hschain.SupplementAddress, s.Hschain.SupplementAddress, denom).Limit(iLimit).Find(&sendTxs).Error; err != nil {
			s.l.Printf("query blocks from db failed")
			return nil, err
		}
		if err := s.db.Order("id DESC").Where("(Recipient = ?) and (Sender <> ? and Recipient <> ?) and denom = ? and height <= ?", address, s.Hschain.SupplementAddress, s.Hschain.SupplementAddress, denom).Limit(iLimit).Find(&recipientTxs).Error; err != nil {
			s.l.Printf("query blocks from db failed")
			return nil, err
		}
		txs = append(txs, sendTxs...)
		txs = append(txs, recipientTxs...)
	}
	for i := 0; i < len(txs); i++ {
		for j := i + 1; j < len(txs); j++ {
			if txs[i].Height < txs[j].Height {
				tx := txs[i]
				txs[i] = txs[j]
				txs[j] = tx
			}
		}
	}

	if (int64)(len(txs)) < iLimit*(page-1) {
		return make([]*schema.Transaction, 0), nil
	}

	if (int64)(len(txs)) <= iLimit*page {
		return txs[iLimit*(page-1) : (int64)(len(txs))-iLimit*(page-1)], nil
	}

	limitTxs := txs[iLimit*(page-1) : iLimit]
	return limitTxs, nil
}

func (s *Server) txs(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit := c.DefaultQuery("limit", "5")
	//timetable := c.DefaultQuery("timetable", "null")
	address := c.DefaultQuery("address", "null")
	denom := c.DefaultQuery("denom", "null")
	txType, _ := strconv.ParseInt(c.DefaultQuery("type", "0"), 10, 64)
	iLimit, _ := strconv.ParseInt(limit, 10, 64)
	if iLimit <= 0 {
		iLimit = 5
	}

	var total = int64(s.cache.GetTotal(address, denom, int(txType)))
	var txs []*schema.Transaction

	ids := s.cache.GetTxids(address, denom, int(txType), page, iLimit)
	for i := len(ids) - 1; i >= 0; i-- {
		var tx schema.Transaction
		if err := s.db.Where("id = ?", ids[i]).First(&tx).Error; err == nil {
			txs = append([]*schema.Transaction{&tx}, txs...)
		}
	}
	/*
		if address == "null" {

			total, err = s.db.QueryTxBlockCount(s.Hschain.SupplementAddress)
			if total == -1 {
				s.l.Print(errors.Wrap(err, "failed to query the latest block height on the active network"))
			}
			if err := s.db.Order("id DESC").Offset((page-1)*iLimit).Where("(Sender <> ? and Recipient <> ?)", s.Hschain.SupplementAddress, s.Hschain.SupplementAddress).Limit(iLimit).Find(&txs).Error; err != nil {
				s.l.Printf("query blocks from db failed")
			}

		} else {

			if txType == 0 {
				txs, err = s.gettxs(address, denom, page, iLimit)
				if err != nil {
					s.l.Printf("query blocks from db failed")
				}
			} else if txType == 1 {
				if denom == "null" {
					if err := s.db.Order("id DESC").Offset((page-1)*iLimit).Where("(Sender = ? ) and (Sender <> ? and Recipient <> ?)", address, s.Hschain.SupplementAddress, s.Hschain.SupplementAddress).Limit(iLimit).Find(&txs).Error; err != nil {
						s.l.Printf("query blocks from db failed")
					}

				} else {
					if err := s.db.Order("id DESC").Offset((page-1)*iLimit).Where("(Sender = ?) and (Sender <> ? and Recipient <> ?) and denom = ?", address, s.Hschain.SupplementAddress, s.Hschain.SupplementAddress, denom).Limit(iLimit).Find(&txs).Error; err != nil {
						s.l.Printf("query blocks from db failed")
					}
				}
			} else if txType == 2 {
				if denom == "null" {
					if err := s.db.Order("id DESC").Offset((page-1)*iLimit).Where("(Recipient = ? ) and (Sender <> ? and Recipient <> ?)", address, s.Hschain.SupplementAddress, s.Hschain.SupplementAddress).Limit(iLimit).Find(&txs).Error; err != nil {
						s.l.Printf("query blocks from db failed")
					}

				} else {
					if err := s.db.Order("id DESC").Offset((page-1)*iLimit).Where("(Recipient = ?) and (Sender <> ? and Recipient <> ?) and denom = ?", address, s.Hschain.SupplementAddress, s.Hschain.SupplementAddress, denom).Limit(iLimit).Find(&txs).Error; err != nil {
						s.l.Printf("query blocks from db failed")
					}
				}
			}
			// if denom == "null" {
			// 	if err := s.db.Order("id DESC").Offset((page-1)*iLimit).Where("(Sender = ? or Recipient = ?) and (Sender <> ? and Recipient <> ?)", address, address, s.Hschain.SupplementAddress, s.Hschain.SupplementAddress).Limit(iLimit).Find(&txs).Error; err != nil {
			// 		s.l.Printf("query blocks from db failed")
			// 	}
			// } else {
			// 	if err := s.db.Order("id DESC").Offset((page-1)*iLimit).Where("(Sender = ? or Recipient = ?) and (Sender <> ? and Recipient <> ?) and denom = ?", address, address, s.Hschain.SupplementAddress, s.Hschain.SupplementAddress, denom).Limit(iLimit).Find(&txs).Error; err != nil {
			// 		s.l.Printf("query blocks from db failed")
			// 	}
			// }

			// if timetable == "null" {
			// 	if denom == "null" {
			// 		if err := s.db.Order("id DESC").Offset((page-1)*iLimit).Where("(Sender = ? or Recipient = ?) and (Sender <> ? and Recipient <> ?)", address, address, s.Hschain.SupplementAddress, s.Hschain.SupplementAddress).Limit(iLimit).Find(&txs).Error; err != nil {
			// 			s.l.Printf("query blocks from db failed")
			// 		}
			// 	} else {
			// 		if err := s.db.Order("id DESC").Offset((page-1)*iLimit).Where("(Sender = ? or Recipient = ?) and (Sender <> ? and Recipient <> ?) and denom = ?", address, address, s.Hschain.SupplementAddress, s.Hschain.SupplementAddress, denom).Limit(iLimit).Find(&txs).Error; err != nil {
			// 			s.l.Printf("query blocks from db failed")
			// 		}
			// 	}
			// }
			// if timetable == "history" {
			// 	if err := s.db.Order("id DESC").Offset((page-1)*iLimit).Where("(Sender = ? or Recipient = ?) and (Sender <> ? and Recipient <> ?) and denom = ?", address, address, s.Hschain.SupplementAddress, s.Hschain.SupplementAddress, denom).Limit(iLimit).Find(&txs).Error; err != nil {
			// 		s.l.Printf("query blocks from db failed")
			// 	}
			// 	s.db.Model(&schema.Transaction{}).Offset((page-1)*iLimit).Where("(Sender = ? and sender_notice = 0) and denom = ?", address, address, denom).Update("sender_notice", 1)
			// 	s.db.Model(&schema.Transaction{}).Offset((page-1)*iLimit).Where("(Recipient = ? and RecipientNotice = 0) and denom = ?", address, address, denom).Update("RecipientNotice", 1)
			// }
			// if timetable == "now" {
			// 	if err := s.db.Order("id DESC").Offset((page-1)*iLimit).Where("((Sender = ? and sender_notice = 0) or (Recipient = ? and RecipientNotice = 0)) and (Sender <> ? and Recipient <> ?)  and denom = ?", address, address, s.Hschain.SupplementAddress, s.Hschain.SupplementAddress, denom).Limit(iLimit).Find(&txs).Error; err != nil {
			// 		s.l.Printf("query blocks from db failed")
			// 	}
			// 	s.db.Model(&schema.Transaction{}).Offset((page-1)*iLimit).Where("and (Sender = ? and sender_notice = 0) and denom = ?", address, address, denom).Update("sender_notice", 1)
			// 	s.db.Model(&schema.Transaction{}).Offset((page-1)*iLimit).Where("(Recipient = ? and RecipientNotice = 0) and denom = ?", address, address, denom).Update("RecipientNotice", 1)
			// }

			/*
				Acount, err := s.db.QueryAddressTxAcount(address)
				if Acount == -1 {
					s.l.Print(errors.Wrap(err, "failed to query the latest block height on the active network"))
				}
				total = Acount

		}
	*/

	s.format(txs)
	Ravl := s.formatRavlTransaction(txs)
	s.txresponse(c, total, Ravl)
}

func (s *Server) tx(c *gin.Context) {
	txid := c.Param("txid")
	var txs []*schema.Transaction

	if err := s.db.Where("tx_hash = ?", txid).First(&txs).Error; err != nil {
		s.l.Printf("query blocks from db failed")
	}

	s.format(txs)
	Ravl := s.formatRavlTransaction(txs)

	//total, err := s.db.QueryTxBlockCount(s.Hschain.SupplementAddress)
	total := int64(s.cache.GetTotal("null", "null", 0))

	s.txresponse(c, total, Ravl)
}

func (s *Server) signedtx(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println(string(body))
	var bodymap interface{}
	if err := json.Unmarshal(body, &bodymap); err != nil {
		fmt.Println(err)
	}
	Ravl, _ := s.client.Signedtx(bodymap)
	s.mintResponse(c, Ravl)
}
