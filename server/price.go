package server

import (
	"encoding/json"
	"hscan/models"
	"strconv"
	"strings"
	"time"
)

func (s *Server) GetDenom() {
	for {
		Account, err := s.getAccount(s.Hschain.DestroyAddress)
		var Accountinfo models.Accountinfo
		err = json.Unmarshal(Account.Body(), &Accountinfo)
		if err != nil {
			time.Sleep(time.Duration(5) * time.Minute)
			continue
		}

		coins := Accountinfo.Result.Value.Coins

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
			Pirce, Priceunit, err := s.queryDenomPrice(k)
			if err != nil {
				time.Sleep(time.Duration(5) * time.Minute)
				continue
			}
			v.Pirce = Pirce.(string)
			v.Priceunit = Priceunit.(string)
		}

		Number, err := s.client.QueryUsersNumber()
		if err != nil {
			s.l.Printf("query Users of Number failed")
			time.Sleep(time.Duration(5) * time.Minute)
			continue
		}

		var result map[string]interface{}
		err = json.Unmarshal(Number.Body(), &result)
		if err != nil {
			time.Sleep(time.Duration(5) * time.Minute)
			continue
		}
		UsersNumber := result["result"].(map[string]interface{})["users_num"].(float64)
		held_by_users := result["result"].(map[string]interface{})["held_by_users"].(string)
		s.UsersNumber = (int32)(UsersNumber)
		s.Held_by_users, _ = strconv.ParseFloat(held_by_users, 64)
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

		pri, err := s.parseResponse(status)
		if err != nil {
			return "0.00000", nil, err
		}
		num := pri["result"].(map[string]interface{})["hst_pri"]
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

func (s *Server) CoinsPrice(Coins []map[string]interface{}) ([]map[string]interface{}, error) {

	Coins = sortCoins(Coins)
	var err error = nil
	for j := 0; j < len(Coins); j++ {
		denom := Coins[j]["denom"]
		Coins[j], err = s.denomPrice(Coins[j], denom.(string))

	}
	return Coins, err
}
