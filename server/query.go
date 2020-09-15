package server

import (
	"encoding/json"
	"hscan/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	resty "github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

func (s *Server) queryResponse(c *gin.Context, response *resty.Response) {
	var body map[string]interface{}
	if response == nil {
		body = nil
	} else {
		body, _ = s.parseResponse(response)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": body,
	})
}

func (s *Server) parseResponse(response *resty.Response) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(response.Body(), &result)
	return result, err
}

func (s *Server) getDenomPrice(denom interface{}) (interface{}, interface{}, error) {

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

func (s *Server) getAccountDenomPrice(response *resty.Response) (interface{}, error) {

	var result models.Accountinfo
	err := json.Unmarshal(response.Body(), &result)
	if err != nil {
		return nil, err
	}

	result.ArrangeInfo()

	for j := 0; j < len(result.Result.Value.Coins); j++ {
		denom := result.Result.Value.Coins[j]["denom"]
		if Priceinto, OK := s.Priceinto[denom]; OK {
			result.Result.Value.Coins[j]["price"] = Priceinto.Pirce
			result.Result.Value.Coins[j]["priceunit"] = Priceinto.Priceunit

		} else {
			var Priceinto models.Priceinto
			num, priceunit, err := s.getDenomPrice(denom)
			result.Result.Value.Coins[j]["price"] = num.(string)
			result.Result.Value.Coins[j]["priceunit"] = priceunit.(string)
			if err == nil {
				Priceinto.Pirce = num.(string)
				Priceinto.Priceunit = priceunit.(string)
				s.Priceinto[denom] = Priceinto
			}

		}

	}
	return result, nil
}

func (s *Server) getAccount(address string) (*resty.Response, error) {

	Account, err := s.client.QueryAccounts(address)
	if err != nil {
		s.l.Print(errors.Wrap(err, "failed to query the latest block height on the active network"))
		Account = nil
	}
	return Account, err
}

func (s *Server) account(c *gin.Context) {

	address := c.Param("address")
	Account, err := s.getAccount(address)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"data": nil,
		})
		return
	}

	body, err := s.getAccountDenomPrice(Account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": body,
	})
}

func (s *Server) mintingStatus(c *gin.Context) {

	status, err := s.client.Mintingstatus()
	if err != nil {
		s.l.Print(errors.Wrap(err, "failed to query the latest block height on the active network"))
		status = nil
	}
	s.queryResponse(c, status)
}

func (s *Server) mintingParams(c *gin.Context) {

	parameters, err := s.client.Mintingparameters()
	if err != nil {
		s.l.Print(errors.Wrap(err, "failed to query the latest block height on the active network"))
		parameters = nil
	}
	s.queryResponse(c, parameters)
}
