package services

import (
	"context"
	"fmt"
	"github.com/upekZ/rest-api-go/internal/database/queries"
	"github.com/upekZ/rest-api-go/internal/types"
)

var uniqueFields = map[string]string{
	"Phone": "phone",
	"Email": "email",
}

type DB interface {
	GetUserByID(context.Context, string) (*types.UserEntity, error)
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context, string, *types.UserEntity) error
	GetUsers(context.Context) ([]queries.User, error)
	CreateUser(context.Context, *types.UserEntity) error
	IsEmailUnique(context.Context, string) (bool, error)
	IsPhoneUnique(context.Context, string) (bool, error)
}

type UserService struct {
	db    DB
	cache Cache
}

func NewUserService(db DB, cache Cache) *UserService {
	return &UserService{
		db:    db,
		cache: cache,
	}
}

func (o *UserService) CreateUser(ctx context.Context, user types.UserEntity) error {

	if state, err := types.ValidateUser(&user); state == false {
		return fmt.Errorf("user validation failure: %v", err)
	}

	if isUnique, err := o.IsUniqueField(ctx, uniqueFields["Phone"], user.Phone); !isUnique {
		o.cache.SetValue(uniqueFields["Phone"], user.Phone, true)
		return fmt.Errorf("user validation failure: %v", err)
	}

	if isUnique, err := o.IsUniqueField(ctx, uniqueFields["Email"], user.Email); !isUnique {
		o.cache.SetValue(uniqueFields["Email"], user.Email, true)
		return fmt.Errorf("user validation failure: %v", err)
	}
	if err := o.db.CreateUser(ctx, &user); err != nil {
		return fmt.Errorf("user creation failure in db: %v", err)
	}

	o.cache.SetValue(uniqueFields["Phone"], user.Phone, true)
	o.cache.SetValue(uniqueFields["Email"], user.Email, true)

	return nil
}

func (o *UserService) ListUsers(ctx context.Context) ([]queries.User, error) {

	users, err := o.db.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("user retrieval failure in db: %w", err)
	}
	return users, nil
}

func (o *UserService) GetUserByID(ctx context.Context, userID string) (*types.UserEntity, error) {
	user, err := o.db.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user retrieval failure in db: %w", err)
	}

	return user, nil
}

func (o *UserService) DeleteUser(ctx context.Context, userID string) error {

	user, err := o.db.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return fmt.Errorf("user not found: %w", err)
	}
	if err := o.db.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("user deletion failure in db: %w", err)
	}

	o.cache.DeleteField(uniqueFields["Phone"], user.Phone)
	o.cache.DeleteField(uniqueFields["Email"], user.Email)
	return nil
}
func (o *UserService) UpdateUser(ctx context.Context, userID string, userManager *types.UserEntity) error {
	if err := o.db.UpdateUser(ctx, userID, userManager); err != nil {
		return fmt.Errorf("user update failure in db: %w", err)
	}
	return nil
}
