package services

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"

	"github.com/thoughtgears/cloudflare-tunnels-poc/models"
)

var ErrUserNotFound = errors.New("user not found")

// userStore holds the in-memory list of users.
var userStore []models.User

// mutex protects concurrent access to userStore.
var storeMutex sync.RWMutex

// UserService defines the contract for user operations.
type UserService interface {
	// GetUsers returns all users currently in the store.
	GetUsers() ([]models.User, error)
	// GetUserByID returns a single user matching the provided ID.
	// Returns ErrUserNotFound if the user does not exist.
	GetUserByID(id string) (*models.User, error)
	// CreateUser adds a new user to the store.
	// It assigns a new ID and sets CreatedAt/UpdatedAt timestamps.
	CreateUser(user models.User) (*models.User, error)
	// UpdateUser updates an existing user identified by ID.
	// Only updates specified fields (excluding ID, CreatedAt). Updates UpdatedAt.
	// Returns ErrUserNotFound if the user does not exist.
	UpdateUser(id string, updatedData models.User) (*models.User, error)
	// DeleteUser removes a user identified by ID from the store.
	// Returns ErrUserNotFound if the user does not exist.
	DeleteUser(id string) error
}

// userServiceImpl provides a concrete implementation of UserService.
type userServiceImpl struct {
	// This implementation directly uses the package-level userStore and mutex.
	// No fields needed here in this specific setup.
	// If there was need for a repository this is where it should be initialized.
}

// NewUserService creates a new instance of the user service.
// In this specific setup with package-level storage, this function
// returns a pointer to a stateless struct that operates on the shared store.
func NewUserService() UserService {
	return &userServiceImpl{}
}

// init initializes the user service with a random number of fake users.
func init() {
	// Determine a random number of users between 5 and 20
	// rand.Intn(16) gives 0-15, adding 5 gives 5-20.
	userCount := rand.Intn(16) + 5 //nolint:gosec // G404: Non-sensitive count for fake data generation

	fmt.Printf("Initializing user service with %d users...\n", userCount)

	// Lock the store for initial population
	storeMutex.Lock()
	defer storeMutex.Unlock() // Ensure unlock even if panic occurs

	// Pre-allocate the slice with the correct size
	userStore = make([]models.User, userCount)
	now := time.Now() // Get current time once for CreatedAt/UpdatedAt

	for i := 0; i < userCount; i++ {
		// Or create a temporary local variable (as you did)
		tempUser := models.User{}

		// Use faker.FakeData based on the struct tags in models.User
		// This should populate fields *without* faker:"-"
		err := faker.FakeData(&tempUser)
		if err != nil {
			fmt.Printf("Error generating fake data for user %d: %v\n", i, err)
			// Continue with potentially partially filled data or handle differently
		}

		// --- Manually set fields AFTER faker ---
		tempUser.ID = uuid.NewString() // Assign unique ID
		tempUser.CreatedAt = now       // Set creation timestamp
		tempUser.UpdatedAt = now       // Set updated timestamp

		// Manually set booleans (overriding any potential faker value if needed, or just setting)
		tempUser.Active = rand.Intn(2) == 1            //nolint:gosec // G404: Non-sensitive fake active status
		tempUser.Preferences.Email = rand.Intn(2) == 1 //nolint:gosec // G404: Non-sensitive fake preference
		tempUser.Preferences.SMS = rand.Intn(2) == 1   //nolint:gosec // G404: Non-sensitive fake preference

		// --- Assign the fully populated tempUser to the slice index ---
		userStore[i] = tempUser // Assign the complete struct to the slice index i
	}
	fmt.Println("User service initialized.")
}

// GetUsers retrieves all users currently stored in memory.
//
// It returns a slice containing copies of the user data. Modifications to the
// returned slice or its elements will not affect the internal user store.
// This function is safe for concurrent use.
//
// Currently, it always returns a nil error, but the signature allows for
// future error handling.
func (s *userServiceImpl) GetUsers() ([]models.User, error) {
	storeMutex.RLock() // Lock for reading
	defer storeMutex.RUnlock()
	usersCopy := make([]models.User, len(userStore))
	copy(usersCopy, userStore)

	return usersCopy, nil
}

