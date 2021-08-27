package server

import (
	"hscan/schema"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (s *Server) blocks(c *gin.Context) {

	height, _ := strconv.ParseInt(c.DefaultQuery("begin", "0"), 10, 64)
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
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

	// if err := s.db.Order("height DESC").Offset((page - 1) * iLimit).Limit(iLimit).Find(&blocks).Error; err != nil {
	// 	s.l.Printf("query blocks from db failed")
	// }

	if err := s.db.Order("height DESC").Where("height <= ?", height-(page-1)*iLimit).Limit(iLimit).Find(&blocks).Error; err != nil {
		s.l.Printf("query blocks from db failed")
	}

	s.blockResponse(c, total, blocks)
}

func (s *Server) block(c *gin.Context) {
	param := c.Param("param")
	var blocks []*schema.Block

	if err := s.db.Where("height = ? or block_hash = ?", param, param).First(&blocks).Error; err != nil {
		s.l.Printf("query blocks from db failed")
	} else {
		var txs []*schema.Transaction
		if err := s.db.Where("height = ?", param).Find(&txs).Error; err != nil {
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
	s.blockResponse(c, total, blocks)
}
