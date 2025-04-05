package handlers

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/thoughtgears/cloudflare-tunnels-poc/services"
)

// UserHandler encapsulates the dependencies, such as the UserService,
// required by the HTTP handler methods related to user operations.
// Methods associated with this struct handle incoming API requests for the /users endpoints.
type UserHandler struct {
	service services.UserService
}

// NewUserHandler is a constructor function that creates and returns a new instance
// of UserHandler. It requires a UserService dependency to be injected, which will
// be used by the handler methods to interact with the underlying user data store
// and business logic.
//
// Parameters:
//   - service: An instance implementing the services.UserService interface.
//
// Returns:
//   - A pointer to a newly created UserHandler instance.
func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// validate holds a package-level instance of the validator engine.
// It is initialized once and reused by handler functions within this package
// to validate incoming request data transfer objects (DTOs) or request structs
// based on the 'validate' struct tags.
var validate = validator.New()

// formatValidationErrors is a helper function that converts validation errors
// returned by the 'go-playground/validator' library into a user-friendly map format,
// suitable for returning as a JSON error response body in the API.
//
// If the input error is of type validator.ValidationErrors, it iterates through
// each field error, creating a descriptive message based on the validation tag
// (e.g., "required", "email", "min"). The map keys are the field names and
// the values are the corresponding error messages.
//
// If the input error is not a validator.ValidationErrors, it returns a generic
// error message under the "error" key.
//
// Parameters:
//   - err: The error returned by the call to validate.Struct().
//
// Returns:
//   - A map[string]string where keys are field names (or "error") and values
//     are user-readable validation error messages.
func formatValidationErrors(err error) map[string]string {
	errorsMap := make(map[string]string)
	var validationErrs validator.ValidationErrors

	// Use errors.As for type assertion, which is generally preferred over direct type assertion.
	if errors.As(err, &validationErrs) {
		for _, fieldErr := range validationErrs {
			// Use fieldErr.Field() for the field name and create a user-friendly message
			// based on fieldErr.Tag()
			fieldName := fieldErr.Field()
			switch fieldErr.Tag() {
			case "required":
				errorsMap[fieldName] = fmt.Sprintf("%s is required", fieldName)
			case "email":
				errorsMap[fieldName] = fmt.Sprintf("%s must be a valid email address", fieldName)
			case "min":
				errorsMap[fieldName] = fmt.Sprintf("%s must be at least %s characters long", fieldName, fieldErr.Param())
			// Add more cases for other validation tags you use
			default:
				errorsMap[fieldName] = fmt.Sprintf("Invalid value for %s (%s)", fieldName, fieldErr.Tag())
			}
		}
	} else {
		// Handle cases where the error is not from the validator (less common after ShouldBindJSON)
		errorsMap["error"] = "Invalid request data structure" // More specific than generic "Invalid request data"
		// Optionally log the actual error `err` here for debugging
		// log.Printf("Non-validation error during request processing: %v", err)
	}

	return errorsMap
}
