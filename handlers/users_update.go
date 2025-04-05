package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/thoughtgears/cloudflare-tunnels-poc/models"
	"github.com/thoughtgears/cloudflare-tunnels-poc/services"
)

// UpdateUserRequest defines the expected JSON payload structure for updating an existing user
// using the PUT method, where all updatable fields are expected to be provided.
// For partial updates (PATCH), pointer fields (*string, *bool etc.) and different
// validation tags (like 'omitempty') would typically be used.
type UpdateUserRequest struct {
	// FirstName is the user's given name. (Required)
	FirstName string `json:"first_name" validate:"required,min=1"`
	// LastName is the user's family name or surname. (Required)
	LastName string `json:"last_name" validate:"required,min=1"`
	// Email is the user's unique email address. (Required, must be valid email format)
	Email string `json:"email" validate:"required,email"`
	// Phone is the user's primary phone number. (Required)
	Phone string `json:"phone" validate:"required"`
	// Address is the user's physical address. (Required)
	Address string `json:"address" validate:"required"`
	// Active indicates whether the user's account should be active.
	Active bool `json:"active"`
	// Preferences contains the user's notification settings (Email/SMS).
	Preferences models.Preferences `json:"preferences"`
}

// UpdateUser handles HTTP PUT requests to the /users/:id endpoint.
// It extracts the user ID from the URL path parameter.
// It binds the incoming JSON request body to an UpdateUserRequest struct.
// It validates the bound request data using struct tags. If validation fails,
// it responds with HTTP 400 Bad Request and detailed validation errors.
// If binding or validation succeeds, it maps the request data to a models.User struct
// and calls the UserService's UpdateUser method with the ID and update data.
// On successful update, it responds with HTTP 200 OK and a JSON object representing
// the updated user.
// If the user specified by the ID is not found, it responds with HTTP 404 Not Found.
// For other update failures, it responds with HTTP 500 Internal Server Error.
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	var req UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})

		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": formatValidationErrors(err)})

		return
	}

	updatedData := models.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Phone:       req.Phone,
		Address:     req.Address,
		Active:      req.Active,
		Preferences: req.Preferences,
	}

	user, err := h.service.UpdateUser(userID, updatedData)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("User with ID '%s' not found", userID)})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		}

		return
	}
	c.JSON(http.StatusOK, user)
}
