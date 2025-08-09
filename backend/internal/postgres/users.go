package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/flexGURU/flower-haven/backend/internal/postgres/generated"
	"github.com/flexGURU/flower-haven/backend/internal/repository"
	"github.com/flexGURU/flower-haven/backend/pkg"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ repository.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	queries *generated.Queries
}

func NewUserRepository(queries *generated.Queries) *UserRepository {
	return &UserRepository{queries: queries}
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *repository.User) (*repository.User, error) {
	params := generated.CreateUserParams{
		Name:         user.Name,
		Email:        user.Email,
		PhoneNumber:  user.PhoneNumber,
		Password:     *user.Password,
		IsAdmin:      user.IsAdmin,
		RefreshToken: pgtype.Text{Valid: false},
		Address:      pgtype.Text{Valid: false},
	}

	if user.Address != nil {
		params.Address = pgtype.Text{
			Valid:  true,
			String: *user.Address,
		}
	}
	if user.RefreshToken != nil {
		params.RefreshToken = pgtype.Text{
			Valid:  true,
			String: *user.RefreshToken,
		}
	}

	generatedUser, err := ur.queries.CreateUser(ctx, params)
	if err != nil {
		if pkg.PgxErrorCode(err) == pkg.UNIQUE_VIOLATION {
			return nil, pkg.Errorf(pkg.ALREADY_EXISTS_ERROR, "%s", err.Error())
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error creating user: %s", err.Error())
	}

	user.ID = uint32(generatedUser.ID)
	user.CreatedAt = generatedUser.CreatedAt
	user.Password = nil
	user.RefreshToken = nil
	user.IsActive = true

	return user, nil
}

func (ur *UserRepository) GetUserByID(ctx context.Context, id int64) (*repository.User, error) {
	generatedUser, err := ur.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "user with ID %d not found", id)
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error fetching user by ID: %s", err.Error())
	}

	user := &repository.User{
		ID:          uint32(generatedUser.ID),
		Name:        generatedUser.Name,
		Email:       generatedUser.Email,
		Address:     nil,
		PhoneNumber: generatedUser.PhoneNumber,
		IsAdmin:     generatedUser.IsAdmin,
		IsActive:    generatedUser.IsActive,
		CreatedAt:   generatedUser.CreatedAt,
	}

	if generatedUser.Address.Valid {
		address := generatedUser.Address.String
		user.Address = &address
	}

	return user, nil
}

func (ur *UserRepository) GetUserInternal(ctx context.Context, id int64, email string) (*repository.User, error) {
	if id == 0 && email == "" {
		return nil, pkg.Errorf(pkg.INVALID_ERROR, "either id or email must be provided")
	}

	var err error
	var generatedUser generated.User

	if id != 0 {
		generatedUser, err = ur.queries.GetUserByID(ctx, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "user with ID %d not found", id)
			}
			return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error fetching user by ID: %s", err.Error())
		}
	} else if email != "" {
		generatedUser, err = ur.queries.GetUserByEmail(ctx, email)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "user with email %s not found", email)
			}
			if errors.Is(err, sql.ErrNoRows) {
				return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "user with ID %d not found", id)
			}
			return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error fetching user by email: %s", err.Error())
		}
	}

	user := &repository.User{
		ID:           uint32(generatedUser.ID),
		Name:         generatedUser.Name,
		Email:        generatedUser.Email,
		Address:      nil,
		PhoneNumber:  generatedUser.PhoneNumber,
		IsAdmin:      generatedUser.IsAdmin,
		IsActive:     generatedUser.IsActive,
		RefreshToken: &generatedUser.RefreshToken.String,
		Password:     &generatedUser.Password,
		CreatedAt:    generatedUser.CreatedAt,
	}

	if generatedUser.Address.Valid {
		address := generatedUser.Address.String
		user.Address = &address
	}

	return user, nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*repository.User, error) {
	generatedUser, err := ur.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "user with email %s not found", email)
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error fetching user by email: %s", err.Error())
	}

	user := &repository.User{
		ID:          uint32(generatedUser.ID),
		Name:        generatedUser.Name,
		Email:       generatedUser.Email,
		Address:     nil,
		PhoneNumber: generatedUser.PhoneNumber,
		IsAdmin:     generatedUser.IsAdmin,
		IsActive:    generatedUser.IsActive,
		CreatedAt:   generatedUser.CreatedAt,
	}

	if generatedUser.Address.Valid {
		address := generatedUser.Address.String
		user.Address = &address
	}

	return user, nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, user *repository.UpdateUser) (*repository.User, error) {
	params := generated.UpdateUserParams{
		ID:           int64(user.ID),
		Name:         pgtype.Text{Valid: false},
		Address:      pgtype.Text{Valid: false},
		PhoneNumber:  pgtype.Text{Valid: false},
		Password:     pgtype.Text{Valid: false},
		RefreshToken: pgtype.Text{Valid: false},
		IsAdmin:      pgtype.Bool{Valid: false},
		IsActive:     pgtype.Bool{Valid: false},
	}

	if user.Name != nil {
		params.Name = pgtype.Text{
			Valid:  true,
			String: *user.Name,
		}
	}
	if user.Address != nil {
		params.Address = pgtype.Text{
			Valid:  true,
			String: *user.Address,
		}
	}
	if user.PhoneNumber != nil {
		params.PhoneNumber = pgtype.Text{
			Valid:  true,
			String: *user.PhoneNumber,
		}
	}
	if user.Password != nil {
		params.Password = pgtype.Text{
			Valid:  true,
			String: *user.Password,
		}
	}
	if user.IsAdmin != nil {
		params.IsAdmin = pgtype.Bool{
			Valid: true,
			Bool:  *user.IsAdmin,
		}
	}
	if user.IsActive != nil {
		params.IsActive = pgtype.Bool{
			Valid: true,
			Bool:  *user.IsActive,
		}
	}
	if user.RefreshToken != nil {
		params.RefreshToken = pgtype.Text{
			Valid:  true,
			String: *user.RefreshToken,
		}
	}

	generatedUser, err := ur.queries.UpdateUser(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "user with ID %d not found", user.ID)
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error updating user: %s", err.Error())
	}

	updatedUser := &repository.User{
		ID:          uint32(generatedUser.ID),
		Name:        generatedUser.Name,
		Email:       generatedUser.Email,
		Address:     nil,
		PhoneNumber: generatedUser.PhoneNumber,
		IsAdmin:     generatedUser.IsAdmin,
		IsActive:    generatedUser.IsActive,
		CreatedAt:   generatedUser.CreatedAt,
	}

	if generatedUser.Address.Valid {
		address := generatedUser.Address.String
		updatedUser.Address = &address
	}

	return updatedUser, nil
}

