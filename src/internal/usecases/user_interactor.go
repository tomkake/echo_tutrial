package usecases

import (
	"context"
	"errors" // For standard errors

	"apiserver/internal/domain"
	"apiserver/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

// UserInteractor defines the interface for user-related business logic.
type UserInteractor interface {
	CreateNewUser(ctx context.Context, name, email, plainPassword string) (*domain.User, error)
	FindUserByID(ctx context.Context, id string) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	UpdateExistingUser(ctx context.Context, id string, name, email *string, plainPassword *string) (*domain.User, error)
	RemoveUser(ctx context.Context, id string) error
}

// userInteractor implements UserInteractor.
type userInteractor struct {
	userRepo repositories.UserRepository
}

// NewUserInteractor creates a new instance of UserInteractor.
func NewUserInteractor(repo repositories.UserRepository) UserInteractor {
	return &userInteractor{userRepo: repo}
}

func (uc *userInteractor) CreateNewUser(ctx context.Context, name, email, plainPassword string) (*domain.User, error) {
	if name == "" || email == "" || plainPassword == "" {
		return nil, errors.New("name, email, and password are required") // Basic validation
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashedPassword := string(hashedPasswordBytes)

	user := &domain.User{
		Name:  name,
		Email: email,
		// ID, CreatedAt, UpdatedAt will be handled by repository/DB
	}

	return uc.userRepo.CreateUser(ctx, user, hashedPassword)
}

func (uc *userInteractor) FindUserByID(ctx context.Context, id string) (*domain.User, error) {
	if id == "" {
		return nil, errors.New("user ID is required")
	}
	return uc.userRepo.GetUserByID(ctx, id)
}

func (uc *userInteractor) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	return uc.userRepo.ListUsers(ctx)
}

func (uc *userInteractor) UpdateExistingUser(ctx context.Context, id string, name, email *string, plainPassword *string) (*domain.User, error) {
	if id == "" {
		return nil, errors.New("user ID is required for update")
	}

	// Construct a domain.User for update, only setting fields if provided
	// The repository's UpdateUser method is responsible for merging with existing data or handling partials.
	updateData := &domain.User{} // Only pass non-nil fields to repo, or let repo handle merge
	hasUpdate := false
	if name != nil {
		updateData.Name = *name
		hasUpdate = true
	}
	if email != nil {
		updateData.Email = *email
		hasUpdate = true
	}

	var newHashedPassword *string
	if plainPassword != nil {
		if *plainPassword == "" {
		    return nil, errors.New("password cannot be updated to empty string")
                }
		hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(*plainPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		h := string(hashedPasswordBytes)
		newHashedPassword = &h
		hasUpdate = true
	}

	if !hasUpdate {
		return nil, errors.New("no update data provided") // Or fetch and return existing user
	}

	return uc.userRepo.UpdateUser(ctx, id, updateData, newHashedPassword)
}

func (uc *userInteractor) RemoveUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("user ID is required")
	}
	return uc.userRepo.DeleteUser(ctx, id)
}
