package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"apiserver/internal/domain"
	"apiserver/internal/repositories/mocks" // Import the mock
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserInteractor_CreateNewUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	interactor := NewUserInteractor(mockRepo)

	name := "Test User"
	email := "test@example.com"
	plainPassword := "password123"

	// For CreateUser, the interactor generates the ID and calls repo, then repo returns the full user (potentially with DB-set fields like CreatedAt)
	// The mock should reflect what the repo's CreateUser is expected to return AFTER a successful creation.
	expectedUserFromRepo := &domain.User{ID: "new-uuid", Name: name, Email: email, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*domain.User"), mock.AnythingOfType("string")).Run(func(args mock.Arguments) {
		userArg := args.Get(1).(*domain.User) // This is the user *before* ID and timestamps are set by DB
		assert.Equal(t, name, userArg.Name)
		assert.Equal(t, email, userArg.Email)

		hashedPassArg := args.Get(2).(string)
		err := bcrypt.CompareHashAndPassword([]byte(hashedPassArg), []byte(plainPassword))
		assert.NoError(t, err, "Password should be hashed correctly")

	}).Return(expectedUserFromRepo, nil).Once() // This is what the repo returns

	createdUser, err := interactor.CreateNewUser(context.Background(), name, email, plainPassword)

	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, expectedUserFromRepo.ID, createdUser.ID)
	assert.Equal(t, expectedUserFromRepo.Name, createdUser.Name)
	assert.Equal(t, expectedUserFromRepo.Email, createdUser.Email)
	mockRepo.AssertExpectations(t)
}

func TestUserInteractor_CreateNewUser_Error_Hashing(t *testing.T) {
    // This test is a bit tricky as bcrypt.GenerateFromPassword rarely fails with valid inputs.
    // We're mostly testing our error propagation if it *were* to fail.
    // For this, we'd need to mock bcrypt or cause it to fail, which is hard.
    // A simpler test: check validation (already done) or repo error.
}

func TestUserInteractor_CreateNewUser_Error_Repo(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	interactor := NewUserInteractor(mockRepo)

	name := "Test User"
	email := "test@example.com"
	plainPassword := "password123"
	repoError := errors.New("repository error")

	mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*domain.User"), mock.AnythingOfType("string")).Return(nil, repoError).Once()

	_, err := interactor.CreateNewUser(context.Background(), name, email, plainPassword)

	assert.Error(t, err)
	assert.Equal(t, repoError, err)
	mockRepo.AssertExpectations(t)
}


func TestUserInteractor_CreateNewUser_Error_Validation(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository) 
	interactor := NewUserInteractor(mockRepo)

	_, err := interactor.CreateNewUser(context.Background(), "", "test@example.com", "password123")
	assert.Error(t, err)
	assert.Equal(t, "name, email, and password are required", err.Error())
}

