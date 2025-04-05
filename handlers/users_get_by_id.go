package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/thoughtgears/cloudflare-tunnels-poc/services"
)

// GetUserByID handles HTTP GET requests to the /users/:id endpoint.
// It extracts the user ID from the URL path parameter.
// It retrieves the specific user by calling the UserService's GetUserByID method.
// On success, it responds with HTTP 200 OK and a JSON object representing the user.
// If the user is not found, it responds with HTTP 404 Not Found.
// For other errors, it responds with HTTP 500 Internal Server Error.
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("User with ID '%s' not found", userID)})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		}

		return
	}

	c.JSON(http.StatusOK, user)
}
