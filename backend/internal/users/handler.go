package users

import (
	"encoding/json"
	"eventBookingSystem/internal/middleware"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
	UserService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Input validation
	if strings.TrimSpace(req.Username) == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Email) == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 8 {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	user, err := h.UserService.CreateUser(req.Username, req.Email, req.Password, "user")
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.UserService.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := generateJWT(user.ID, user.Role == "admin")
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the request context
	userID := r.Context().Value(middleware.UserIDKey).(string)

	// Get the user from the database
	user, err := h.UserService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func generateJWT(userID string, isAdmin bool) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set in .env file")
		return "", fmt.Errorf("JWT_SECRET is not set")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":  userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"isAdmin": isAdmin,                               //  isAdmin claim
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (h *UserHandler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Input validation
	if strings.TrimSpace(req.Username) == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Email) == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 8 {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	user, err := h.UserService.CreateUser(req.Username, req.Email, req.Password, "admin")
	if err != nil {
		http.Error(w, "Failed to create admin user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Setup(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if system is already initialized
	users, err := h.UserService.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to check system initialization status", http.StatusInternalServerError)
		return
	}

	if len(users) > 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":     false,
			"message":     "System is already initialized",
			"initialized": true,
		})
		return
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Enhanced validation
	validationErrors := make(map[string]string)

	if username := strings.TrimSpace(req.Username); username == "" {
		validationErrors["username"] = "Username is required"
	} else if len(username) < 3 {
		validationErrors["username"] = "Username must be at least 3 characters"
	}

	if email := strings.TrimSpace(req.Email); email == "" {
		validationErrors["email"] = "Email is required"
	} else if _, err := mail.ParseAddress(email); err != nil {
		validationErrors["email"] = "Invalid email format"
	}

	if password := req.Password; password == "" {
		validationErrors["password"] = "Password is required"
	} else if len(password) < 8 {
		validationErrors["password"] = "Password must be at least 8 characters"
	}

	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"errors":  validationErrors,
		})
		return
	}

	// Create the initial admin user
	user, err := h.UserService.CreateUser(
		req.Username,
		req.Email,
		req.Password,
		"admin", // First user is always admin
	)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to create admin user",
			"details": err.Error(),
		})
		return
	}

	// Generate token for the new admin
	token, err := generateJWT(user.ID, true)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to generate authentication token",
			"details": err.Error(),
		})
		return
	}

	// Return success response with user details and token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "System initialized successfully",
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
		"token": token,
	})
}
