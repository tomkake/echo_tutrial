package domain // Changed from models

import "time"

// User represents the core domain entity for a user.
type User struct {
    ID        string    // Assuming UUID stored as string
    Name      string
    Email     string
    Password  string    // This is part of the domain, but might not be exposed directly
    CreatedAt time.Time
    UpdatedAt time.Time // Note: Schema had 'UpdatedAt'
}
