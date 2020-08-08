package scanner

import (
	"fmt"
	"log"
	"time"

	"hscan/client"
	"hscan/schema"

	"hscan/db"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
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

	for {
		select {}
	}
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

	//for test
	dbHeight = 11306

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
	_, err = s.getTxs(txs, block)
	if err != nil {
		return fmt.Errorf("failed to get schema txs: %s", err)
	}

	return fmt.Errorf("failed to insert scanned data: %s", "ok")

	err = s.db.InsertScannedData(schemaBlock)
	if err != nil {
		return fmt.Errorf("failed to insert scanned data: %s", err)
	}

	/*
		valSet, err := s.client.ValidatorSet(block.Block.LastCommit.Height())
		if err != nil {
			return fmt.Errorf("failed to query validator set using rpc client: %s", err)
		}

		vals, err := s.client.Validators()
		if err != nil {
			return fmt.Errorf("failed to query validators using rpc client: %s", err)
		}

		resultBlock, err := s.getBlock(block) // TODO: Reward Fees Calculation
		if err != nil {
			return fmt.Errorf("failed to get block: %s", err)
		}

		resultTxs, err := s.getTxs(block)
		if err != nil {
			return fmt.Errorf("failed to get transactions: %s", err)
		}

		resultValidators, err := s.getValidators(vals)
		if err != nil {
			return fmt.Errorf("failed to get validators: %s", err)
		}

		resultPreCommits, err := s.getPreCommits(block.Block.LastCommit, valSet)
		if err != nil {
			return fmt.Errorf("failed to get precommits: %s", err)
		}

		err = ex.db.InsertExportedData(resultBlock, resultTxs, resultValidators, resultPreCommits)
		if err != nil {
			return fmt.Errorf("failed to insert exporterd data: %s", err)
		}

	*/
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
	parseTx := func(cdc *codec.Codec, txBytes []byte) (sdk.Tx, error) {
		var tx types.StdTx

		err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
		if err != nil {
			return nil, err
		}

		return tx, nil
	}

	format := func(cdc *codec.Codec, resTx *tmctypes.ResultTx, resBlock *tmctypes.ResultBlock) (*sdk.TxResponse, error) {
		tx, err := parseTx(cdc, resTx.Tx)
		if err != nil {
			s.l.Printf("parseTx failed")
			return nil, err
		}

		resp := sdk.NewResponseResultTx(resTx, tx, resBlock.Block.Time.Format(time.RFC3339))
		return &resp, nil
	}

	var err error
	out := make([]*sdk.TxResponse, len(txs))
	for i := range txs {
		s.l.Printf("raw tx is %s", txs[i].Tx)
		out[i], err = format(s.cdc, txs[i], resBlock)
		if err != nil {
			return nil, err
		}
	}

	s.l.Printf("Txs is %+v", out)

	return nil, nil
}

/*

// getTxs parses transactions and wrap into Transaction schema struct
func (s *s) getTxs(block *models.Block) ([]*schema.Transaction, error) {
	transactions := make([]*schema.Transaction, 0)

	txs, err := ex.client.Txs(block)
	if err != nil {
		return nil, err
	}

	if len(txs) > 0 {
		for _, tx := range txs {
			var stdTx txtypes.StdTx
			s.cdc.UnmarshalBinaryLengthPrefixed([]byte(tx.Tx), &stdTx)

			msgsBz, err := s.cdc.MarshalJSON(stdTx.GetMsgs())
			if err != nil {
				return nil, err
			}

			sigs := make([]types.Signature, len(stdTx.Signatures), len(stdTx.Signatures))

			for i, sig := range stdTx.Signatures {
				consPubKey, err := ctypes.Bech32ifyConsPub(sig.PubKey)
				if err != nil {
					return nil, err
				}

				sigs[i] = types.Signature{
					Address:       sig.Address().String(), // hex string
					AccountNumber: sig.AccountNumber,
					Pubkey:        consPubKey,
					Sequence:      sig.Sequence,
					Signature:     base64.StdEncoding.EncodeToString(sig.Signature), // encode base64
				}
			}

			sigsBz, err := s.cdc.MarshalJSON(sigs)
			if err != nil {
				return nil, err
			}

			tempTransaction := &schema.Transaction{
				Height:     tx.Height,
				TxHash:     tx.Hash.String(),
				Code:       tx.TxResult.Code, // 0 is success
				Messages:   string(msgsBz),
				Signatures: string(sigsBz),
				Memo:       stdTx.Memo,
				GasWanted:  tx.TxResult.GasWanted,
				GasUsed:    tx.TxResult.GasUsed,
				Timestamp:  block.Block.Time,
			}

			transactions = append(transactions, tempTransaction)

		}
	}

	return transactions, nil
}
*/
