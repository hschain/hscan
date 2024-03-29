package server

import (
	"encoding/json"
	"hscan/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func sortCoins(Coins []map[string]interface{}) []map[string]interface{} {

	for i := 0; i < len(Coins); {

		if Coins[i]["denom"] == "syscoin" || Coins[i]["denom"] == "SYSCOIN" {
			Coins = append(Coins[:i], Coins[i+1:]...)

			continue
		}

		if Coins[i]["denom"] == "hst" || Coins[i]["denom"] == "uhst" {
			a := Coins[i]
			Coins[i] = Coins[0]
			Coins[0] = a
		}
		i++
	}

	if len(Coins) == 0 {
		hst := make(map[string]interface{}, 1)
		hst["amount"] = "0"
		hst["denom"] = "uhst"
		Coins = make([]map[string]interface{}, 1)
		Coins[0] = hst
		return Coins
	}

	if Coins[0]["denom"] != "uhst" {
		hst := make([]map[string]interface{}, 1)
		hst[0] = make(map[string]interface{}, 1)
		hst[0]["amount"] = "0"
		hst[0]["denom"] = "uhst"
		Coins = append(hst, Coins...)
		return Coins
	}
	return Coins
}

func (s *Server) totals(c *gin.Context) {

	status, err := s.client.Querytotals()
	if err != nil {
		s.l.Print(errors.Wrap(err, "failed to query the latest block height on the active network"))
		s.interfaceResponse(c, nil)
		return
	}

	var result models.TotalInfo
	err = json.Unmarshal(status.Body(), &result)
	if err != nil {
		s.interfaceResponse(c, nil)
		return
	}

	result.Result, _ = s.CoinsPrice(result.Result)
	for i := 0; i < len(result.Result); i++ {
		if result.Result[i]["denom"].(string) == "uhst" {
			result.Result[i]["minted_supply"] = (int64)(s.HeldByUsers * 1000000)
		} else if result.Result[i]["denom"].(string) == "uhsc" || result.Result[i]["denom"].(string) == "hsc" {
			amount := result.Result[i]["amount"].(string)
			IntAmount, _ := strconv.ParseInt(amount, 10, 64)
			result.Result[i]["minted_supply"] = IntAmount - (int64)(s.HeldByHsc)
		} else {
			amount := result.Result[i]["amount"].(string)
			denom := result.Result[i]["denom"].(string)
			IntAmount, _ := strconv.ParseInt(amount, 10, 64)
			if _, OK := s.Destory[denom]; OK {
				result.Result[i]["minted_supply"] = IntAmount - s.Destory[denom]
			} else {
				result.Result[i]["minted_supply"] = IntAmount
			}
		}
	}
	s.interfaceResponse(c, result)
}

func (s *Server) total(c *gin.Context) {

	denomination := c.Param("denomination")
	parameters, err := s.client.Querytotal(denomination)
	if err != nil {
		s.l.Print(errors.Wrap(err, "failed to query the latest block height on the active network"))
		s.interfaceResponse(c, nil)
		return
	}

	body, _ := s.parseResponse(parameters)
	body, _ = s.denomPrice(body, denomination)

	s.interfaceResponse(c, body)
}
