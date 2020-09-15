package models

type Priceinto struct {
	Pirce     string `json:"pirce"`
	Priceunit string `json:"priceunit"`
}

type Accountinfo struct {
	Height string `json:"height"`
	Result struct {
		Type  string `json:"type"`
		Value struct {
			Address   string              `json:"address"`
			Coins     []map[string]string `json:"coins"`
			PublicKey struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"public_key"`
			AccountNumber string `json:"account_number"`
			Sequence      string `json:"sequence"`
		} `json:"value"`
	} `json:"result"`
}

func (info *Accountinfo) ArrangeInfo() {

	for i := 0; i < len(info.Result.Value.Coins); {

		if info.Result.Value.Coins[i]["denom"] == "syscoin" || info.Result.Value.Coins[i]["denom"] == "SYSCOIN" {
			info.Result.Value.Coins = append(info.Result.Value.Coins[:i], info.Result.Value.Coins[i+1:]...)

			continue
		}

		if info.Result.Value.Coins[i]["denom"] == "hst" || info.Result.Value.Coins[i]["denom"] == "uhst" {
			a := info.Result.Value.Coins[i]
			info.Result.Value.Coins[i] = info.Result.Value.Coins[0]
			info.Result.Value.Coins[0] = a
		}
		i++
	}

	if len(info.Result.Value.Coins) == 0 {
		hst := make(map[string]string, 1)
		hst["amount"] = "0"
		hst["denom"] = "uhst"
		info.Result.Value.Coins = make([]map[string]string, 1)
		info.Result.Value.Coins[0] = hst
		return
	}

	if info.Result.Value.Coins[0]["denom"] != "uhst" {
		hst := make([]map[string]string, 1)
		hst[0]["amount"] = "0"
		hst[0]["denom"] = "uhst"
		info.Result.Value.Coins = append(hst, info.Result.Value.Coins...)
		return
	}
}
