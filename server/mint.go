package server

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (s *Server) mintingStatus(c *gin.Context) {

	status, err := s.client.Mintingstatus()
	if err != nil {
		s.l.Print(errors.Wrap(err, "failed to query the latest block height on the active network"))
		status = nil
	}
	s.mintResponse(c, status)
}

func (s *Server) mintingParams(c *gin.Context) {

	parameters, err := s.client.Mintingparameters()
	if err != nil {
		s.l.Print(errors.Wrap(err, "failed to query the latest block height on the active network"))
		parameters = nil
	}
	s.mintResponse(c, parameters)
}
