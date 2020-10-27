package websocket

import (
	"encoding/json"
	"fmt"
	"hscan/db"
	"hscan/schema"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// ClientManager is a websocket manager
type ClientManager struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	db         *db.Database
}

type AppMessage struct {
	Address string `json:"address,omitempty"`
	App     string `json:"app,omitempty"`
	Page    string `json:"page,omitempty"`
	Signal  string `json:"signal,omitempty"`
}

// Client is a websocket client
type Client struct {
	Id      string
	Address string
	Socket  *websocket.Conn
	Message *AppMessage
	Send    chan []byte
	Signal  chan []byte
}

// Manager define a ws server manager
var Manager = ClientManager{
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Clients:    make(map[string]*Client),
}

func Setdb(db *db.Database) {
	Manager.db = db
}

// Start is to start a ws server
func (manager *ClientManager) Start() {
	for {
		select {
		case conn := <-manager.Register:
			manager.Clients[conn.Address] = conn

		case conn := <-manager.Unregister:
			if _, ok := manager.Clients[conn.Address]; ok {
				close(conn.Send)
				delete(manager.Clients, conn.Address)
			}
		}
	}
}

func (manager *ClientManager) Send(Signal []byte, address string) {
	if conn, ok := manager.Clients[address]; ok {
		conn.Signal <- Signal
	}
}

func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		c.Socket.Close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			Manager.Unregister <- c
			c.Socket.Close()
			break
		}
		fmt.Println(string(message))
		var clientMessage *AppMessage
		err = json.Unmarshal(message, &clientMessage)

		if err != nil {
			continue
		}

		c.Message = clientMessage
		if clientMessage.Page == "app" && clientMessage.Signal == "connect" {
			c.Address = clientMessage.Address
			Manager.Register <- c
			c.Send <- []byte("{\"message\":\"The connection was successful\"}")
		} else if clientMessage.Page == "tx" && clientMessage.Signal == "in" {
			c.Nowtx()
		} else {
			c.Send <- []byte("{\"message\":\"Unknown request\"}")
		}
		fmt.Println(c)
		//jsonMessage, _ := json.Marshal(&Message{Sender: c.Id, Content: string(message)})
		//Manager.Broadcast <- jsonMessage
		//c.Send <- message

	}
}

func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func (c *Client) Push() {
	defer func() {
		c.Socket.Close()
	}()
	for {
		select {
		case signal, _ := <-c.Signal:
			if string(signal) == "tx" && c.Message.Signal == "in" && c.Message.Page == "tx" {
				c.Nowtx()
			}
		}
	}
}

func (c *Client) Nowtx() {
	var txs []*schema.Transaction
	if err := Manager.db.Order("id DESC").Where(" (Sender = ? and sender_notice = 0) or (Recipient = ? and RecipientNotice = 0)", c.Address, c.Address).Find(&txs).Error; err != nil {
		fmt.Printf("query blocks from db failed")
		return
	}

	tx := gin.H{
		"paging": map[string]interface{}{
			"total": len(txs),
			"end":   txs[len(txs)-1].ID,
			"begin": txs[0].ID,
		},
		"data": txs,
	}

	message, err := json.Marshal(tx)
	if err != nil {
		return
	}
	c.Send <- message

	if len(txs) > 0 {
		Manager.db.Model(&schema.Transaction{}).Where("Sender = ? and sender_notice = 0", c.Address, c.Address).Update("sender_notice", 1)
		Manager.db.Model(&schema.Transaction{}).Where("Recipient = ? and RecipientNotice = 0", c.Address, c.Address).Update("RecipientNotice", 1)
	}

}

func WsPage(c *gin.Context) {
	// change the reqest to websocket model
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)

	if error != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	// websocket connect
	client := &Client{
		Id:      uuid.NewV4().String(),
		Address: "",
		Socket:  conn,
		Send:    make(chan []byte),
		Signal:  make(chan []byte),
	}
	//Manager.Register <- client

	go client.Read()
	go client.Write()
	go Manager.Start()

}
