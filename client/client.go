package client

import (
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
