package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/upekZ/rest-api-go/internal/database/models"
	"github.com/upekZ/rest-api-go/internal/database/queries"
	"github.com/upekZ/rest-api-go/internal/database/sqlc"
	"github.com/upekZ/rest-api-go/internal/types"
	"os"
	"runtime"
	"time"
)

type PostgresConn struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func NewPostgresConn() (*PostgresConn, error) {

	dsn := os.Getenv("DATABASE_DSN")

	if dsn == "" {
		dsn = "host=localhost port=5432 user=postgres dbname=postgres password=justadummy sslmode=disable"
	}

	poolconfig, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		return nil, err
	}

	poolconfig.MaxConns = int32(runtime.NumCPU() * 2)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolconfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	queries := sqlc.New(pool)

	return &PostgresConn{
		pool:    pool,
		queries: queries,
	}, nil
}

func (pgConn *PostgresConn) CreateUser(ctx context.Context, user *types.UserManager) error {

	tx, err := pgConn.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer handleRollBack(ctx, tx)

	params := user.SetUserParams()
	err = pgConn.queries.WithTx(tx).CreateUser(ctx,
		queries.CreateUserParams{FirstName: params.FirstName, LastName: params.LastName,
			Email: params.Email, Phone: params.Phone, Age: params.Age, Status: params.Status})

	if err != nil {
		return fmt.Errorf("error in user creation: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("DB commit error: %w", err)
	}

	return nil
}

func (pgConn *PostgresConn) GetUsers(ctx context.Context) ([]models.User, error) {

	users, err := pgConn.queries.ListUsers(ctx)

	if err != nil {
		return nil, fmt.Errorf("DB select error: %w", err)
	}

	return users, nil
}

func (pgConn *PostgresConn) UpdateUser(ctx context.Context, uID string, user *types.UserManager) error {

	tx, err := pgConn.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer handleRollBack(ctx, tx)

	var uuidVal pgtype.UUID
	err = uuidVal.Scan(uID)
	if err != nil {
		return fmt.Errorf("user id parsing failure: %w", err)
	}

	params := user.SetUserParams()

	err = pgConn.queries.WithTx(tx).UpdateUser(ctx,
		queries.UpdateUserParams{FirstName: params.FirstName, LastName: params.LastName,
			Email: params.Email, Phone: params.Phone, Age: params.Age, Status: params.Status, Userid: uuidVal})

	if err != nil {
		return fmt.Errorf("user update failure: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("DB commit error: %w", err)
	}

	return err
}

func (pgConn *PostgresConn) DeleteUser(ctx context.Context, id string) error {

	tx, err := pgConn.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer handleRollBack(ctx, tx)

	var uuidVal pgtype.UUID
	err = uuidVal.Scan(id)
	if err != nil {
		return fmt.Errorf("user id parsing failure: %w", err)
	}

	_, err = pgConn.queries.GetUser(ctx, uuidVal)
	if err != nil {
		return fmt.Errorf("user: [%s] not found. error: %w", id, err)
	}

	err = pgConn.queries.WithTx(tx).DeleteUser(context.Background(), uuidVal)
	if err != nil {
		return fmt.Errorf("user deletion failure: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("DB commit error: %w", err)
	}

	return err
}

func (pgConn *PostgresConn) GetUserByID(ctx context.Context, id string) (*types.UserManager, error) {

	var uuidVal pgtype.UUID
	if err := uuidVal.Scan(id); err != nil {
		return nil, fmt.Errorf("user id parsing failure: %w", err)
	}

	user, err := pgConn.queries.GetUser(ctx, uuidVal)
	if err != nil {
		return nil, fmt.Errorf("query execution failure for account [%s] error: %w", id, err)
	}
	userManager := types.CreateUserMgrFromParams(&user)
	return userManager, nil
}

func handleRollBack(ctx context.Context, trx pgx.Tx) {
	if err := trx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		fmt.Printf("transaction roll-back failure: %v", err)
	}
}
