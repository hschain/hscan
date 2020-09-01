package server

import (
	"hscan/schema"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	sdk "github.com/zxs-paryada/hschain/types"
)

func (s *Server) format(txs []*schema.Transaction) {

	for i := range txs {
		var logs sdk.ABCIMessageLogs
		var messages []*schema.Message
		s.cdc.UnmarshalJSON([]byte(txs[i].RawMessages), &logs)

		s.l.Printf("log is %+v", logs)

		for j := 0; j < len(logs); j++ {
			index := 0

			//put message event in the head
			for k := 0; k < len(logs[j].Events); k++ {
				if logs[j].Events[k].Type == "message" {
					for l := 0; l < len(logs[j].Events[k].Attributes); l++ {
						if logs[j].Events[k].Attributes[l].Key == "action" {
							txs[i].Action = logs[j].Events[k].Attributes[l].Value
							break
						}
					}
					index = k
					break
				}
			}

			if index != 0 {
				var arr sdk.StringEvents
				arr = append(arr, logs[j].Events[index])
				arr = append(arr, logs[j].Events[0:index]...)
				if index != len(logs[j].Events) {
					arr = append(arr, logs[j].Events[index+1:]...)
				}
				logs[j].Events = arr
			}

			//convert
			msg := &schema.Message{
				MsgIndex: logs[j].MsgIndex,
				Success:  logs[j].Success,
				Log:      logs[j].Log,
			}

			for k := 0; k < len(logs[j].Events); k++ {
				evt := struct {
					Type       string      `json:"type"`
					Attributes interface{} `json:"attributes"`
				}{
					Type:       logs[j].Events[k].Type,
					Attributes: make(map[string]string),
				}

				for l := 0; l < len(logs[j].Events[k].Attributes); l++ {
					if logs[j].Events[k].Attributes[l].Key == "amount" {
						if coin, err := sdk.ParseCoin(logs[j].Events[k].Attributes[l].Value); err == nil {
							evt.Attributes.(map[string]string)["amount"] = coin.Amount.String()
							evt.Attributes.(map[string]string)["denom"] = strings.ToUpper(coin.Denom)
							continue
						}

					}
					evt.Attributes.(map[string]string)[logs[j].Events[k].Attributes[l].Key] = logs[j].Events[k].Attributes[l].Value
				}

				msg.Events = append(msg.Events, evt)
			}

			messages = append(messages, msg)

		}

		txs[i].Messages = messages
	}

}

func (s *Server) txs(c *gin.Context) {
	height := c.DefaultQuery("height", "0")
	limit := c.DefaultQuery("limit", "5")
	iLimit, _ := strconv.ParseInt(limit, 10, 64)
	if iLimit <= 0 {
		iLimit = 5
	}

	var txs []*schema.Transaction

	total, err := s.db.QueryLatestTxBlockHeight()
	if total == -1 {
		s.l.Fatal(errors.Wrap(err, "failed to query the latest block height on the active network"))
	}

	if err := s.db.Order("height DESC").Where(" height >= ?", height).Limit(iLimit).Find(&txs).Error; err != nil {
		s.l.Printf("query blocks from db failed")
	}

	s.format(txs)

	c.JSON(http.StatusOK, gin.H{
		"paging": map[string]interface{}{
			"total":  total,
			"before": txs[len(txs)-1].Height,
			"after":  txs[0].Height,
		},
		"data": txs,
	})

}

func (s *Server) tx(c *gin.Context) {
	txid := c.Param("txid")
	var txs []*schema.Transaction

	total, err := s.db.QueryLatestTxBlockHeight()
	if total == -1 {
		s.l.Fatal(errors.Wrap(err, "failed to query the latest block height on the active network"))
	}

	if err := s.db.Where("tx_hash = ?", txid).First(&txs).Error; err != nil {
		s.l.Printf("query blocks from db failed")
	} else {

	}

	s.format(txs)

	c.JSON(http.StatusOK, gin.H{
		"paging": map[string]interface{}{
			"total":  total,
			"before": txs[len(txs)-1].Height,
			"after":  txs[0].Height,
		},
		"data": txs,
	})
}
