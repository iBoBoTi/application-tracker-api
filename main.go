package main

import (
	"expvar"
	"io"
	"log"
	"os"

	"github.com/iBoBoTi/ats/internal/config"
	logger "github.com/iBoBoTi/ats/internal/log"
	"github.com/iBoBoTi/ats/internal/models"
	"github.com/iBoBoTi/ats/routers"
	"github.com/iBoBoTi/ats/server"

	"github.com/unidoc/unipdf/v3/common/license"
)

// build time variables
var (
	buildTime string
	version   string
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	cfg, err := config.Load(".")
	if err != nil {
		log.Fatal(err)
	}
	err = license.SetMeteredKey(cfg.UNIDOC_LICENSE_API_KEY)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	expvar.NewString("version").Set(version)

	cfg, err := config.Load(".")
	if err != nil {
		log.Fatal(err)
	}

	var logWriter io.Writer = os.Stdout

	db := models.GetDB(cfg)

	zeroLogger := logger.NewZeroLogger(logWriter, logger.LevelInfo)
	srv, err := server.NewServer(cfg, db, zeroLogger)
	if err != nil {
		log.Fatal(err)
	}

	srv.BuildTime = buildTime
	srv.Version = version

	routers.SetupRouter(srv)
	server.RunGinServer(srv)

}
