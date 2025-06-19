package server

import (
	"fmt"
	// "os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo *echo.Echo
	port int
}

// NewServer creates a new Server
func NewServer(port int) *Server {
	s := &Server{}

	if port <= 0 {
		panic("EASYINVESTING_PORT is a invalid number (port <= 0)")
	}
	s.port = port

	e := echo.New()
	s.echo = e

	s.echo.Use(getLoggerFunc())
	s.echo.Use(middleware.Recover())

	// frontendURL := os.Getenv("EASYINVESTING_FRONTEND_URL")
	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"*",
		},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
	}))

	s.setRoutes()

	return s
}

// Start starts the server using ´port´
func (s *Server) Start() {
	addr := fmt.Sprintf(":%d", s.port)
	err := s.echo.Start(addr)
	s.echo.Logger.Fatal(err)
}
