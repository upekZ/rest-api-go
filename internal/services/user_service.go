package services

import (
	"context"
	"encoding/json"
	"fmt"
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
	GetUsers(context.Context) ([]model.UserEntity, error) //queries.User to be replaced with model.UserEntity
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
	if isUnique, _ := o.IsUniqueField(ctx, uniqueFields["Phone"], user.Phone); !isUnique {
		o.cache.SetValue(uniqueFields["Phone"], user.Phone, true)
		return fmt.Errorf("phone number already attached to a user")
	}

	if isUnique, _ := o.IsUniqueField(ctx, uniqueFields["Email"], user.Email); !isUnique {
		o.cache.SetValue(uniqueFields["Email"], user.Email, true)
		return fmt.Errorf("email already attached to a user")
	}
	if err := o.db.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("user creation failure in db")
	}

	o.cache.SetValue(uniqueFields["Phone"], user.Phone, true)
	o.cache.SetValue(uniqueFields["Email"], user.Email, true)
	o.broadcastUserEvent("created", *user)
	return nil
}

func (o *UserService) ListUsers(ctx context.Context) ([]model.UserEntity, error) {
	users, err := o.db.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("user retrieval failure in db")
	}
	return users, nil
}

func (o *UserService) GetUserByID(ctx context.Context, userID string) (*model.UserEntity, error) {
	user, err := o.db.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user retrieval failure in db")
	}

	return user, nil
}

func (o *UserService) DeleteUser(ctx context.Context, userID string) (*model.UserEntity, error) {

	user, err := o.db.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user [%s] not found", userID)
	}
	if err := o.db.DeleteUser(ctx, userID); err != nil {
		return nil, fmt.Errorf("user deletion failure in db")
	}

	o.cache.DeleteField(uniqueFields["Phone"], user.Phone)
	o.cache.DeleteField(uniqueFields["Email"], user.Email)

	return user, nil
}

func (o *UserService) UpdateUser(ctx context.Context, userID string, userManager *model.UserEntity) (*model.UserEntity, error) {
	if err := o.db.UpdateUser(ctx, userID, userManager); err != nil {
		return nil, fmt.Errorf("user update failure in db")
	}
	return userManager, nil
}

func (o *UserService) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	if err := o.wsHandler.HandleWebSocket(w, r); err != nil {
		http.Error(w, "websocket handler failure", http.StatusInternalServerError)
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
		fmt.Printf("json marshal for web-sockets failure:") // ToDo convert to logging
		return
	}
	o.wsHandler.Broadcast(data)
}
