package server

import (
	"hscan/schema"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (s *Server) response(c *gin.Context, total int64, blocks []*schema.Block) {

	if len(blocks) <= 0 {
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
			"end":   blocks[len(blocks)-1].Height,
			"begin": blocks[0].Height,
		},
		"data": blocks,
	})
}

func (s *Server) blocks(c *gin.Context) {

	height, _ := strconv.ParseInt(c.DefaultQuery("begin", "0"), 10, 64)
	limit := c.DefaultQuery("limit", "5")
	iLimit, _ := strconv.ParseInt(limit, 10, 64)
	if iLimit <= 0 {
		iLimit = 5
	}

	total, err := s.db.QueryLatestBlockHeight()
	if total == -1 {
		s.l.Print(errors.Wrap(err, "failed to query the latest block height on the active network"))
	}

	if height <= 0 {
		height = total
	}

	var blocks []*schema.Block

	if err := s.db.Order("height DESC").Where(" height <= ?", height).Limit(iLimit).Find(&blocks).Error; err != nil {
		s.l.Printf("query blocks from db failed")
	}
	s.response(c, total, blocks)
}

func (s *Server) block(c *gin.Context) {
	height := c.Param("height")
	var blocks []*schema.Block

	if err := s.db.Where("height = ?", height).First(&blocks).Error; err != nil {
		s.l.Printf("query blocks from db failed")
	} else {
		var txs []*schema.Transaction
		if err := s.db.Where("height = ?", height).Find(&txs).Error; err != nil {
			s.l.Printf("query txs from db failed")
		}

		if len(txs) > 0 {
			s.format(txs)
			Ravl := s.formatRavlTransaction(txs)
			blocks[0].Txs = Ravl
		}

	}

	total, err := s.db.QueryLatestBlockHeight()
	if total == -1 {
		s.l.Print(errors.Wrap(err, "failed to query the latest block height on the active network"))
	}
	s.response(c, total, blocks)
}
