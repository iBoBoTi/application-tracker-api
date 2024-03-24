package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iBoBoTi/ats/internal/config"
	"github.com/iBoBoTi/ats/internal/log"
	"github.com/iBoBoTi/ats/internal/models"
	"github.com/iBoBoTi/ats/internal/security"
)

type Server struct {
	Router    *gin.Engine
	Config    config.Config
	Logger    log.Logger
	wg        sync.WaitGroup
	DB        *models.Database
	BuildTime string
	Version   string
	// Authenticator           *security.Authenticator
	TokenMaker security.Maker
}

func NewServer(cfg config.Config, db *models.Database, logger log.Logger) (*Server, error) {

	tokenMaker, err := security.NewPasetoMaker(cfg.TokenSymmetricKey)
	if err != nil {
		return nil, errors.New("cannot create a token maker")
	}

	server := &Server{
		Config:     cfg,
		Logger:     logger,
		DB:         db,
		TokenMaker: tokenMaker,
	}
	return server, nil
}

func (srv *Server) Start(address string) error {
	return srv.Router.Run(address)
}

func RunGinServer(srv *Server) {
	if srv == nil {
		srv.Logger.Fatal(errors.New("server instance cannot be nil"), nil)
	}

	httpServer := &http.Server{
		Addr:         srv.Config.HttpServerAddress,
		Handler:      srv.Router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// create a shutdown channel
	shutdownError := make(chan error)

	go func() {

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		srv.Logger.Info("shutting down server", map[string]interface{}{
			"signal": s.String(),
		})

		// create a context with a 5-seconds timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := httpServer.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		srv.Logger.Info("completing background tasks", nil)

		srv.wg.Wait()
		shutdownError <- nil
	}()

	// Likewise log a "starting server" message.
	srv.Logger.Info("starting server", map[string]interface{}{
		"addr": srv.Config.HttpServerAddress,
	})

	// Start the server as normal, returning any error.
	err := httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		srv.Logger.Fatal(err, nil)
	}

	err = <-shutdownError
	if err != nil {
		srv.Logger.Fatal(err, nil)
	}

	srv.Logger.Info("server stopped", nil)

}

func (srv *Server) GetConfig() config.Config {
	return srv.Config
}
