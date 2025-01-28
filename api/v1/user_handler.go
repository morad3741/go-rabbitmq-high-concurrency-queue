package v1

import (
	"encoding/json"
	"hexArchitectureProject/internal/user"
	"net/http"
)

// UserHandler handles HTTP requests related to user operations.
type UserHandler struct {
	userService user.UserServiceInterface
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userService user.UserServiceInterface) *UserHandler {
	return &UserHandler{userService: userService}
}

// RegisterUser handles the user registration endpoint.
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the service to register the user
	if err := h.userService.AddUser(req.Name, req.Email, req.Password); err != nil {
		http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}
