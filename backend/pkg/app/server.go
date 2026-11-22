package app

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Server is the single binary entrypoint.
type Server struct {
	Engine *gin.Engine
}

// New creates the server from the core.
func New(core *CoreContext) (*Server, error) {
	return &Server{Engine: core.Engine}, nil
}

// Run serves the engine on host:port.
func (s *Server) Run(host string, port int) error {
	srv := &http.Server{Addr: host + ":" + strconv.Itoa(port), Handler: s.Engine}
	return srv.ListenAndServe()
}
