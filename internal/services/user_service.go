package services

import (
	"context"
	"fmt"
	"github.com/upekZ/rest-api-go/internal/database/queries"
	"github.com/upekZ/rest-api-go/internal/types"
)

type DB interface {
	GetUserByID(context.Context, string) (*types.UserManager, error)
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context, string, *types.UserManager) error
	GetUsers(context.Context) ([]queries.User, error)
	CreateUser(context.Context, *types.UserManager) error
}

type UserService struct {
	db DB
}

func NewUserService(db DB) *UserService {
	return &UserService{db: db}
}

func (o *UserService) CreateUser(ctx context.Context, user types.UserManager) error {

	if state := types.ValidateUser(&user); state == false {
		return fmt.Errorf("user validation failure")
	}
	if err := o.db.CreateUser(ctx, &user); err != nil {
		return fmt.Errorf("user creation failure in db: %w", err)
	}
	return nil
}

func (o *UserService) ListUsers(ctx context.Context) ([]queries.User, error) {

	users, err := o.db.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("user retrieval failure in db: %w", err)
	}
	return users, nil
}

func (o *UserService) GetUserByID(ctx context.Context, userID string) (*types.UserManager, error) {
	users, err := o.db.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user retrieval failure in db: %w", err)
	}

	return users, nil
}

func (o *UserService) DeleteUser(ctx context.Context, userID string) error {
	if err := o.db.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("user deletion failure in db: %w", err)
	}
	return nil
}
func (o *UserService) UpdateUser(ctx context.Context, userID string, userManager *types.UserManager) error {
	if err := o.db.UpdateUser(ctx, userID, userManager); err != nil {
		return fmt.Errorf("user update failure in db: %w", err)
	}
	return nil
}
