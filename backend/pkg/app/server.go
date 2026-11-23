package app

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
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
func (s *Server) Run(ctx context.Context, host string, port int) error {
	srv := &http.Server{
		Addr:         host + ":" + strconv.Itoa(port),
		Handler:      s.Engine,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	g, ctx := errgroup.WithContext(ctx)
	
	// Start the server
	g.Go(func() error {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	
	// Wait for context cancellation or server shutdown
	<-ctx.Done()
	
	// Shutdown the server gracefully
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	
	return g.Wait()
}

// Healthz returns a simple health check endpoint.
func (s *Server) Healthz(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}