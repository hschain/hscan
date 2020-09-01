package server

import (
	"hscan/schema"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (s *Server) blocks(c *gin.Context) {

	height, _ := strconv.ParseInt(c.DefaultQuery("height", "0"), 10, 64)
	limit := c.DefaultQuery("limit", "60")
	iLimit, _ := strconv.ParseInt(limit, 10, 64)
	if iLimit <= 0 {
		iLimit = 5
	}

	if height < 0 {
		height = 0
	}

	var blocks []*schema.Block

	if err := s.db.Order("height DESC").Where(" height >= ?", height).Limit(iLimit).Find(&blocks).Error; err != nil {
		s.l.Printf("query blocks from db failed")
	}

	total, err := s.db.QueryLatestBlockHeight()
	if total == -1 {
		s.l.Fatal(errors.Wrap(err, "failed to query the latest block height on the active network"))
	}

	c.JSON(http.StatusOK, gin.H{
		"paging": map[string]interface{}{
			"total":  total,
			"before": blocks[len(blocks)-1].Height,
			"after":  blocks[0].Height,
		},
		"data": blocks,
	})

}

func (s *Server) block(c *gin.Context) {
	height := c.Param("height")
	var blocks []*schema.Block

	if err := s.db.Where("height = ?", height).First(&blocks).Error; err != nil {
		s.l.Printf("query blocks from db failed")
	}

	total, err := s.db.QueryLatestBlockHeight()
	if total == -1 {
		s.l.Fatal(errors.Wrap(err, "failed to query the latest block height on the active network"))
	}

	c.JSON(http.StatusOK, gin.H{
		"paging": map[string]interface{}{
			"total":  total,
			"before": blocks[len(blocks)-1].Height,
			"after":  blocks[0].Height,
		},
		"data": blocks,
	})
}
