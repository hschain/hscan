package server

import (
	"hscan/db"
	"log"

	"github.com/gin-gonic/gin"

	"hscan/client"
	"hscan/models"
	"hscan/websocket"

	"github.com/hschain/hschain/codec"
)

type Server struct {
	addr          string
	e             *gin.Engine
	l             *log.Logger
	db            *db.Database
	cdc           *codec.Codec
	client        *client.Client
	Priceinto     map[string]models.PriceInto
	UsersNumber   int32
	Held_by_users float64
}

func NewServer(addr string, l *log.Logger, db *db.Database, cdc *codec.Codec, client *client.Client) *Server {
	return &Server{
		addr,
		gin.Default(),
		l,
		db,
		cdc,
		client,
		make(map[string]models.PriceInto, 1),
		0,
		0,
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

	go s.updatePriceinto()

	s.l.Printf("web runnig at %s", s.addr)

	r := s.e.Group("/api/v1")
	r.Use(s.cros)

	r.GET("/ws", websocket.WsPage)
	websocket.Setdb(s.db)

	r.GET("/tps", s.tps)
	r.GET("/nodes", s.nodes)
	r.GET("/addnodes", s.addNodes)
	r.GET("/frame", s.frame)
	r.GET("/usersnumber", s.usersNumber)
	r.GET("/version", s.version)
	r.GET("/addversion", s.addVersion)

	r.GET("/blocks", s.blocks)
	r.GET("/blocks/:param", s.block)

	r.GET("/txs", s.txs)
	r.GET("/txs/:txid", s.tx)
	r.POST("/txs", s.signedtx)

	r.GET("/total", s.totals)
	r.GET("/total/:denomination", s.total)

	r.GET("/topaccounts", s.getTopAccounts)
	r.GET("/account/:address", s.account)

	r.GET("/minting/status", s.mintingStatus)
	r.GET("/minting/params", s.mintingParams)

	s.e.Run(s.addr)

	return nil
}
