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
