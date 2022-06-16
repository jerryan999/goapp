package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jerryan999/goapp/internal/api"
)

// Handlers struct has all the dependencies required for HTTP handlers
type Handlers struct {
	api *api.API
}

func (h *Handlers) Health(c *gin.Context) {
	d, err := h.api.Health()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, d)
}
