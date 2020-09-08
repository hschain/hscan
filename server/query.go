package server

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	resty "github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

func (s *Server) Hsresponse(c *gin.Context, response *resty.Response) {
	var body map[string]interface{}
	if response == nil {
		body = nil
	} else {
		body, _ = s.ParseResponse(response)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": body,
	})
}

func (s *Server) ParseResponse(response *resty.Response) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(response.Body(), &result)
	return result, err
}

func (s *Server) GetdenomPri(denom interface{}) (interface{}, interface{}, error) {

	if denom == "hst" || denom == "uhst" {
		status, err := s.client.Queryexchangerate("hst_pri")
		if err != nil {
			return "0.00000", nil, err
		}

		pri, err := s.ParseResponse(status)
		if err != nil {
			return "0.00000", nil, err
		}
		num := pri["result"].(map[string]interface{})["hst_pri"]
		return num, "$/hst", nil
	} else {
		return "0.00000", nil, nil
	}

}

func (s *Server) getaccountdenomPri(response *resty.Response) (map[string]interface{}, error) {

	var body map[string]interface{}
	if response == nil {
		body = nil
		return body, nil
	} else {
		body, _ = s.ParseResponse(response)
	}

	coinsMap := body["result"].(map[string]interface{})["value"].(map[string]interface{})["coins"]

	for j := 0; j < len(coinsMap.([]interface{})); j++ {

		coins := coinsMap.([]interface{})[j]
		denom := coins.(map[string]interface{})["denom"]

		num, priceunit, _ := s.GetdenomPri(denom)
		body["result"].(map[string]interface{})["value"].(map[string]interface{})["coins"].([]interface{})[j].(map[string]interface{})["price"] = num
		body["result"].(map[string]interface{})["value"].(map[string]interface{})["coins"].([]interface{})[j].(map[string]interface{})["priceunit"] = priceunit

	}
	return body, nil
}

func (s *Server) getaccount(address string) (*resty.Response, error) {

	Account, err := s.client.QueryAccounts(address)
	if err != nil {
		s.l.Print(errors.Wrap(err, "failed to query the latest block height on the active network"))
		Account = nil
	}
	return Account, err
}

func (s *Server) account(c *gin.Context) {

	address := c.Param("address")
	Account, _ := s.getaccount(address)
	body, _ := s.getaccountdenomPri(Account)
	c.JSON(http.StatusOK, gin.H{
		"data": body,
	})
}

func (s *Server) mintingstatus(c *gin.Context) {

	status, err := s.client.Mintingstatus()
	if err != nil {
		s.l.Print(errors.Wrap(err, "failed to query the latest block height on the active network"))
		status = nil
	}
	s.Hsresponse(c, status)
}

func (s *Server) mintingparams(c *gin.Context) {

	parameters, err := s.client.Mintingparameters()
	if err != nil {
		s.l.Print(errors.Wrap(err, "failed to query the latest block height on the active network"))
		parameters = nil
	}
	s.Hsresponse(c, parameters)
}
