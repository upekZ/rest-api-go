package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestCreateUser(t *testing.T) {
	user := map[string]interface{}{
		"firstName": "create",
		"lastName":  "test",
		"email":     "create1.test@test.com",
		"phone":     "901129210921",
		"age":       25,
		"status":    "Active",
	}
	body, _ := json.Marshal(user)

	appHost := os.Getenv("APP_HOST")
	appPort := os.Getenv("APP_PORT")
	baseURL := fmt.Sprintf("http://%s:%s/users", appHost, appPort)

	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	assert.NoError(t, err)
	defer response.Body.Close()

	assert.Equal(t, http.StatusCreated, response.StatusCode)

	var createdUser map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&createdUser)
	assert.NoError(t, err)

	assert.Equal(t, "create", createdUser["firstName"])
	assert.Equal(t, float64(25), createdUser["age"])
}

func TestGetUser(t *testing.T) {
	appHost := os.Getenv("APP_HOST")
	appPort := os.Getenv("APP_PORT")
	baseURL := fmt.Sprintf("http://%s:%s/users", appHost, appPort)
	fmt.Println("Get User Logs")
	req, err := http.NewRequest("GET", baseURL, nil)
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
	assert.Equal(t, "create", users[0]["first_name"])
	assert.Equal(t, 25, users[0]["age"])
}
