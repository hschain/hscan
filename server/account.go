package server

import (
	"encoding/json"
	"hscan/models"
	"hscan/schema"
	"net/http"
	"strconv"

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

func (s *Server) getTopAccounts(c *gin.Context) {
	var Alassets []schema.PersonAlassets
	limit := c.DefaultQuery("limit", "5")
	page := c.DefaultQuery("page", "1")
	denom := c.DefaultQuery("denom", "uhst")

	ilimit, _ := strconv.ParseInt(limit, 10, 64)
	if ilimit <= 0 {
		ilimit = 0
	}

	ipage, _ := strconv.ParseInt(page, 10, 64)
	if ipage <= 0 {
		ipage = 0
	}

	if err := s.db.Order("denom Desc").Where("denom = ?", denom).Limit(500).Find(&Alassets).Error; err != nil {
		s.l.Printf("query blocks from db failed")
	}

	Begim := (ipage - 1) * ilimit
	End := ipage * ilimit
	if End > int64(len(Alassets)) {
		End = int64(len(Alassets))
	}

	TopAccounts := Alassets[Begim:End]

	c.JSON(http.StatusOK, gin.H{
		"paging": map[string]interface{}{
			"total": len(Alassets),
			"end":   End,
			"begin": Begim,
		},
		"data": TopAccounts,
	})

}
