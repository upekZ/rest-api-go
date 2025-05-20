package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	dsn := "host=localhost port=5433 user=postgres password=justadummy dbname=testdb sslmode=disable"
	const maxRetries = 10
	const retryDelay = 2 * time.Second
	var pool *pgxpool.Pool

	for attempt := 1; attempt <= maxRetries; attempt++ {
		var err error
		pool, err = pgxpool.New(context.Background(), dsn)
		if err == nil {
			err = pool.Ping(context.Background())
			if err == nil {
				break
			}
		}
		time.Sleep(retryDelay)
	}

	if pool == nil {
		panic("Failed to connect to test database")
	}
	defer pool.Close()

	_, err := pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS "user" (
            userId SERIAL PRIMARY KEY,
            first_name VARCHAR(50),
            last_name VARCHAR(50),
            email VARCHAR(100),
            phone VARCHAR(20),
            age INTEGER,
            "status" VARCHAR(20)
        );
        INSERT INTO "user" (first_name, last_name, email, phone, age, "status")
        VALUES ('upeka', 'W', 'test.create@example.com', '09809889214', 29, 'active');
    `)
	if err != nil {
		panic("Failed to initialize test data: " + err.Error())
	}
	m.Run()
}

func TestGetUser(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:3001/users", nil)
	assert.NoError(t, err)

	client := &http.Client{}
	response, err := client.Do(req)
	assert.NoError(t, err)
	defer response.Body.Close()

	assert.Equal(t, http.StatusOK, response.StatusCode)

	var users []map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&users)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, len(users), 1)
	assert.Equal(t, "upeka", users[0]["first_name"])
	assert.Equal(t, 29, users[0]["age"])
}

func TestCreateUser(t *testing.T) {
	user := map[string]interface{}{
		"first_name": "create",
		"last_name":  "test",
		"email":      "create.test@test.com",
		"phone":      "90129210921",
		"age":        25,
		"status":     "active",
	}
	body, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "http://localhost:3001/users", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	assert.NoError(t, err)
	defer response.Body.Close()

	assert.Equal(t, http.StatusCreated, response.StatusCode)

	pool, err := pgxpool.New(context.Background(), "host=localhost port=5433 user=postgres password=justadummy dbname=testdb sslmode=disable")
	assert.NoError(t, err)
	defer pool.Close()

	var count int
	err = pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM "user" WHERE email = $1`, "create.test@test.com").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}
