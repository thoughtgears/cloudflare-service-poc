package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/thoughtgears/cloudflare-tunnels-poc/models"
)

// CreateUserRequest defines the expected JSON payload structure for creating a new user.
// It includes the fields a client should provide. ID, CreatedAt, and UpdatedAt are
// excluded as they are managed by the system/service layer. Active status defaults
// to true within the handler logic.
type CreateUserRequest struct {
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
	// Preferences contains the user's notification settings (Email/SMS).
	Preferences models.Preferences `json:"preferences"`
}

// CreateUser handles HTTP POST requests to the /users endpoint.
// It binds the incoming JSON request body to a CreateUserRequest struct.
// It validates the bound request data using struct tags. If validation fails,
// it responds with HTTP 400 Bad Request and detailed validation errors.
// If binding or validation succeeds, it maps the request data to a models.User struct
// (setting Active to true by default) and calls the UserService's CreateUser method.
// On successful creation, it responds with HTTP 201 Created and a JSON object
// representing the newly created user (including system-generated fields like ID).
// On failure during user creation, it responds with HTTP 500 Internal Server Error.
// @Summary		Create a new user
// @Description	add a new user to the store based on JSON payload
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			user	body		handlers.CreateUserRequest	true	"User data to create"
// @Success		201		{object}	models.User					"Successfully created user"
// @Failure		400		{object}	map[string]any				"Validation Error or Invalid Request Format"
// @Failure		500		{object}	map[string]string			"Internal Server Error"
// @Router			/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})

		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": formatValidationErrors(err)})

		return
	}

	newUser := models.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Phone:       req.Phone,
		Address:     req.Address, // Assuming string address model
		Active:      true,
		Preferences: req.Preferences,
	}

	createdUser, err := h.service.CreateUser(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})

		return
	}

	c.JSON(http.StatusCreated, createdUser)
}