func TestUserInteractor_FindUserByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	interactor := NewUserInteractor(mockRepo)

	userID := "test-id"
	expectedUser := &domain.User{ID: userID, Name: "Found User", Email: "found@example.com"}

	mockRepo.On("GetUserByID", mock.Anything, userID).Return(expectedUser, nil).Once()

	user, err := interactor.FindUserByID(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestUserInteractor_FindUserByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	interactor := NewUserInteractor(mockRepo)

	userID := "not-found-id"
	mockRepo.On("GetUserByID", mock.Anything, userID).Return(nil, nil).Once() 

	user, err := interactor.FindUserByID(context.Background(), userID)

	assert.NoError(t, err) 
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

func TestUserInteractor_FindUserByID_Error_Repo(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	interactor := NewUserInteractor(mockRepo)

	userID := "test-id"
	repoError := errors.New("repository error")
	mockRepo.On("GetUserByID", mock.Anything, userID).Return(nil, repoError).Once()

	_, err := interactor.FindUserByID(context.Background(), userID)

	assert.Error(t, err)
	assert.Equal(t, repoError, err)
	mockRepo.AssertExpectations(t)
}

func TestUserInteractor_FindUserByID_Error_Validation(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	interactor := NewUserInteractor(mockRepo)

	_, err := interactor.FindUserByID(context.Background(), "")
    assert.Error(t, err)
    assert.Equal(t, "user ID is required", err.Error())
}


// Tests for GetAllUsers
func TestUserInteractor_GetAllUsers_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	interactor := NewUserInteractor(mockRepo)

	expectedUsers := []domain.User{
		{ID: "id1", Name: "User One", Email: "one@example.com"},
		{ID: "id2", Name: "User Two", Email: "two@example.com"},
	}
	mockRepo.On("ListUsers", mock.Anything).Return(expectedUsers, nil).Once()

	users, err := interactor.GetAllUsers(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	mockRepo.AssertExpectations(t)
}

func TestUserInteractor_GetAllUsers_Error_Repo(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	interactor := NewUserInteractor(mockRepo)

	repoError := errors.New("repository error")
	mockRepo.On("ListUsers", mock.Anything).Return(nil, repoError).Once()

	_, err := interactor.GetAllUsers(context.Background())

	assert.Error(t, err)
	assert.Equal(t, repoError, err)
	mockRepo.AssertExpectations(t)
}

// Tests for UpdateExistingUser
func TestUserInteractor_UpdateExistingUser_Success_AllFields(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	interactor := NewUserInteractor(mockRepo)

	userID := "user-to-update"
	newName := "Updated Name"
	newEmail := "updated@example.com"
	newPlainPassword := "newPassword123"

	updatedDomainUser := &domain.User{Name: newName, Email: newEmail} // This is what's passed to repo (without ID)
	expectedUserFromRepo := &domain.User{ID: userID, Name: newName, Email: newEmail} // This is what repo returns

	mockRepo.On("UpdateUser", mock.Anything, userID, mock.MatchedBy(func(du *domain.User) bool {
		return du.Name == newName && du.Email == newEmail
	}), mock.AnythingOfType("*string")).Run(func(args mock.Arguments) {
		hashedPassArg := args.Get(3).(*string)
		assert.NotNil(t, hashedPassArg)
		err := bcrypt.CompareHashAndPassword([]byte(*hashedPassArg), []byte(newPlainPassword))
		assert.NoError(t, err, "Password should be hashed correctly for update")
	}).Return(expectedUserFromRepo, nil).Once()

	user, err := interactor.UpdateExistingUser(context.Background(), userID, &newName, &newEmail, &newPlainPassword)

	assert.NoError(t, err)
	assert.Equal(t, expectedUserFromRepo, user)
	mockRepo.AssertExpectations(t)
}

func TestUserInteractor_UpdateExistingUser_Success_PartialUpdate_NameOnly(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	interactor := NewUserInteractor(mockRepo)

	userID := "user-to-update"
	newName := "Just Name Updated"
    
	updatedDomainUser := &domain.User{Name: newName} // This is what's passed to repo (without ID)
	expectedUserFromRepo := &domain.User{ID: userID, Name: newName, Email: "original@example.com"} 

	mockRepo.On("UpdateUser", mock.Anything, userID, mock.MatchedBy(func(du *domain.User) bool {
		return du.Name == newName && du.Email == "" // Email in updateData will be empty
	}), (*string)(nil)).Return(expectedUserFromRepo, nil).Once() // No password update

	user, err := interactor.UpdateExistingUser(context.Background(), userID, &newName, nil, nil)

	assert.NoError(t, err)
	assert.Equal(t, expectedUserFromRepo, user)
	mockRepo.AssertExpectations(t)
}


func TestUserInteractor_UpdateExistingUser_Error_Validation_NoID(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	interactor := NewUserInteractor(mockRepo)
	someName := "name"
	_, err := interactor.UpdateExistingUser(context.Background(), "", &someName, nil, nil)
	assert.Error(t, err)
	assert.Equal(t, "user ID is required for update", err.Error())
}

func TestUserInteractor_UpdateExistingUser_Error_Validation_NoData(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	interactor := NewUserInteractor(mockRepo)
	_, err := interactor.UpdateExistingUser(context.Background(), "some-id", nil, nil, nil)
	assert.Error(t, err)
	assert.Equal(t, "no update data provided", err.Error())
}

func TestUserInteractor_UpdateExistingUser_Error_Repo(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	interactor := NewUserInteractor(mockRepo)
	
	userID := "user-to-update"
	name := "name"
	repoError := errors.New("repo update error")

	mockRepo.On("UpdateUser", mock.Anything, userID, mock.AnythingOfType("*domain.User"), (*string)(nil)).Return(nil, repoError).Once()

	_, err := interactor.UpdateExistingUser(context.Background(), userID, &name, nil, nil)
	assert.Error(t, err)
	assert.Equal(t, repoError, err)
	mockRepo.AssertExpectations(t)
}

func TestUserInteractor_UpdateExistingUser_Error_PasswordEmpty(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	interactor := NewUserInteractor(mockRepo)
	userID := "user-to-update"
	emptyPassword := ""
	_, err := interactor.UpdateExistingUser(context.Background(), userID, nil, nil, &emptyPassword)
	assert.Error(t, err)
	assert.Equal(t, "password cannot be updated to empty string", err.Error())
}


// Tests for RemoveUser
func TestUserInteractor_RemoveUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	interactor := NewUserInteractor(mockRepo)

	userID := "user-to-delete"
	mockRepo.On("DeleteUser", mock.Anything, userID).Return(nil).Once()

	err := interactor.RemoveUser(context.Background(), userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserInteractor_RemoveUser_Error_Validation(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	interactor := NewUserInteractor(mockRepo)

	err := interactor.RemoveUser(context.Background(), "")
	assert.Error(t, err)
	assert.Equal(t, "user ID is required", err.Error())
}

func TestUserInteractor_RemoveUser_Error_Repo(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	interactor := NewUserInteractor(mockRepo)

	userID := "user-to-delete"
	repoError := errors.New("repo delete error")
	mockRepo.On("DeleteUser", mock.Anything, userID).Return(repoError).Once()

	err := interactor.RemoveUser(context.Background(), userID)
	assert.Error(t, err)
	assert.Equal(t, repoError, err)
	mockRepo.AssertExpectations(t)
}
