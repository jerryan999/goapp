package http

import (
	"fmt"
	"net/http"
	"time"

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
	h.server.ListenAndServe()
}

// Config holds all the configuration required to start the HTTP server
type Config struct {
	Host               string `json:"host"`
	Port               int    `json:"port"`
	ReadTimeoutSecond  int    `json:"read_timeout_second"`
	WriteTimeoutSecond int    `json:"write_timeout_second"`
	DialTimeoutSecond  int    `json:"dial_timeout_second"`
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
		user_group.GET("/retrieve", h.ReadUserByEmail)
	}

	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:           router,
		ReadTimeout:       time.Second * time.Duration(cfg.ReadTimeoutSecond),
		ReadHeaderTimeout: time.Second * time.Duration(cfg.ReadTimeoutSecond),
		WriteTimeout:      time.Second * time.Duration(cfg.WriteTimeoutSecond),
		IdleTimeout:       time.Second * time.Duration(cfg.DialTimeoutSecond),
	}

	return &HTTP{
		server: httpServer,
		cfg:    cfg,
	}, nil
}
