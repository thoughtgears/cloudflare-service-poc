package models

import (
	"time"
)

// Preferences holds settings related to user notifications.
// It specifies whether a user wishes to receive notifications via Email and/or SMS.
type Preferences struct {
	// Email indicates if the user wants email notifications (true = yes, false = no).
	Email bool `json:"email" faker:"-"`
	// SMS indicates if the user wants SMS text message notifications (true = yes, false = no).
	SMS bool `json:"sms" faker:"-"`
}

// User defines the structure for storing user information within the system.
// It includes personal details, contact information, address, account status,
// notification preferences, and timestamps for record management.
type User struct {
	// ID is the unique identifier for the user, typically a UUID.
	ID string `json:"id" faker:"uuid_hyphenated"` // Note: faker tag generates UUID for fake data
	// FirstName is the user's given name.
	FirstName string `json:"first_name" faker:"first_name"`
	// LastName is the user's family name or surname.
	LastName string `json:"last_name" faker:"last_name"`
	// Email is the user's unique email address, used for login and communication.
	Email string `json:"email" faker:"email"`
	// Phone is the user's primary phone number.
	Phone string `json:"phone" faker:"phone_number"`
	// Address is the user's physical address (currently stored as a single string).
	// Consider using a structured Address type for more detail if needed in the future.
	Address string `json:"address" faker:"real_address"`
	// Active indicates whether the user's account is currently active (true) or inactive (false).
	Active bool `json:"active" faker:"-"`
	// Preferences embeds the notification settings for the user.
	Preferences Preferences `json:"preferences" faker:"-"`
	// CreatedAt records the exact date and time when the user record was created in the system.
	CreatedAt time.Time `json:"created_at" faker:"-"`
	// UpdatedAt records the exact date and time when the user record was last modified.
	UpdatedAt time.Time `json:"updated_at" faker:"-"`
}
