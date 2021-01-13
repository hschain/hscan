package models

type Frame struct {
	UsersNumber            int32   `json:"usersNumber"`
	Tps                    int32   `json:"tps"`
	CurrentDayProvisions   float64 `json:"current_day_provisions"`
	TotalCirculationSupply int64   `json:"total_circulation_supply"`
}

type Accountinfo struct {
	Height string `json:"height"`
	Result struct {
		Type  string `json:"type"`
		Value struct {
			Address string `json:"address"`
			Coins   []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"coins"`
			PublicKey     interface{} `json:"public_key"`
			AccountNumber string      `json:"account_number"`
			Sequence      string      `json:"sequence"`
		} `json:"value"`
	} `json:"result"`
}
