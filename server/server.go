package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/TxChar/isekai-shop-api/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type echoServer struct {
	app *echo.Echo
	db *gorm.DB
	conf *config.Config
}

var (
	once sync.Once
	server *echoServer
)

func NewEchoServer(conf *config.Config, db *gorm.DB) *echoServer {
	echoApp := echo.New()
	echoApp.Logger.SetLevel(log.DEBUG)

	once.Do(func() {
			server = &echoServer{
				app: echoApp,
				db: db,
				conf: conf,
			}
		})

	return server
}

func (s *echoServer) Start() {
	// Initialize all middlewares
	corsMiddleware := getCORSMiddleware(s.conf.Server.AllowOrigins)
	bodyLimitMiddleware := getBodyLimitMiddleware(s.conf.Server.BodyLimit)
	timeOutMiddleware := getTimeOutMiddleware(s.conf.Server.Timeout)

	// Prevent application from crashing
	s.app.Use(middleware.Recover())

	s.app.Use(middleware.Logger())
	s.app.Use(timeOutMiddleware)
	s.app.Use(corsMiddleware)
	s.app.Use(bodyLimitMiddleware)
	
	s.app.GET("/v1/health", s.healthCheck)
	s.httpListening()
}

func (s *echoServer) httpListening() {
	url := fmt.Sprintf(":%d", s.conf.Server.Port)

	if err := s.app.Start(url); err != nil && err != http.ErrServerClosed {
		s.app.Logger.Fatalf("Error: %s", err.Error())
	}
}

func (s *echoServer) healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

// func getLoggerMiddleware() echo.MiddlewareFunc {
// 	return middleware.Logger()
// }

func getTimeOutMiddleware(timeout time.Duration) echo.MiddlewareFunc {
	return middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper: middleware.DefaultSkipper,
		ErrorMessage: "Request Timeout",
		Timeout: timeout * time.Second,
	})
}

func getCORSMiddleware(allowOrigins []string) echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: allowOrigins,
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	})
}

func getBodyLimitMiddleware(bodyLimit string) echo.MiddlewareFunc {
	return middleware.BodyLimit(bodyLimit)
}