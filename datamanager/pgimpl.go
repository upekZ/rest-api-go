package datamanager

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/upekZ/rest-api-go/sqlc"
)

type PostgresConn struct {
	queries *sqlc.Queries
}

func NewPostgresConn() (*PostgresConn, error) {

	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		dsn = "host=localhost port=5432 user=postgres dbname=postgres password=justadummy sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	queries := sqlc.New(pool)

	return &PostgresConn{
		queries: queries,
	}, nil
}

func (pgConn *PostgresConn) CreateUser(user *UserManager) error {

	params := user.SetUserParams()
	err := pgConn.queries.CreateUser(context.Background(),
		sqlc.CreateUserParams{FirstName: params.FirstName, LastName: params.LastName,
			Email: params.Email, Phone: params.Phone, Age: params.Age, Status: params.Status})

	if err != nil {
		return fmt.Errorf("error in user creation: %w", err)
	}

	return nil
}

func (pgConn *PostgresConn) GetUsers() ([]sqlc.User, error) {

	users, err := pgConn.queries.ListUsers(context.Background())

	if err != nil {
		return nil, fmt.Errorf("DB select error: %w", err)
	}

	return users, nil
}

func (pgConn *PostgresConn) UpdateUser(uID string, user *UserManager) error {

	var uuidVal pgtype.UUID
	err := uuidVal.Scan(uID)
	if err != nil {
		return fmt.Errorf("user id parsing failure: %w", err)
	}

	params := user.SetUserParams()

	err = pgConn.queries.UpdateUser(context.Background(),
		sqlc.UpdateUserParams{FirstName: params.FirstName, LastName: params.LastName,
			Email: params.Email, Phone: params.Phone, Age: params.Age, Status: params.Status, Userid: uuidVal})

	if err != nil {
		return fmt.Errorf("user update failure: %w", err)
	}

	return err
}

func (pgConn *PostgresConn) DeleteUser(id string) error {

	var uuidVal pgtype.UUID
	err := uuidVal.Scan(id)
	if err != nil {
		return fmt.Errorf("user id parsing failure: %w", err)
	}

	err = pgConn.queries.DeleteUser(context.Background(), uuidVal)
	if err != nil {
		return fmt.Errorf("user deletion error: %w", err)
	}
	return err
}

func (pgConn *PostgresConn) GetUserByID(id string) (*UserManager, error) {

	var uuidVal pgtype.UUID
	err := uuidVal.Scan(id)
	if err != nil {
		return nil, fmt.Errorf("user id parsing failure: %w", err)
	}

	user, err := pgConn.queries.GetUser(context.Background(), uuidVal)
	if err != nil {
		return nil, fmt.Errorf("query execusion failure for account [%s] error: %w", id, err)
	}
	userManager := CreateUserMgrFromParams(&user)
	return userManager, nil
}
