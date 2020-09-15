package server

import (
	"encoding/json"
	"hscan/schema"
	"net/http"

	"github.com/gin-gonic/gin"
	resty "github.com/go-resty/resty/v2"
)

func (s *Server) parseResponse(response *resty.Response) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(response.Body(), &result)
	return result, err
}

func (s *Server) mintResponse(c *gin.Context, response *resty.Response) {
	var body map[string]interface{}
	if response == nil {
		body = nil
	} else {
		body, _ = s.parseResponse(response)
	}

	s.interfaceResponse(c, body)
}

func (s *Server) interfaceResponse(c *gin.Context, face interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": face,
	})
}

func (s *Server) blockResponse(c *gin.Context, total int64, blocks []*schema.Block) {

	if len(blocks) <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"paging": map[string]interface{}{
				"total": total,
				"end":   0,
				"begin": 0,
			},
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"paging": map[string]interface{}{
			"total": total,
			"end":   blocks[len(blocks)-1].Height,
			"begin": blocks[0].Height,
		},
		"data": blocks,
	})
}
