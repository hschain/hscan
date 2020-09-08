package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (s *Server) totals(c *gin.Context) {

	status, err := s.client.Querytotals()
	if err != nil {
		s.l.Print(errors.Wrap(err, "failed to query the latest block height on the active network"))
		status = nil
	}

	body, _ := s.ParseResponse(status)
	coinsMap := body["result"].([]interface{})
	for i := 0; i < len(coinsMap); i++ {
		denom := coinsMap[i].(map[string]interface{})["denom"]
		num, priceunit, _ := s.GetdenomPri(denom)
		body["result"].([]interface{})[i].(map[string]interface{})["price"] = num
		body["result"].([]interface{})[i].(map[string]interface{})["priceunit"] = priceunit
	}
	c.JSON(http.StatusOK, gin.H{
		"data": body,
	})
}

func (s *Server) total(c *gin.Context) {

	denomination := c.Param("denomination")
	parameters, err := s.client.Querytotal(denomination)
	if err != nil {
		s.l.Print(errors.Wrap(err, "failed to query the latest block height on the active network"))
		parameters = nil
	}
	body, _ := s.ParseResponse(parameters)

	num, priceunit, _ := s.GetdenomPri(denomination)
	body["price"] = num
	body["priceunit"] = priceunit

	c.JSON(http.StatusOK, gin.H{
		"data": body,
	})
}
