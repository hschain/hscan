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

type HschainInfo struct {
	Height string `json:"height"`
	Result Result `json:"result"`
}
type MintPlans struct {
	Period         string `json:"period"`
	TotalPerPeriod string `json:"total_per_period"`
	TotalPerDay    string `json:"total_per_day"`
}
type BlockProvision struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}
type DestoryAmount struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}
type ConversionRates struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}
type Status struct {
	TotalMintedSupply       string            `json:"total_minted_supply"`
	TotalMintingSupply      string            `json:"total_minting_supply"`
	TotalDistrSupply        string            `json:"total_distr_supply"`
	TotalCirculationSupply  string            `json:"total_circulation_supply"`
	CurrentDayProvisions    string            `json:"current_day_provisions"`
	NextPeriodDayProvisions string            `json:"next_period_day_provisions"`
	NextPeroidStartTime     string            `json:"next_peroid_startTime"`
	BlockProvision          BlockProvision    `json:"block_provision"`
	BurnAmount              []interface{}     `json:"burn_amount"`
	DestoryAmount           []DestoryAmount   `json:"destory_amount"`
	ConversionRates         []ConversionRates `json:"conversion_rates"`
}
type Result struct {
	MintPlans []MintPlans `json:"mint_plans"`
	Status    Status      `json:"status"`
}
