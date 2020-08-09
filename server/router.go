package server

import (
	"hscan/db"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/cosmos/cosmos-sdk/codec"
)

type Server struct {
	addr string
	e    *gin.Engine
	l    *log.Logger
	db   *db.Database
	cdc  *codec.Codec
}

func NewServer(addr string, l *log.Logger, db *db.Database, cdc *codec.Codec) *Server {
	return &Server{
		addr,
		gin.Default(),
		l,
		db,
		cdc,
	}
}

func (s *Server) cros(c *gin.Context) {
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("X-XSS-Protection", "1; mode=block")
	c.Header("Access-Control-Allow-Methods", "POST")
	c.Header("Access-Control-Allow-Headers", "x-auth-token, content-type")
	c.Header("Access-Control-Expose-Headers", "x-auth-token")
	c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-Frame-Options", "DENY")
	c.Header("Vary", "Origin")
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Header("Connection", "keep-alive")

}

func (s *Server) Start() error {
	s.l.Printf("web runnig at %s", s.addr)

	r := s.e.Group("/api/v1")
	r.Use(s.cros)

	r.GET("/blocks", s.blocks)
	r.GET("/txs", s.txs)

	s.e.Run(s.addr)
	return nil
}
