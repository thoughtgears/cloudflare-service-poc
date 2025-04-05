package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/thoughtgears/cloudflare-tunnels-poc/services"
)

// DeleteUser handles HTTP DELETE requests to the /users/:id endpoint.
// It extracts the user ID from the URL path parameter.
// It attempts to delete the user by calling the UserService's DeleteUser method.
// On successful deletion, it responds with HTTP 204 No Content.
// If the user specified by the ID is not found, it responds with HTTP 404 Not Found.
// For other deletion failures, it responds with HTTP 500 Internal Server Error.
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	err := h.service.DeleteUser(userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("User with ID '%s' not found", userID)})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		}

		return
	}

	c.Status(http.StatusNoContent)
}
