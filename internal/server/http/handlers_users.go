package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jerryan999/goapp/internal/users"
)

// CreateUser is the HTTP handler to create a new user
// This handler does not use any framework, instead just the standard library
func (h *Handlers) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()
	u := new(users.User)
	if err := c.BindJSON(u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.api.CreateUser(ctx, u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// ReadUserByEmail is the HTTP handler to read an existing user by email
func (h *Handlers) ReadUserByEmail(c *gin.Context) {
	ctx := c.Request.Context()
	email := c.Param("email")

	u, err := h.api.ReadUserByEmail(ctx, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, u)
}
