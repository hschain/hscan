package server

import (
	"hscan/schema"
	"net/http"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gin-gonic/gin"
)

func (s *Server) format(txs []*schema.Transaction) {

	for i := range txs {
		var logs sdk.ABCIMessageLogs
		var messages []*schema.Message
		s.cdc.UnmarshalJSON([]byte(txs[i].RawMessages), &logs)

		s.l.Printf("log is %+v", logs)

		for j := 0; j < len(logs); j++ {

			//convert
			msg := &schema.Message{
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
	limit := c.DefaultQuery("limit", "5")
	iLimit, _ := strconv.ParseInt(limit, 10, 64)
	if iLimit <= 0 {
		iLimit = 5
	}

	var txs []*schema.Transaction

	if err := s.db.Order("id DESC").Limit(iLimit).Find(&txs).Error; err != nil {
		s.l.Printf("query blocks from db failed")
	}

	s.format(txs)

	c.JSON(http.StatusOK, gin.H{
		"paging": map[string]interface{}{
			"total":  1,
			"before": 2,
			"after":  3,
		},
		"data": txs,
	})

}

func (s *Server) tx(c *gin.Context) {
	txid := c.Param("txid")
	var txs []*schema.Transaction
	var tx0 *schema.Transaction
	if err := s.db.Where("tx_hash = ?", txid).First(&txs).Error; err != nil {
		s.l.Printf("query blocks from db failed")
	}

	s.format(txs)
	if len(txs) == 1 {
		tx0 = txs[0]
	}

	c.JSON(http.StatusOK, gin.H{
		"paging": map[string]interface{}{
			"total":  1,
			"before": 2,
			"after":  3,
		},
		"data": tx0,
	})
}
