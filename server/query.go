package server

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	resty "github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

func (s *Server) ParseResponse(response *resty.Response) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(response.Body(), &result)
	return result, err
}

func (s *Server) account(c *gin.Context) {
	address := c.Param("address")

	Account, err := s.client.QueryAccounts(address)
	if err != nil {
		s.l.Fatal(errors.Wrap(err, "failed to query the latest block height on the active network"))
	}

	body, _ := s.ParseResponse(Account)
	c.JSON(http.StatusOK, gin.H{
		"data": body,
	})
}

func (s *Server) mintingstatus(c *gin.Context) {

	Account, err := s.client.Mintingstatus()
	if err != nil {
		s.l.Fatal(errors.Wrap(err, "failed to query the latest block height on the active network"))
	}

	body, _ := s.ParseResponse(Account)
	c.JSON(http.StatusOK, gin.H{
		"data": body,
	})
}

func (s *Server) mintingparams(c *gin.Context) {

	Account, err := s.client.Mintingparameters()
	if err != nil {
		s.l.Fatal(errors.Wrap(err, "failed to query the latest block height on the active network"))
	}

	body, _ := s.ParseResponse(Account)
	c.JSON(http.StatusOK, gin.H{
		"data": body,
	})
}
