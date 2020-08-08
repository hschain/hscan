package models

//tmtypes "github.com/tendermint/tendermint/types"

//Block by tm
/*
type Block struct {
	BlockMeta tmtypes.BlockMeta `json:"block_meta"`
	Block     tmtypes.Block     `json:"block"`
}
*/

/*
type Block struct {
	BlockMeta struct {
		Header struct {
			ChainID     string          `json:"chain_id"`
			Height      decimal.Decimal `json:"height"`
			Time        time.Time       `json:"time"`
			NumTxs      decimal.Decimal `json:"num_txs"`
			LastBlockID struct {
				Hash  string `json:"hash"`
				Parts struct {
					Total decimal.Decimal `json:"total"`
					Hash  string          `json:"hash"`
				} `json:"parts"`
			} `json:"last_block_id"`
			TotalTxs           decimal.Decimal `json:"total_txs"`
			LastCommitHash     string          `json:"last_commit_hash"`
			DataHash           string          `json:"data_hash"`
			ValidatorsHash     string          `json:"validators_hash"`
			NextValidatorsHash string          `json:"next_validators_hash"`
			ConsensusHash      string          `json:"consensus_hash"`
			AppHash            string          `json:"app_hash"`
			LastResultsHash    string          `json:"last_results_hash"`
			EvidenceHash       string          `json:"evidence_hash"`
			ProposerAddress    string          `json:"proposer_address"`
			Version            struct {
				Block decimal.Decimal `json:"block"`
				App   decimal.Decimal `json:"app"`
			} `json:"version"`
		} `json:"header"`
		BlockID struct {
			Hash  string `json:"hash"`
			Parts struct {
				Total decimal.Decimal `json:"total"`
				Hash  string          `json:"hash"`
			} `json:"parts"`
		} `json:"block_id"`
	} `json:"block_meta"`
	Block struct {
		Header struct {
			ChainID     string          `json:"chain_id"`
			Height      decimal.Decimal `json:"height"`
			Time        time.Time       `json:"time"`
			NumTxs      decimal.Decimal `json:"num_txs"`
			LastBlockID struct {
				Hash  string `json:"hash"`
				Parts struct {
					Total decimal.Decimal `json:"total"`
					Hash  string          `json:"hash"`
				} `json:"parts"`
			} `json:"last_block_id"`
			TotalTxs           decimal.Decimal `json:"total_txs"`
			LastCommitHash     string          `json:"last_commit_hash"`
			DataHash           string          `json:"data_hash"`
			ValidatorsHash     string          `json:"validators_hash"`
			NextValidatorsHash string          `json:"next_validators_hash"`
			ConsensusHash      string          `json:"consensus_hash"`
			AppHash            string          `json:"app_hash"`
			LastResultsHash    string          `json:"last_results_hash"`
			EvidenceHash       string          `json:"evidence_hash"`
			ProposerAddress    string          `json:"proposer_address"`
			Version            struct {
				Block decimal.Decimal `json:"block"`
				App   decimal.Decimal `json:"app"`
			} `json:"version"`
		} `json:"header"`
		Data struct {
			Txs []string `json:"txs"`
		} `json:"data"`
		Evidence struct {
			Evidence []string `json:"evidence"`
		} `json:"evidence"`
		LastCommit struct {
			BlockID struct {
				Hash  string `json:"hash"`
				Parts struct {
					Total decimal.Decimal `json:"total"`
					Hash  string          `json:"hash"`
				} `json:"parts"`
			} `json:"block_id"`
			Precommits []struct {
				ValidatorAddress string          `json:"validator_address"`
				ValidatorIndex   string          `json:"validator_index"`
				Height           string          `json:"height"`
				Round            string          `json:"round"`
				Timestamp        time.Time       `json:"timestamp"`
				Type             decimal.Decimal `json:"type"`
				BlockID          struct {
					Hash  string `json:"hash"`
					Parts struct {
						Total decimal.Decimal `json:"total"`
						Hash  string          `json:"hash"`
					} `json:"parts"`
				} `json:"block_id"`
				Signature string `json:"signature"`
			} `json:"precommits"`
		} `json:"last_commit"`
	} `json:"block"`
}
*/
