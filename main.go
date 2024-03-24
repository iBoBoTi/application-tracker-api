package main

import (
	"expvar"
	"io"
	"log"
	"os"

	"github.com/iBoBoTi/ats/config"
	logger "github.com/iBoBoTi/ats/log"
	"github.com/iBoBoTi/ats/models"
	"github.com/iBoBoTi/ats/routers"
	"github.com/iBoBoTi/ats/server"
)

// build time variables
var (
	buildTime string
	version   string
)

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
