package models

type PriceInto struct {
	Pirce     string `json:"pirce"`
	Priceunit string `json:"priceunit"`
}

type AccountInfo struct {
	Height string `json:"height"`
	Result struct {
		Type  string `json:"type"`
		Value struct {
			Address   string                   `json:"address"`
			Coins     []map[string]interface{} `json:"coins"`
			PublicKey struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"public_key"`
			AccountNumber string `json:"account_number"`
			Sequence      string `json:"sequence"`
		} `json:"value"`
	} `json:"result"`
}

type TotalInfo struct {
	Height string                   `json:"height"`
	Result []map[string]interface{} `json:"result"`
}
