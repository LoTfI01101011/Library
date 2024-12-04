package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/LoTfI01101011/Library/controllers"
	"github.com/LoTfI01101011/Library/initial"
	"github.com/LoTfI01101011/Library/middleware"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/stretchr/testify/assert"
)

type TestUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// type TestBook struct {
// 	Title       string `json:"title"`
// 	Author      string `json:"author"`
// 	Pages       int    `json:"pages"`
// 	Description string `json:"description"`
// }

var token string

func init() {
	initial.ConnectToDb()
	initial.InitRedis()
	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost:8000/api/google/auth/callback?provider=google"),
	)
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Middleware and routes setup
	r.POST("/api/auth/register", controllers.SignUpUser)
	r.POST("/api/auth/login", controllers.LoginUser)
	r.POST("/api/auth/logout", middleware.AuthMiddelware, controllers.Logout)
	r.POST("/api/book", middleware.AuthMiddelware, controllers.CreateBook)
	r.GET("/api/book", middleware.AuthMiddelware, controllers.GetBooks)
	r.GET("/api/book/:id", middleware.AuthMiddelware, controllers.GetBookById)
	r.PATCH("/api/book/:id", middleware.AuthMiddelware, controllers.UpdateBook)
	r.DELETE("/api/book/:id", middleware.AuthMiddelware, controllers.DeleteBook)

	return r
}

func loginAndGetToken(r *gin.Engine, t *testing.T) string {
	if token != "" {
		return token
	}

	w := httptest.NewRecorder()
	loginData := map[string]string{
		"email":    "testuser@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(loginData)

	req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "token")
	token = "Bearer " + response["token"].(string)

	return token
}

func sendRequest(r *gin.Engine, method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	reqBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, path, bytes.NewReader(reqBody))
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	r.ServeHTTP(w, req)
	return w
}

func TestSignUpUser(t *testing.T) {
	r := setupRouter()

	user := TestUser{
		Email:    "testuser@example.com",
		Password: "password123",
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	w := sendRequest(r, http.MethodPost, "/api/auth/register", user, headers)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "token")
}

func TestLoginAndLogoutUser(t *testing.T) {
	r := setupRouter()

	// Login
	user := TestUser{
		Email:    "testuser@example.com",
		Password: "password123",
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	loginResponse := sendRequest(r, http.MethodPost, "/api/auth/login", user, headers)
	assert.Equal(t, http.StatusOK, loginResponse.Code)

	var loginData map[string]string
	json.Unmarshal(loginResponse.Body.Bytes(), &loginData)
	assert.Contains(t, loginData, "token")

	token := "Bearer " + loginData["token"]

	// Logout
	logoutHeaders := map[string]string{
		"Authorization": token,
	}

	logoutResponse := sendRequest(r, http.MethodPost, "/api/auth/logout", nil, logoutHeaders)
	assert.Equal(t, http.StatusOK, logoutResponse.Code)
}

// func TestCreateBook(t *testing.T) {
// 	r := setupRouter()
// 	token := loginAndGetToken(r, t)
// 	log.Println(token)
// 	book := TestBook{
// 		Title:       "Go Programming",
// 		Author:      "John Doe",
// 		Pages:       350,
// 		Description: "A book about Go programming.",
// 	}

// 	headers := map[string]string{
// 		"Authorization": token,
// 		"Content-Type":  "application/json",
// 	}

// 	w := sendRequest(r, http.MethodPost, "/api/book", book, headers)
// 	assert.Equal(t, http.StatusOK, w.Code)
// 	var response map[string]interface{}
// 	json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.Contains(t, response, "book")
// }

// func TestGetBooks(t *testing.T) {
// 	r := setupRouter()
// 	token := loginAndGetToken(r, t)

// 	headers := map[string]string{
// 		"Authorization": token,
// 	}

// 	w := sendRequest(r, http.MethodGet, "/api/book", nil, headers)
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var response map[string]interface{}
// 	json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.Contains(t, response, "books")
// }

// func TestGetBookById(t *testing.T) {
// 	r := setupRouter()
// 	token := loginAndGetToken(r, t)
// 	bookID := "mock_book_id"

// 	headers := map[string]string{
// 		"Authorization": token,
// 	}

// 	w := sendRequest(r, http.MethodGet, "/api/book/"+bookID, nil, headers)
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var response map[string]interface{}
// 	json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.Contains(t, response, "book")
// }

// func TestUpdateBook(t *testing.T) {
// 	r := setupRouter()
// 	token := loginAndGetToken(r, t)
// 	bookID := "mock_book_id"

// 	updateData := map[string]interface{}{
// 		"Title": "Updated Title",
// 	}

// 	headers := map[string]string{
// 		"Authorization": token,
// 		"Content-Type":  "application/json",
// 	}

// 	w := sendRequest(r, http.MethodPatch, "/api/book/"+bookID, updateData, headers)
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var response map[string]interface{}
// 	json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.Contains(t, response, "book")
// 	assert.Equal(t, "Updated Title", response["book"].(map[string]interface{})["Title"])
// }

// func TestDeleteBook(t *testing.T) {
// 	r := setupRouter()
// 	token := loginAndGetToken(r, t)
// 	bookID := "mock_book_id"

// 	headers := map[string]string{
// 		"Authorization": token,
// 	}

// 	w := sendRequest(r, http.MethodDelete, "/api/book/"+bookID, nil, headers)
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var response map[string]interface{}
// 	json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.Contains(t, response, "response")
// 	assert.Equal(t, "the book was deleted succesfuly", response["response"])
// }
