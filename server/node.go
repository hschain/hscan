package server

import (
	"encoding/json"
	"fmt"
	"hscan/models"
	"hscan/schema"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (s *Server) addNodes(c *gin.Context) {

	name := c.DefaultQuery("name", "hscan")
	url := c.DefaultQuery("url", "127.0.0.1")

	node := schema.NodeInfo{
		NodeName: name,
		NodeUrl:  url,
	}

	s.db.Insertnodes(node)

	c.JSON(http.StatusOK, gin.H{
		"message": map[string]interface{}{
			"type": "success",
		},
	})
}

func (s *Server) nodes(c *gin.Context) {
	var infos []*schema.NodeInfo

	if err := s.db.Select("*").Find(&infos).Error; err != nil {
		s.l.Printf("query blocks from db failed")
	}

	s.interfaceResponse(c, infos)
}

func (s *Server) tps(c *gin.Context) {
	var blocks []*schema.Block

	if err := s.db.Order("height DESC").Limit(1).Find(&blocks).Error; err != nil {
		s.l.Printf("query blocks from db failed")
	}
	tps := (float32)(blocks[0].NumTxs)/5.0 + 0.9
	s.interfaceResponse(c, (int)(tps))
}

func (s *Server) usersNumber(c *gin.Context) {

	status, err := s.client.QueryUsersNumber()
	if err != nil {
		s.l.Printf("query Users of Number failed")
		s.interfaceResponse(c, 0)
		return
	}
	s.interfaceResponse(c, status)
}

func (s *Server) frame(c *gin.Context) {

	status, err := s.client.Mintingstatus()
	if err != nil {
		s.l.Print(errors.Wrap(err, "failed to query the latest block height on the active network"))
		status = nil
	}
	var result map[string]interface{}
	err = json.Unmarshal(status.Body(), &result)
	currentDayProvisions := result["result"].(map[string]interface{})["status"].(map[string]interface{})["current_day_provisions"].(string)
	intCurrentDayProvisions, _ := strconv.ParseFloat(currentDayProvisions, 64)
	intTotalCirculationSupply := s.HeldByUsers

	var blocks []*schema.Block
	if err := s.db.Order("height DESC").Limit(1).Find(&blocks).Error; err != nil {
		s.l.Printf("query blocks from db failed")
	}
	tps := (float32)(blocks[0].NumTxs)/5.0 + 0.9

	frame := models.Frame{
		Tps:                    (int32)(tps),
		UsersNumber:            s.UsersNumber,
		CurrentDayProvisions:   intCurrentDayProvisions / 1000000,
		TotalCirculationSupply: (int64)(intTotalCirculationSupply),
	}
	fmt.Println(frame)
	s.interfaceResponse(c, frame)
}

func (s *Server) version(c *gin.Context) {
	var infos schema.VersionControl

	address := c.DefaultQuery("address", "null")
	app := c.DefaultQuery("app", "hscan")
	platform := c.DefaultQuery("platform", "pc")
	version := c.DefaultQuery("version", "v1.0.0")

	if address != "null" {
		Version := schema.UserVersion{
			Address:   address,
			App:       app,
			Platform:  platform,
			Version:   version,
			Timestamp: time.Now(),
		}
		s.db.InsertUserVersion(Version)
	}

	if err := s.db.Select("*").Where("app=? and platform=?", app, platform).Find(&infos).Error; err != nil {
		s.l.Printf("query blocks from db failed")
	}

	s.interfaceResponse(c, infos)

}

func (s *Server) addVersion(c *gin.Context) {
	app := c.DefaultQuery("app", "hscan")
	platform := c.DefaultQuery("platform", "pc")
	version := c.DefaultQuery("version", "pc")
	url := c.DefaultQuery("url", "127.0.0.1")
	synchronization := c.DefaultQuery("synchronization", "null")
	dbsynchronization := false
	if synchronization == "true" {
		dbsynchronization = true
	}
	Version := schema.VersionControl{
		App:             app,
		Platform:        platform,
		Version:         version,
		Url:             url,
		Synchronization: dbsynchronization,
	}

	s.db.InsertVersionControl(Version)

	c.JSON(http.StatusOK, gin.H{
		"message": map[string]interface{}{
			"type": "success",
		},
	})
}
