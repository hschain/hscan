package server

import (
	"hscan/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (s *Server) totals(c *gin.Context) {

	status, err := s.client.Querytotals()
	if err != nil {
		s.l.Print(errors.Wrap(err, "failed to query the latest block height on the active network"))
		status = nil
		c.JSON(http.StatusOK, gin.H{
			"data": nil,
		})
		return
	}

	body, _ := s.parseResponse(status)

	for i := 0; i < len(body["result"].([]interface{})); {
		denom := body["result"].([]interface{})[i].(map[string]interface{})["denom"]

		if denom == "syscoin" || denom == "SYSCOIN" {
			body["result"] = append(body["result"].([]interface{})[:i], body["result"].([]interface{})[i+1:]...)

			continue
		}

		if Priceinto, OK := s.Priceinto[denom.(string)]; OK {
			body["result"].([]interface{})[i].(map[string]interface{})["price"] = Priceinto.Pirce
			body["result"].([]interface{})[i].(map[string]interface{})["priceunit"] = Priceinto.Priceunit

		} else {
			var Priceinto models.Priceinto
			num, priceunit, err := s.getDenomPrice(denom)
			body["result"].([]interface{})[i].(map[string]interface{})["price"] = num.(string)
			body["result"].([]interface{})[i].(map[string]interface{})["priceunit"] = priceunit.(string)
			if err == nil {
				Priceinto.Pirce = num.(string)
				Priceinto.Priceunit = priceunit.(string)
				s.Priceinto[denom.(string)] = Priceinto
			}

		}
		i++
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
		c.JSON(http.StatusOK, gin.H{
			"data": nil,
		})
		return
	}
	body, _ := s.parseResponse(parameters)

	if Priceinto, OK := s.Priceinto[denomination]; OK {
		body["price"] = Priceinto.Pirce
		body["priceunit"] = Priceinto.Priceunit

	} else {
		var Priceinto models.Priceinto
		num, priceunit, err := s.getDenomPrice(denomination)
		body["price"] = num.(string)
		body["priceunit"] = priceunit.(string)
		if err == nil {
			Priceinto.Pirce = num.(string)
			Priceinto.Priceunit = priceunit.(string)
			s.Priceinto[denomination] = Priceinto
		}

	}

	c.JSON(http.StatusOK, gin.H{
		"data": body,
	})
}
