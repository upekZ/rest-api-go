package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/upekZ/rest-api-go/internal/database/queries" //To be removed after moving usage of queries.User --> model.UserEntity
	"github.com/upekZ/rest-api-go/internal/model"
	"net/http"
)

var uniqueFields = map[string]string{
	"Phone": "phone",
	"Email": "email",
}

type DB interface {
	GetUserByID(context.Context, string) (*model.UserEntity, error)
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context, string, *model.UserEntity) error
	GetUsers(context.Context) ([]queries.User, error) //queries.User to be replaced with model.UserEntity
	CreateUser(context.Context, *model.UserEntity) error
	IsEmailUnique(context.Context, string) (bool, error)
	IsPhoneUnique(context.Context, string) (bool, error)
}

type WebSocketHandler interface {
	HandleWebSocket(w http.ResponseWriter, r *http.Request) error
	Broadcast(message []byte)
}

type UserService struct {
	db        DB
	cache     Cache
	wsHandler WebSocketHandler
}

func NewUserService(db DB, cache Cache, wsHandler WebSocketHandler) *UserService {
	return &UserService{
		db:        db,
		cache:     cache,
		wsHandler: wsHandler,
	}
}

func (o *UserService) CreateUser(ctx context.Context, user *model.UserEntity) error {

	if state, err := model.ValidateUser(user); state == false {
		return fmt.Errorf("user validation failure: %v", err)
	}

	//To Do: Iterate through Unique fields (ie: Phone and Email) to validate uniqueness
	if isUnique, err := o.IsUniqueField(ctx, uniqueFields["Phone"], user.Phone); !isUnique {
		o.cache.SetValue(uniqueFields["Phone"], user.Phone, true)
		return fmt.Errorf("user validation failure: %v", err)
	}

	if isUnique, err := o.IsUniqueField(ctx, uniqueFields["Email"], user.Email); !isUnique {
		o.cache.SetValue(uniqueFields["Email"], user.Email, true)
		return fmt.Errorf("user validation failure: %v", err)
	}
	if err := o.db.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("user creation failure in db: %v", err)
	}

	o.cache.SetValue(uniqueFields["Phone"], user.Phone, true)
	o.cache.SetValue(uniqueFields["Email"], user.Email, true)
	o.broadcastUserEvent("created", *user)
	return nil
}

func (o *UserService) ListUsers(ctx context.Context) ([]queries.User, error) {
	fmt.Println("t12e1SCscS2st")
	users, err := o.db.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("user retrieval failure in db: %w", err)
	}
	return users, nil
}

func (o *UserService) GetUserByID(ctx context.Context, userID string) (*model.UserEntity, error) {
	user, err := o.db.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user retrieval failure in db: %w", err)
	}

	return user, nil
}

func (o *UserService) DeleteUser(ctx context.Context, userID string) (*model.UserEntity, error) {

	user, err := o.db.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if err := o.db.DeleteUser(ctx, userID); err != nil {
		return nil, fmt.Errorf("user deletion failure in db: %w", err)
	}

	o.cache.DeleteField(uniqueFields["Phone"], user.Phone)
	o.cache.DeleteField(uniqueFields["Email"], user.Email)

	return user, nil
}

func (o *UserService) UpdateUser(ctx context.Context, userID string, userManager *model.UserEntity) (*model.UserEntity, error) {
	if err := o.db.UpdateUser(ctx, userID, userManager); err != nil {
		return nil, fmt.Errorf("user update failure in db: %w", err)
	}
	return userManager, nil
}

func (o *UserService) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	if err := o.wsHandler.HandleWebSocket(w, r); err != nil {
		http.Error(w, "Could not handle WebSocket", http.StatusBadRequest)
		return
	}
}

func (o *UserService) broadcastUserEvent(eventType string, user model.UserEntity) {
	event := map[string]interface{}{
		"event": eventType,
		"user":  user,
	}
	data, err := json.Marshal(event)
	if err != nil {
		fmt.Printf("broadcastUserEvent marshal failure: %v", err)
		return
	}
	o.wsHandler.Broadcast(data)
}
