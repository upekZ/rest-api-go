package services

import (
	"context"
	"encoding/json"
	"github.com/upekZ/rest-api-go/internal/database/queries"
	"github.com/upekZ/rest-api-go/internal/model"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserDB struct {
	mock.Mock
}

func (m *MockUserDB) CreateUser(ctx context.Context, user *model.UserEntity) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserDB) UpdateUser(ctx context.Context, uID string, user *model.UserEntity) error {
	args := m.Called(uID, user)
	return args.Error(1)
}

func (m *MockUserDB) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserDB) GetUserByID(ctx context.Context, id string) (*model.UserEntity, error) {
	args := m.Called(id)
	return args.Get(0).(*model.UserEntity), args.Error(1)
}

func (m *MockUserDB) GetUsers(context.Context) ([]queries.User, error) {
	//ToDo
	return nil, nil
}

func (m *MockUserDB) IsEmailUnique(context.Context, string) (bool, error) {
	//ToDo
	return true, nil
}
func (m *MockUserDB) IsPhoneUnique(context.Context, string) (bool, error) {
	//ToDo
	return true, nil
}

// MockWebSocketService mocks the WebSocketService interface.
type MockWebSocketService struct {
	mock.Mock
}

func (m *MockWebSocketService) Broadcast(message []byte) {
	m.Called(message)
}

func (m *MockWebSocketService) Run() {
}

func (m *MockWebSocketService) HandleWebSocket(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type MockUserCache struct {
	mock.Mock
}

func (c *MockUserCache) IsValueTaken(key string, value string) bool {
	//ToDo
	return false
}

func (c *MockUserCache) DeleteField(key string, value string) {
	//ToDo
}

func (c *MockUserCache) SetValue(key string, value string, exists bool) {
	//ToDo
}

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name          string
		user          *model.UserEntity
		dbResult      model.UserEntity
		dbError       error
		wsCalled      bool
		expectedError string
	}{
		//Most validation is done at service level. --> Less chance to propagate invalid fields to DB
		{
			name:     "Valid user",
			user:     &model.UserEntity{FirstName: "namefirst", LastName: "namelast", Email: "valid@accepted.com", Phone: "17689899899", Age: 15, Status: "Active"},
			dbResult: model.UserEntity{FirstName: "namefirst", LastName: "namelast", Email: "valid@accepted.com", Phone: "17689899899", Age: 15, Status: "Active"},
			dbError:  nil,
			wsCalled: true,
		},
		{
			name:          "Empty name",
			user:          &model.UserEntity{FirstName: "", LastName: "name_last", Email: "valid@accepted.com", Phone: "0768989899", Age: 15, Status: "Active"},
			expectedError: "user validation failure: invalid entries for fields",
		},
		{
			name:          "Invalid email",
			user:          &model.UserEntity{FirstName: "name_first", LastName: "name_last", Email: "invalid.com", Phone: "0768989899", Age: 15, Status: "Active"},
			expectedError: "user validation failure: invalid entries for fields",
		},
		{
			name:          "Status",
			user:          &model.UserEntity{FirstName: "name_first", LastName: "name_last", Email: "invalid.com", Phone: "0768989899", Age: 10, Status: "pending"},
			expectedError: "user validation failure: invalid entries for fields",
		},
		{
			name:          "Database error",
			user:          &model.UserEntity{FirstName: "name_first", Email: "invalid.com", Phone: "0768989899", Age: 15, Status: "Active"},
			expectedError: "user validation failure: invalid entries for fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockUserDB{}
			mockWS := &MockWebSocketService{}
			mockCache := &MockUserCache{}
			service := NewUserService(mockDB, mockCache, mockWS)

			if tt.name == "Valid user" {
				mockDB.On("CreateUser", mock.Anything, mock.AnythingOfType("*model.UserEntity")).Return(nil)
			}

			if tt.wsCalled {
				event := map[string]interface{}{"event": "created", "user": tt.user}
				data, _ := json.Marshal(event)
				mockWS.On("Broadcast", data).Once()
			}

			err := service.CreateUser(context.Background(), tt.user)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			mockDB.AssertExpectations(t)
			mockWS.AssertExpectations(t)
		})
	}
}
