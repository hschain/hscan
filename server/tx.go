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
	"github.com/pkg/errors"
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

func (s *Server) txs(c *gin.Context) {
	height, _ := strconv.ParseInt(c.DefaultQuery("begin", "0"), 10, 64)
	limit := c.DefaultQuery("limit", "5")
	timetable := c.DefaultQuery("timetable", "null")
	address := c.DefaultQuery("address", "null")
	iLimit, _ := strconv.ParseInt(limit, 10, 64)
	if iLimit <= 0 {
		iLimit = 5
	}

	total, err := s.db.QueryLatestTxBlockHeight()
	if total == -1 {
		s.l.Print(errors.Wrap(err, "failed to query the latest block height on the active network"))
	}

	if height <= 0 {
		height = total
	}
	var txs []*schema.Transaction
	if address == "null" {
		if err := s.db.Order("id DESC").Where(" id <= ?", height).Limit(iLimit).Find(&txs).Error; err != nil {
			s.l.Printf("query blocks from db failed")
		}
	} else {
		if err := s.db.Order("id DESC").Where(" id <= ? and (Sender = ? or Recipient = ?)", height, address, address).Limit(iLimit).Find(&txs).Error; err != nil {
			s.l.Printf("query blocks from db failed")
		}
		if timetable == "history" {
			s.db.Model(&schema.Transaction{}).Where("Sender = ? and sender_notice = 0", address, address).Update("sender_notice", 1)
			s.db.Model(&schema.Transaction{}).Where("Recipient = ? and RecipientNotice = 0", address, address).Update("RecipientNotice", 1)
		}

		Acount, err := s.db.QueryAddressTxAcount(address)
		if Acount == -1 {
			s.l.Print(errors.Wrap(err, "failed to query the latest block height on the active network"))
		}
		total = Acount
	}

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

	total, err := s.db.QueryLatestTxBlockHeight()
	if total == -1 {
		s.l.Print(errors.Wrap(err, "failed to query the latest block height on the active network"))
	}
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