// GetUserByID searches for and returns a single user based on their unique ID.
//
// Parameters:
//   - id: The UUID string of the user to retrieve.
//
// Returns:
//   - A pointer to a copy of the found models.User struct if a user with the
//     specified ID exists. Modifications to the returned user will not affect
//     the internal user store.
//   - nil and ErrUserNotFound if no user matches the provided ID.
//   - nil and potentially other errors in the future (currently only returns ErrUserNotFound on failure).
//
// This function is safe for concurrent use.
func (s *userServiceImpl) GetUserByID(id string) (*models.User, error) {
	storeMutex.RLock() // Lock for reading
	defer storeMutex.RUnlock()
	for i := range userStore {
		if userStore[i].ID == id {
			userCopy := userStore[i]

			return &userCopy, nil
		}
	}

	return nil, ErrUserNotFound
}

// CreateUser adds a new user to the in-memory store.
//
// It takes a models.User struct as input. The ID, CreatedAt, and UpdatedAt
// fields of the input struct are ignored; new values will be generated and assigned.
// A new UUID is generated for the ID. CreatedAt and UpdatedAt are set to the
// current time (formatted as a string based on time.Now().String()).
// The new user is appended to the internal user store.
//
// Parameters:
//   - user: A models.User struct containing the desired data for the new user.
//     ID, CreatedAt, and UpdatedAt fields will be overwritten.
//
// Returns:
//   - A pointer to a copy of the newly created user struct, including the
//     assigned ID and timestamps.
//   - A nil error (currently no specific error conditions are handled during creation).
//
// This function is safe for concurrent use.
func (s *userServiceImpl) CreateUser(user models.User) (*models.User, error) {
	storeMutex.Lock() // Lock for writing
	defer storeMutex.Unlock()
	now := time.Now()
	user.ID = uuid.NewString()
	user.CreatedAt = now
	user.UpdatedAt = now
	userStore = append(userStore, user)
	createdUserCopy := user

	return &createdUserCopy, nil
}

// UpdateUser finds a user by ID and replaces their stored data (except for ID
// and CreatedAt) with the provided data.
//
// It searches for the user matching the given ID. If found, it updates the
// user's fields (FirstName, LastName, Email, Phone, Address, Active, Preferences)
// in the internal store with the values from the updatedData parameter.
// The user's UpdatedAt field is set to the current time (formatted as a string).
// The user's ID and CreatedAt fields remain unchanged.
//
// Parameters:
//   - id: The UUID string of the user to update.
//   - updatedData: A models.User struct containing the new data for the user.
//     ID and CreatedAt fields from this parameter are ignored.
//
// Returns:
//   - A pointer to a copy of the updated models.User struct as it exists in the store
//     after the update, including the new UpdatedAt timestamp.
//   - nil and ErrUserNotFound if no user matches the provided ID.
//   - nil and potentially other errors in the future.
//
// This function is safe for concurrent use.
func (s *userServiceImpl) UpdateUser(id string, updatedData models.User) (*models.User, error) {
	storeMutex.Lock() // Lock for writing
	defer storeMutex.Unlock()
	foundIndex := -1
	for i := range userStore {
		if userStore[i].ID == id {
			foundIndex = i

			break
		}
	}
	if foundIndex == -1 {
		return nil, ErrUserNotFound
	}
	originalUser := &userStore[foundIndex]
	originalUser.FirstName = updatedData.FirstName
	originalUser.LastName = updatedData.LastName
	originalUser.Email = updatedData.Email
	originalUser.Phone = updatedData.Phone
	originalUser.Address = updatedData.Address
	originalUser.Active = updatedData.Active
	originalUser.Preferences = updatedData.Preferences
	originalUser.UpdatedAt = time.Now()
	updatedUserCopy := *originalUser

	return &updatedUserCopy, nil
}

// DeleteUser removes a user from the in-memory store based on their unique ID.
//
// Parameters:
//   - id: The UUID string of the user to delete.
//
// Returns:
//   - nil if the user was successfully found and removed.
//   - ErrUserNotFound if no user matches the provided ID.
//
// This function modifies the internal user store and is safe for concurrent use.
func (s *userServiceImpl) DeleteUser(id string) error {
	storeMutex.Lock() // Lock for writing
	defer storeMutex.Unlock()
	foundIndex := -1
	for i := range userStore {
		if userStore[i].ID == id {
			foundIndex = i

			break
		}
	}
	if foundIndex == -1 {
		return ErrUserNotFound
	}
	userStore = append(userStore[:foundIndex], userStore[foundIndex+1:]...)

	return nil
}
