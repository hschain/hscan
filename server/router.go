package server

import (
	"hscan/db"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"hscan/client"
	"hscan/models"

	"github.com/zxs-paryada/hschain/codec"
)

type Server struct {
	addr      string
	e         *gin.Engine
	l         *log.Logger
	db        *db.Database
	cdc       *codec.Codec
	client    *client.Client
	Priceinto map[string]models.Priceinto
}

func NewServer(addr string, l *log.Logger, db *db.Database, cdc *codec.Codec, client *client.Client) *Server {
	return &Server{
		addr,
		gin.Default(),
		l,
		db,
		cdc,
		client,
		make(map[string]models.Priceinto, 1),
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

func (s *Server) updatePriceinto() {

	for {
		for k, v := range s.Priceinto {
			Pirce, Priceunit, err := s.getDenomPrice(k)
			if err != nil {
				continue
			}
			v.Pirce = Pirce.(string)
			v.Priceunit = Priceunit.(string)
		}
		time.Sleep(time.Duration(5) * time.Minute)
	}
}

func (s *Server) Start() error {
	s.l.Printf("web runnig at %s", s.addr)

	r := s.e.Group("/api/v1")
	r.Use(s.cros)

	r.GET("/blocks", s.blocks)
	r.GET("/blocks/:param", s.block)
	r.GET("/txs", s.txs)
	r.GET("/txs/:txid", s.tx)
	r.GET("/total", s.totals)
	r.GET("/total/:denomination", s.total)
	r.GET("/account/:address", s.account)
	r.GET("/minting/status", s.mintingStatus)
	r.GET("/minting/params", s.mintingParams)
	r.POST("/txs", s.signedtx)
	s.e.Run(s.addr)

	go s.updatePriceinto()
	return nil
}
