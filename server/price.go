package server

import (
	"encoding/json"
	"fmt"
	"hscan/models"
	"strconv"
	"strings"
	"time"
)

func (s *Server) SynchronismWallet() {
	for {
		result, err := s.client.QueryHscInfo()
		if err != nil {
			s.l.Printf("query Users of Number failed")
			time.Sleep(time.Duration(5) * time.Minute)
			continue
		}

		hscinfo := models.WalletHscInfo{}
		err = json.Unmarshal(result.Body(), &hscinfo)
		if err != nil {
			time.Sleep(time.Duration(5) * time.Minute)
			continue
		}

		if hscinfo.Code == 200 {

		}

		time.Sleep(time.Duration(5) * time.Minute)
	}
}

func (s *Server) GetDenom() {
	for {
		account, err := s.getAccount(s.Hschain.DestroyAddress)

		if err != nil {
			time.Sleep(time.Duration(5) * time.Minute)
			continue
		}

		fmt.Println(s.Hschain.DestroyAddress)
		accountinfo := models.Accountinfo{}
		fmt.Println(string(account.Body()))
		err = json.Unmarshal(account.Body(), &accountinfo)
		if err != nil {
			time.Sleep(time.Duration(5) * time.Minute)
			continue
		}

		coins := accountinfo.Result.Value.Coins

		for i := 0; i < len(coins); i++ {
			amount, _ := strconv.ParseInt(coins[i].Amount, 10, 64)
			s.Destory[coins[i].Denom] = amount
		}
		time.Sleep(time.Duration(5) * time.Minute)

	}
}

func (s *Server) updatePriceinto() {

	for {
		for k, v := range s.Priceinto {
			pirce, priceunit, err := s.queryDenomPrice(k)
			if err != nil {
				time.Sleep(time.Duration(5) * time.Minute)
				continue
			}
			v.Pirce = pirce.(string)
			v.Priceunit = priceunit.(string)
		}

		number, err := s.client.QueryUsersNumber()
		if err != nil {
			s.l.Printf("query Users of Number failed")
			time.Sleep(time.Duration(5) * time.Minute)
			continue
		}

		result := models.UsersNumber{}
		err = json.Unmarshal(number.Body(), &result)
		if err != nil {
			time.Sleep(time.Duration(5) * time.Minute)
			continue
		}

		if result.Code == 200 {
			s.UsersNumber = (int32)(result.Result.UsersNum)
			s.HeldByUsers, _ = strconv.ParseFloat(result.Result.HeldByUsers, 64)
		}

		time.Sleep(time.Duration(5) * time.Minute)
	}
}

func (s *Server) queryDenomPrice(denom interface{}) (interface{}, interface{}, error) {

	nom := strings.Replace(denom.(string), "u", "", 1)
	nom = strings.ToTitle(nom)
	if denom == "hst" || denom == "uhst" || denom == "hst0" || denom == "uhst0" {
		status, err := s.client.Queryexchangerate("hst_pri")
		if err != nil {
			return "0.00000", nil, err
		}

		result := models.HstExchangeRate{}
		err = json.Unmarshal(status.Body(), &result)
		if err != nil {
			return "0.00000", nil, err
		}
		if result.Code != 200 {
			return "0.00000", nil, fmt.Errorf("server get date error")
		}
		num := result.Result.HstPri
		return num, "/" + nom, nil
	} else {
		return "0.00000", "/" + nom, nil
	}

}

func (s *Server) denomPrice(denoms map[string]interface{}, denom string) (map[string]interface{}, error) {

	if Priceinto, OK := s.Priceinto[denom]; OK {
		denoms["price"] = Priceinto.Pirce
		denoms["priceunit"] = Priceinto.Priceunit

	} else {
		var Priceinto models.PriceInto
		num, priceunit, err := s.queryDenomPrice(denom)
		denoms["price"] = num
		denoms["priceunit"] = priceunit
		if err == nil {
			Priceinto.Pirce = num.(string)
			Priceinto.Priceunit = priceunit.(string)
			s.Priceinto[denom] = Priceinto

		}
		return denoms, err
	}
	return denoms, nil
}

func (s *Server) CoinsPrice(coins []map[string]interface{}) ([]map[string]interface{}, error) {

	coins = sortCoins(coins)
	var err error = nil
	for j := 0; j < len(coins); j++ {
		denom := coins[j]["denom"]
		coins[j], err = s.denomPrice(coins[j], denom.(string))

	}
	return coins, err
}
