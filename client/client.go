package client

import (
	"fmt"
	"strconv"
	"time"

	"hscan/config"

	resty "github.com/go-resty/resty/v2"
	rpccli "github.com/tendermint/tendermint/rpc/client"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// Client wraps for both Tendermint RPC and other API clients that
// are needed for this project
type Client struct {
	lcdClient *resty.Client
	rpcClient *rpccli.HTTP
	cfg       *config.NodeConfig
}

// NewClient creates a new client with the given config
func NewClient(cfg config.NodeConfig) *Client {

	lcdClient := resty.New().
		SetHostURL(cfg.LCDServerEndpoint).
		SetTimeout(30 * time.Second)

	rpcClient := rpccli.NewHTTP(cfg.NodeServerEndPoint, "/websocket")

	return &Client{
		lcdClient,
		rpcClient,
		&cfg,
	}
}

//LatestBlockHeight
func (c *Client) LatestBlockHeight() (int64, error) {
	status, err := c.rpcClient.Status()

	if err != nil {
		return -1, err
	}

	return status.SyncInfo.LatestBlockHeight, nil
}

func (c *Client) GetBlock(height int64) (*tmctypes.ResultBlock, error) {
	return c.rpcClient.Block(&height)
}

// Txs queries for all the transactions in a block.
// It uses `Tx` RPC method to query for the transaction
func (c *Client) GetTxs(block *tmctypes.ResultBlock) ([]*tmctypes.ResultTx, error) {
	txs := make([]*tmctypes.ResultTx, len(block.Block.Txs), len(block.Block.Txs))

	for i, tmTx := range block.Block.Txs {

		//log.Printf("tx is %X", tmTx.Hash())

		tx, err := c.rpcClient.Tx(tmTx.Hash(), true)
		if err != nil {
			return nil, err
		}

		txs[i] = tx
	}

	return txs, nil
}

func (c *Client) RestyGet(path string, param string) (*resty.Response, error) {
	var data string = path + param
	fmt.Println(data)
	return resty.New().R().EnableTrace().Get(data)
}

func (c *Client) RestyPost(path string, param interface{}) (*resty.Response, error) {
	resp := resty.New().R()
	resp.SetHeader("Content-Type", "application/json")
	m := param.(map[string]interface{})
	resp.SetBody(m)
	return resp.Post(path)
}

func (c *Client) QueryAccounts(address string) (*resty.Response, error) {

	return c.RestyGet(c.cfg.LCDServerEndpoint+"/auth/accounts/", address)
}

func (c *Client) Mintingparameters() (*resty.Response, error) {

	return c.RestyGet(c.cfg.LCDServerEndpoint+"/minting/parameters", "")
}

func (c *Client) Mintingstatus() (*resty.Response, error) {

	return c.RestyGet(c.cfg.LCDServerEndpoint+"/minting/status", "")
}

func (c *Client) Mintingbonus(Height int64) (*resty.Response, error) {

	height := strconv.FormatInt(Height, 10)
	return c.RestyGet(c.cfg.LCDServerEndpoint+"/minting/bonus/", height)
}

func (c *Client) Signedtx(parameters interface{}) (*resty.Response, error) {

	return c.RestyPost(c.cfg.LCDServerEndpoint+"/txs", parameters)
}

func (c *Client) Querytotal(address string) (*resty.Response, error) {

	return c.RestyGet(c.cfg.LCDServerEndpoint+"/supply/total/", address)
}

func (c *Client) Querytotals() (*resty.Response, error) {

	return c.RestyGet(c.cfg.LCDServerEndpoint+"/supply/total", "")
}

func (c *Client) Queryexchangerate(denom string) (*resty.Response, error) {

	return c.RestyGet(c.cfg.PriServerEndpoint+"/h5/", denom)
}

func (c *Client) QueryUsersNumber() (*resty.Response, error) {

	return c.RestyGet(c.cfg.PriServerEndpoint+"/h5/hsc_users_num", "")
}
