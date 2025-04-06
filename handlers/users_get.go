package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUsers handles HTTP GET requests to the /users endpoint.
// It retrieves all users by calling the UserService's GetUsers method.
// On success, it responds with HTTP 200 OK and a JSON array of user objects.
// On failure, it responds with HTTP 500 Internal Server Error and a JSON error message.
// @Summary		List all users
// @Description	get all users currently stored
// @Tags			users
// @Accept			json
// @Produce		json
// @Success		200	{array}		models.User			"Successfully retrieved list of users"
// @Failure		500	{object}	map[string]string	"Internal Server Error"
// @Router			/users [get]
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
