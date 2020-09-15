package server

import (
	"encoding/json"
	"hscan/models"

	"github.com/gin-gonic/gin"
	resty "github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

func (s *Server) getAccount(address string) (*resty.Response, error) {

	Account, err := s.client.QueryAccounts(address)
	if err != nil {
		s.l.Print(errors.Wrap(err, "failed to query the latest block height on the active network"))
		Account = nil
	}
	return Account, err
}

func (s *Server) getAccountDenomPrice(response *resty.Response) (interface{}, error) {

	var result models.AccountInfo
	err := json.Unmarshal(response.Body(), &result)
	if err != nil {
		return nil, err
	}

	result.Result.Value.Coins, _ = s.CoinsPrice(result.Result.Value.Coins)
	return result, nil
}

func (s *Server) account(c *gin.Context) {

	address := c.Param("address")
	Account, err := s.getAccount(address)
	if err != nil {
		s.interfaceResponse(c, nil)
		return
	}

	body, err := s.getAccountDenomPrice(Account)
	if err != nil {
		s.interfaceResponse(c, nil)
		return
	}

	s.interfaceResponse(c, body)
}
