package models

type UsersNumber struct {
	Code   int `json:"code"`
	Result struct {
		UsersNum    int    `json:"users_num"`
		HeldByUsers string `json:"held_by_users"`
	} `json:"result"`
}

type HstExchangeRate struct {
	Code   int `json:"code"`
	Result struct {
		HstPri string `json:"hst_pri"`
	} `json:"result"`
}

type WalletHscInfo struct {
	Code   int `json:"code"`
	Result struct {
		TotalHst        string `json:"total_hst"`
		HeldByUsers     string `json:"held_by_users"`
		Destroyed       string `json:"destroyed"`
		Produced        string `json:"produced"`
		Sum             string `json:"sum"`
		SumRate         string `json:"sum_rate"`
		TodExptProd     string `json:"tod_expt_prod"`
		NextCycExptProd string `json:"next_cyc_expt_prod"`
		NumOfPool       int    `json:"num_of_pool"`
		NumOfRig        int    `json:"num_of_rig"`
		TodCmptPowCoe   string `json:"tod_cmpt_pow_coe"`
		LstCmptPowCoe   string `json:"lst_cmpt_pow_coe"`
		CurProd         string `json:"cur_prod"`
		CurUse          string `json:"cur_use"`
		RedRate         string `json:"red_rate"`
		ExptNextCpe     string `json:"expt_next_cpe"`
		Mggs            string `json:"mggs"`
		TotRigPurUsdt   string `json:"tot_rig_pur_usdt"`
		DestroyHashrate string `json:"destroy_hashrate"`
	} `json:"result"`
}
