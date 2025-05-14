package types

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/upekZ/rest-api-go/internal/sqlc"
	"regexp"
)

type UserManager struct {
	UID       string          `json:"userId"`
	FirstName string          `json:"firstName"`
	LastName  string          `json:"lastName"`
	Email     string          `json:"email"`
	Phone     string          `json:"phone"`
	Age       uint32          `json:"age"`
	Status    sqlc.UserStatus `json:"status"`
}

func ValidateUser(user *UserManager) bool {

	var validReqFields = (IsValidEmail(user.Email)) && IsValidName(user.FirstName) && IsValidName(user.LastName)

	if validReqFields {
		if !(IsValidPhone(user.Phone)) {
			user.Phone = ""
		}
	}

	return validReqFields
}

func (manager *UserManager) SetUserParams() *sqlc.User {
	user := sqlc.User{
		FirstName: manager.FirstName,
		LastName:  manager.LastName,
		Email:     manager.Email,
		Phone: pgtype.Text{
			String: manager.Phone,
			Valid:  IsValidPhone(manager.Phone),
		},
		Age: pgtype.Int4{
			Int32: int32(manager.Age),
			Valid: true,
		},
		Status: sqlc.NullUserStatus{
			UserStatus: manager.Status,
			Valid:      manager.Status.Valid(),
		},
	}

	return &user
}

func CreateUserMgrFromParams(params *sqlc.User) *UserManager {
	return &UserManager{
		UID:       params.Userid.String(),
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Phone:     params.Phone.String,
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