func (ur *UserRepository) ListUsers(ctx context.Context, filter *repository.UserFilter) ([]*repository.User, *pkg.Pagination, error) {
	paramListUsers := generated.ListUsersParams{
		Limit:    int32(filter.Pagination.PageSize),
		Offset:   pkg.Offset(filter.Pagination.Page, filter.Pagination.PageSize),
		Search:   pgtype.Text{Valid: false},
		IsActive: pgtype.Bool{Valid: false},
		IsAdmin:  pgtype.Bool{Valid: false},
	}

	paramCountUsers := generated.ListUsersCountParams{
		Search:   pgtype.Text{Valid: false},
		IsActive: pgtype.Bool{Valid: false},
		IsAdmin:  pgtype.Bool{Valid: false},
	}

	if filter.Search != nil {
		search := strings.ToLower(*filter.Search)
		paramListUsers.Search = pgtype.Text{
			Valid:  true,
			String: "%" + search + "%",
		}
		paramCountUsers.Search = pgtype.Text{
			Valid:  true,
			String: "%" + search + "%",
		}
	}

	if filter.IsAdmin != nil {
		paramListUsers.IsAdmin = pgtype.Bool{
			Valid: true,
			Bool:  *filter.IsAdmin,
		}
		paramCountUsers.IsAdmin = pgtype.Bool{
			Valid: true,
			Bool:  *filter.IsAdmin,
		}
	}

	if filter.IsActive != nil {
		paramListUsers.IsActive = pgtype.Bool{
			Valid: true,
			Bool:  *filter.IsActive,
		}
		paramCountUsers.IsActive = pgtype.Bool{
			Valid: true,
			Bool:  *filter.IsActive,
		}
	}

	generatedUsers, err := ur.queries.ListUsers(ctx, paramListUsers)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error listing users: %s", err.Error())
	}

	totalCount, err := ur.queries.ListUsersCount(ctx, paramCountUsers)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error counting users: %s", err.Error())
	}

	userList := make([]*repository.User, len(generatedUsers))
	for idx, generatedUser := range generatedUsers {
		user := &repository.User{
			ID:          uint32(generatedUser.ID),
			Name:        generatedUser.Name,
			Email:       generatedUser.Email,
			Address:     nil,
			PhoneNumber: generatedUser.PhoneNumber,
			IsAdmin:     generatedUser.IsAdmin,
			IsActive:    generatedUser.IsActive,
			CreatedAt:   generatedUser.CreatedAt,
		}

		if generatedUser.Address.Valid {
			address := generatedUser.Address.String
			user.Address = &address
		}

		userList[idx] = user
	}

	return userList, pkg.CalculatePagination(uint32(totalCount), filter.Pagination.PageSize, filter.Pagination.Page), nil
}
