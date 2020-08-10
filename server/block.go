package server

import (
	"hscan/schema"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) blocks(c *gin.Context) {

	limit := c.DefaultQuery("limit", "5")
	iLimit, _ := strconv.ParseInt(limit, 10, 64)
	if iLimit <= 0 {
		iLimit = 5
	}

	var blocks []*schema.Block

	if err := s.db.Order("height DESC").Limit(iLimit).Find(&blocks).Error; err != nil {
		s.l.Printf("query blocks from db failed")
	}

	c.JSON(http.StatusOK, gin.H{
		"paging": map[string]interface{}{
			"total":  1,
			"before": 2,
			"after":  3,
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

	c.JSON(http.StatusOK, gin.H{
		"paging": map[string]interface{}{
			"total":  1,
			"before": 2,
			"after":  3,
		},
		"data": blocks,
	})
}
