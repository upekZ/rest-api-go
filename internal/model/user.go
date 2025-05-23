package model

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/upekZ/rest-api-go/internal/database/queries"
	"regexp"
)

type dbUser = queries.User

type UserEntity struct {
	UID       string             `json:"userId"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Email     string             `json:"email"`
	Phone     string             `json:"phone"`
	Age       uint32             `json:"age"`
	Status    queries.UserStatus `json:"status"`
}

func ValidateUser(user *UserEntity) (bool, error) {

	if !(IsValidName(user.FirstName) && IsValidName(user.LastName)) {
		return false, fmt.Errorf("invalid format for name")
	}
	if !IsValidEmail(user.Email) {
		return false, fmt.Errorf("invalid email")
	}

	return true, nil
}

func (manager *UserEntity) SetUserParams() *queries.User {
	user := queries.User{
		FirstName: manager.FirstName,
		LastName:  manager.LastName,
		Email:     manager.Email,
		Phone:     manager.Phone,
		Age: pgtype.Int4{
			Int32: int32(manager.Age),
			Valid: true,
		},
		Status: queries.NullUserStatus{
			UserStatus: manager.Status,
			Valid:      manager.Status.Valid(),
		},
	}

	return &user
}

func CreateUserMgrFromParams(params *queries.User) *UserEntity {
	return &UserEntity{
		UID:       params.Userid.String(),
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Phone:     params.Phone,
		Age:       uint32(params.Age.Int32),
		Status:    params.Status.UserStatus,
	}
}

func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func IsValidPhone(number string) bool {
	re := regexp.MustCompile(`^(\+?\d{10,15})$`)
	return re.MatchString(number)
}

func IsValidName(name string) bool {
	re := regexp.MustCompile(`^[a-zA-Z\s'-]{2,50}$`)
	return re.MatchString(name)
}

func ConvertUsersToEntities(users []dbUser) []UserEntity {

	entities := make([]UserEntity, 0, len(users))

	for _, user := range users {
		entity := UserEntity{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Phone:     user.Phone,
		}

		if user.Userid.Valid {
			entity.UID = user.Userid.String()
		} else {
			entity.UID = ""
		}
		if user.Age.Valid {
			if user.Age.Int32 >= 0 {
				entity.Age = uint32(user.Age.Int32)
			} else {
				entity.Age = 0
			}
		} else {
			entity.Age = 0
		}

		if user.Status.Valid {
			entity.Status = user.Status.UserStatus
		} else {
			entity.Status = ""
		}

		entities = append(entities, entity)
	}

	return entities
}
