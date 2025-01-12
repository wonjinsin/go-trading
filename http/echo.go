package http

import (
	"fmt"
	"log"
	"magmar/config"
	mw "magmar/middleware"
	"magmar/repository"
	"magmar/router"
	"magmar/service"
	"magmar/util"
	"os"

	"github.com/dimiro1/banner"
	"github.com/labstack/echo/v4"
)

var zlog *util.Logger

func init() {
	var err error
	zlog, err = util.NewLogger()
	if err != nil {
		log.Fatalf("InitLog module[main] err[%s]", err.Error())
		os.Exit(1)
	}

	zlog.Infow("logger started")
	// bannerInit()
}

// EchoHandler ...
func EchoHandler() *echo.Echo {
	magmar := config.Magmar
	e := echo.New()
	e.Use(mw.SetTRID())
	e.Use(mw.RequestLogger(zlog))
	e.Use(mw.Auth(magmar.GetString(util.HeaderXAPIKey)))
	e.HideBanner = true

	log.Println("Initializing repositories...")
	repo, err := repository.Init(magmar)
	if err != nil {
		fmt.Printf("Error when Start repository: %v\n", err)
		os.Exit(1)
	}

	log.Println("Initializing services...")
	svc, err := service.Init(magmar, repo)
	if err != nil {
		log.Printf("Error when Start service: %v\n", err)
		os.Exit(1)
	}

	router.Init(e, svc)
	return e
}

func bannerInit() {
	isEnabled := true
	isColorEnabled := true
	in, err := os.Open("banner.txt")
	if in == nil || err != nil {
		os.Exit(1)
	}

	banner.Init(os.Stdout, isEnabled, isColorEnabled, in)
}
