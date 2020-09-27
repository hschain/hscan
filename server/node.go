package server

import (
	"hscan/schema"

	"github.com/gin-gonic/gin"
)

func (s *Server) nodes(c *gin.Context) {
	var infos []*schema.NodeInfo

	if err := s.db.Select("*").Find(&infos).Error; err != nil {
		s.l.Printf("query blocks from db failed")
	}

	s.interfaceResponse(c, infos)
}

func (s *Server) tps(c *gin.Context) {
	var blocks []*schema.Block

	if err := s.db.Order("height DESC").Limit(1).Find(&blocks).Error; err != nil {
		s.l.Printf("query blocks from db failed")
	}
	tps := (float32)(blocks[0].NumTxs)/5.0 + 0.9
	s.interfaceResponse(c, (int)(tps))
}
