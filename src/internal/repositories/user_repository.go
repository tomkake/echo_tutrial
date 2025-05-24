package repositories

import (
	"context"
	"database/sql" // For sql.Result, and potentially for db connection if not abstracted by sqlc Querier fully

	"apiserver/internal/domain"       // Our domain model
	db "apiserver/internal/db/sqlc" // sqlc generated package, aliased to db
	"github.com/google/uuid"
)

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User, hashedPassword string) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	ListUsers(ctx context.Context) ([]domain.User, error)
	UpdateUser(ctx context.Context, id string, user *domain.User, hashedPassword *string) (*domain.User, error) // hashedPassword is a pointer to allow optional update
	DeleteUser(ctx context.Context, id string) error
}

// sqlcUserRepository implements UserRepository using sqlc generated code.
type sqlcUserRepository struct {
	querier db.Querier // sqlc generated Querier interface
	dbConn  *sql.DB    // The underlying DB connection, needed for Querier usually.
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(conn *sql.DB) UserRepository {
	return &sqlcUserRepository{
		querier: db.New(conn), // db.New(conn) is the typical constructor for sqlc's Queries struct which implements Querier
		dbConn:  conn,
	}
}

// Helper to convert sqlc.User to domain.User
func toDomainUser(sqlcUser db.User) *domain.User {
	domainUser := &domain.User{
		ID:        sqlcUser.ID.String(),
		Password:  "", // Password is not exposed from DB to domain generally
		CreatedAt: sqlcUser.CreatedAt,
		UpdatedAt: sqlcUser.Updatedat, // Note: sqlc generated 'Updatedat'
	}
	if sqlcUser.Name.Valid {
		domainUser.Name = sqlcUser.Name.String
	}
	if sqlcUser.Email.Valid {
		domainUser.Email = sqlcUser.Email.String
	}
	// sqlcUser.Password.String could be assigned if needed, but typically not to domain model
	return domainUser
}

// Helper to convert a slice of sqlc.User to a slice of domain.User
func toDomainUserSlice(sqlcUsers []db.User) []domain.User {
    domainUsers := make([]domain.User, len(sqlcUsers))
    for i, su := range sqlcUsers {
        domainUsers[i] = *toDomainUser(su)
    }
    return domainUsers
}

func (r *sqlcUserRepository) CreateUser(ctx context.Context, user *domain.User, hashedPassword string) (*domain.User, error) {
	userID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	params := db.CreateUserParams{
		ID:       userID,
		Name:     sql.NullString{String: user.Name, Valid: user.Name != ""},
		Email:    sql.NullString{String: user.Email, Valid: user.Email != ""},
		Password: sql.NullString{String: hashedPassword, Valid: hashedPassword != ""},
	}

	_, err = r.querier.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	// Return the user by fetching it, so CreatedAt/UpdatedAt are populated
	return r.GetUserByID(ctx, userID.String())
}

func (r *sqlcUserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err // Invalid UUID format
	}

	sqlcUser, err := r.querier.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or a custom domain.ErrNotFound
		}
		return nil, err
	}
	return toDomainUser(sqlcUser), nil
}

func (r *sqlcUserRepository) ListUsers(ctx context.Context) ([]domain.User, error) {
	sqlcUsers, err := r.querier.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	return toDomainUserSlice(sqlcUsers), nil
}

func (r *sqlcUserRepository) UpdateUser(ctx context.Context, id string, user *domain.User, hashedPassword *string) (*domain.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err // Invalid UUID format
	}

	currentUser, err := r.querier.GetUserByID(ctx, userID)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // Or a custom domain.ErrNotFound indicating user to update not found
        }
        return nil, err // Other DB error
    }

    // Prepare params with current values, then overwrite with new ones if provided
	params := db.UpdateUserParams{
		ID:       userID,
		Name:     currentUser.Name,
		Email:    currentUser.Email,
		Password: currentUser.Password, // This is the current hashed password from DB
	}

	if user.Name != "" { // Assuming empty string means no update for name
		params.Name = sql.NullString{String: user.Name, Valid: true}
	}
	if user.Email != "" { // Assuming empty string means no update for email
		params.Email = sql.NullString{String: user.Email, Valid: true}
	}
	
	if hashedPassword != nil {
		params.Password = sql.NullString{String: *hashedPassword, Valid: *hashedPassword != ""}
	}

	_, err = r.querier.UpdateUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return r.GetUserByID(ctx, id) // Return the updated user
}

func (r *sqlcUserRepository) DeleteUser(ctx context.Context, id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return err // Invalid UUID format
	}
	_, err = r.querier.DeleteUser(ctx, userID)
	return err
}
