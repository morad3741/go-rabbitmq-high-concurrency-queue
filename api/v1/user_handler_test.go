package v1

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockUserService is a mock implementation of the UserService interface.
type MockUserService struct{}

func (m *MockUserService) AddUser(name, email, password string) error {
	// Simulate a successful user registration
	return nil
}

// TestRegisterUser_Success tests the successful user registration scenario.
func TestRegisterUser_Success(t *testing.T) {
	// Create a mock user service
	mockUserService := &MockUserService{}

	// Create the user handler with the mock service
	userHandler := NewUserHandler(mockUserService)

	// Create the request body for user registration
	reqBody := map[string]string{
		"name":     "John Doe",
		"email":    "johndoe@example.com",
		"password": "securepassword",
	}
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Error marshalling request body: %v", err)
	}

	// Create the HTTP request using httptest
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", ioutil.NopCloser(bytes.NewReader(reqBodyJSON)))

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the RegisterUser handler
	userHandler.RegisterUser(rr, req)

	// Assert the status code is 201 Created
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, status)
	}

	// Assert the response body is as expected
	expectedBody := "User registered successfully"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, rr.Body.String())
	}
}
