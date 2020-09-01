package scanner

import (
	"fmt"
	"log"
	"time"

	"hscan/client"
	"hscan/schema"

	"hscan/db"

	"github.com/zxs-paryada/hschain/codec"
	sdk "github.com/zxs-paryada/hschain/types"
	"github.com/zxs-paryada/hschain/x/auth/types"

	"github.com/pkg/errors"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// Scanner wraps the required params to scan blockchain
type Scanner struct {
	l      *log.Logger
	client *client.Client
	db     *db.Database
	cdc    *codec.Codec
}

// NewScanner returns Scanner
func NewScanner(l *log.Logger, client *client.Client, db *db.Database, cdc *codec.Codec) *Scanner {
	return &Scanner{
		l,
		client,
		db,
		cdc,
	}
}

// Start starts to synchronize blockchain data
func (s *Scanner) Start() error {
	go func() {
		for {
			s.l.Println("start - sync blockchain")
			err := s.sync()
			if err != nil {
				s.l.Printf("error - sync blockchain: %v\n", err)
			}
			s.l.Println("finish - sync blockchain")
			time.Sleep(time.Second)
		}
	}()

	return nil
}

// sync compares block height between the height saved in your database and
// the latest block height on the active chain and calls process to start ingesting data.
func (s *Scanner) sync() error {
	// Query latest block height saved in database
	dbHeight, err := s.db.QueryLatestBlockHeight()
	if dbHeight == -1 {
		s.l.Fatal(errors.Wrap(err, "failed to query the latest block height saved in database"))
	}

	// Query latest block height on the active network
	latestBlockHeight, err := s.client.LatestBlockHeight()
	if latestBlockHeight == -1 {
		s.l.Fatal(errors.Wrap(err, "failed to query the latest block height on the active network"))
	}

	// Synchronizing blocks from the scratch will return 0 and will ingest accordingly.
	// Skip the first block since it has no pre-commits
	if dbHeight == 0 {
		dbHeight = 1
	}

	//dbHeight = 11240

	s.l.Printf("dbHeight is %v, latestBlockHeight is %v", dbHeight, latestBlockHeight)

	// Ingest all blocks up to the latest height
	for i := dbHeight + 1; i <= latestBlockHeight; i++ {
		err := s.process(i)

		if err != nil {
			return err
		}
		s.l.Printf("synced block %d/%d \n", i, latestBlockHeight)
	}

	return nil
}

// process ingests chain data, such as block, transaction, validator set information
// and save them in database
func (s *Scanner) process(height int64) error {
	s.l.Printf("start process block %v", height)

	//Get block info from blockchain
	block, err := s.client.GetBlock(height)
	if err != nil {
		return fmt.Errorf("failed to query block using rpc client: %s", err)
	}

	//handle the block info
	s.l.Printf("block is %+v", *block)
	schemaBlock, err := s.getBlock(block)

	if err != nil {
		return fmt.Errorf("failed to get block: %s", err)
	}

	//Get txs in the block from blockchain
	txs, err := s.client.GetTxs(block)

	if err != nil {
		return fmt.Errorf("failed to get txs: %s", err)
	}

	//handle the txs
	transactions, err := s.getTxs(txs, block)
	if err != nil {
		return fmt.Errorf("failed to get schema txs: %s", err)
	}

	for i, trx := range transactions {
		s.l.Printf("transactions[%d] is %+v", i, *trx)
	}

	err = s.db.InsertScannedData(schemaBlock, transactions)
	if err != nil {
		return fmt.Errorf("failed to insert scanned data: %s", err)
	}

	return nil
}

// getBlock parses block information and wrap into Block schema struct
func (s *Scanner) getBlock(block *tmctypes.ResultBlock) ([]*schema.Block, error) {
	blocks := make([]*schema.Block, 0)

	tempBlock := &schema.Block{
		Height:    block.Block.Height,
		Proposer:  block.Block.ProposerAddress.String(),
		Moniker:   "super node",
		BlockHash: block.BlockMeta.BlockID.Hash.String(),
		//BlockHash:     block.BlockID.Hash.String(),
		ParentHash:    block.BlockMeta.Header.LastBlockID.Hash.String(),
		NumPrecommits: int64(len(block.Block.LastCommit.Precommits)),
		NumTxs:        block.Block.NumTxs,
		TotalTxs:      block.Block.TotalTxs,
		Timestamp:     block.Block.Time,
	}

	blocks = append(blocks, tempBlock)

	return blocks, nil
}

// getTxs parses transactions and wrap into Transaction schema struct
func (s *Scanner) getTxs(txs []*tmctypes.ResultTx, resBlock *tmctypes.ResultBlock) ([]*schema.Transaction, error) {
	transactions := make([]*schema.Transaction, 0)
	for i := range txs {
		var stdTx types.StdTx

		err := s.cdc.UnmarshalBinaryLengthPrefixed(txs[i].Tx, &stdTx)
		if err != nil {
			return nil, err
		}

		//s.l.Printf("stdTx is %+v", stdTx)

		resp := sdk.NewResponseResultTx(txs[i], stdTx, resBlock.Block.Time.Format(time.RFC3339))

		msgsBz, err := s.cdc.MarshalJSON(resp.Logs)
		if err != nil {
			return nil, err
		}

		tempTransaction := &schema.Transaction{
			Height:      resp.Height,
			TxHash:      resp.TxHash,
			Code:        resp.Code, // 0 is success
			RawMessages: string(msgsBz),
			Fee:         string(stdTx.Fee.Bytes()),
			Memo:        stdTx.GetMemo(),
			GasWanted:   resp.GasWanted,
			GasUsed:     resp.GasUsed,
			Timestamp:   resBlock.Block.Time,
		}

		transactions = append(transactions, tempTransaction)
	}

	return transactions, nil
}
