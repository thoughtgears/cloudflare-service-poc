package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUsers handles HTTP GET requests to the /users endpoint.
// It retrieves all users by calling the UserService's GetUsers method.
// On success, it responds with HTTP 200 OK and a JSON array of user objects.
// On failure, it responds with HTTP 500 Internal Server Error and a JSON error message.
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.service.GetUsers()
	if err != nil {
		// Log the actual error internally if possible
		// log.Error().Err(err).Msg("Failed to retrieve users from service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, users)
}
