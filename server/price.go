package server

import (
	"hscan/models"
	"strings"
	"time"
)

func (s *Server) updatePriceinto() {

	for {
		for k, v := range s.Priceinto {
			Pirce, Priceunit, err := s.queryDenomPrice(k)
			if err != nil {
				continue
			}
			v.Pirce = Pirce.(string)
			v.Priceunit = Priceunit.(string)
		}
		time.Sleep(time.Duration(5) * time.Minute)
	}
}

func (s *Server) queryDenomPrice(denom interface{}) (interface{}, interface{}, error) {

	nom := strings.Replace(denom.(string), "u", "", 1)
	nom = strings.ToTitle(nom)
	if denom == "hst" || denom == "uhst" {
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
