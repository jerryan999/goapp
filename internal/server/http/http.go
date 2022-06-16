package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bnkamalesh/webgo/v6"

	"github.com/gin-gonic/gin"
	"github.com/jerryan999/goapp/internal/api"
)

// HTTP struct holds all the dependencies required for starting HTTP server
type HTTP struct {
	server *http.Server
	cfg    *Config
}

// Start starts the HTTP server
func (h *HTTP) Start() {
	webgo.LOGHANDLER.Info("HTTP server, listening on", h.cfg.Host, h.cfg.Port)
	h.server.ListenAndServe()
}

// Config holds all the configuration required to start the HTTP server
type Config struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	DialTimeout  time.Duration
}

// NewService returns an instance of HTTP with all its dependencies set
func NewService(cfg *Config, a *api.API) (*HTTP, error) {
	h := &Handlers{
		api: a,
	}
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	router.GET("/health", h.Health)

	// User groups
	user_group := router.Group("/users")
	{
		user_group.POST("/create", h.CreateUser)
		user_group.GET("/read/:email", h.ReadUserByEmail)
	}

	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler:           router,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.ReadTimeout * 2,
	}

	return &HTTP{
		server: httpServer,
		cfg:    cfg,
	}, nil
}
