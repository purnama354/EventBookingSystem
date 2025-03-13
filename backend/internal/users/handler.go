package users

import (
	"encoding/json"
	"eventBookingSystem/internal/middleware"
	"eventBookingSystem/internal/types"
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

// Update the Register method to use standardized responses
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		types.SendError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
		return
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		types.SendError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	// Input validation with detailed error messages
	validationErrors := make(map[string]string)

	if strings.TrimSpace(req.Username) == "" {
		validationErrors["username"] = "Username is required"
	} else if len(req.Username) < 3 {
		validationErrors["username"] = "Username must be at least 3 characters long"
	}

	if strings.TrimSpace(req.Email) == "" {
		validationErrors["email"] = "Email is required"
	} else if _, err := mail.ParseAddress(req.Email); err != nil {
		validationErrors["email"] = "Invalid email format"
	}

	if len(req.Password) < 8 {
		validationErrors["password"] = "Password must be at least 8 characters long"
	}

	if len(validationErrors) > 0 {
		types.SendError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed", validationErrors)
		return
	}

	user, err := h.UserService.CreateUser(req.Username, req.Email, req.Password, "user")
	if err != nil {
		types.SendError(w, http.StatusInternalServerError, "REGISTRATION_FAILED", "Failed to create user", nil)
		return
	}

	// Return just the public user info, not including password hash
	userData := map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	}

	types.SendSuccess(w, http.StatusCreated, "User registered successfully", userData)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		types.SendError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		types.SendError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	// Validate input
	validationErrors := make(map[string]string)

	if strings.TrimSpace(req.Email) == "" {
		validationErrors["email"] = "Email is required"
	}

	if strings.TrimSpace(req.Password) == "" {
		validationErrors["password"] = "Password is required"
	}

	if len(validationErrors) > 0 {
		types.SendError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed", validationErrors)
		return
	}

	user, err := h.UserService.Login(req.Email, req.Password)
	if err != nil {
		types.SendError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password", nil)
		return
	}

	// Generate JWT token
	token, err := generateJWT(user.ID, user.Role == "admin")
	if err != nil {
		types.SendError(w, http.StatusInternalServerError, "TOKEN_GENERATION_FAILED", "Failed to generate token", nil)
		return
	}

	// Return user data and token
	responseData := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	}

	types.SendSuccess(w, http.StatusOK, "Login successful", responseData)
}

// Updated GetProfile handler
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the request context
	userID := r.Context().Value(middleware.UserIDKey).(string)

	// Get the user from the database
	user, err := h.UserService.GetUserByID(userID)
	if err != nil {
		types.SendError(w, http.StatusNotFound, "USER_NOT_FOUND", "User not found", nil)
		return
	}

	// Return user data without password hash
	userData := map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	}

	types.SendSuccess(w, http.StatusOK, "Profile retrieved successfully", userData)
}

func generateJWT(userID string, isAdmin bool) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set in .env file")
		return "", fmt.Errorf("JWT_SECRET is not set")
	}

	role := "user"
	if isAdmin {
		role = "admin"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
		"role":   role,
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Updated CreateAdmin handler
func (h *UserHandler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		types.SendError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
		return
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		types.SendError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	// Input validation with detailed error messages
	validationErrors := make(map[string]string)

	if strings.TrimSpace(req.Username) == "" {
		validationErrors["username"] = "Username is required"
	} else if len(req.Username) < 3 {
		validationErrors["username"] = "Username must be at least 3 characters long"
	}

	if strings.TrimSpace(req.Email) == "" {
		validationErrors["email"] = "Email is required"
	} else if _, err := mail.ParseAddress(req.Email); err != nil {
		validationErrors["email"] = "Invalid email format"
	}

	if len(req.Password) < 8 {
		validationErrors["password"] = "Password must be at least 8 characters long"
	}

	if len(validationErrors) > 0 {
		types.SendError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed", validationErrors)
		return
	}

	user, err := h.UserService.CreateUser(req.Username, req.Email, req.Password, "admin")
	if err != nil {
		types.SendError(w, http.StatusInternalServerError, "ADMIN_CREATION_FAILED", "Failed to create admin user", nil)
		return
	}

	// Return admin user data without password hash
	userData := map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	}

	types.SendSuccess(w, http.StatusCreated, "Admin user created successfully", userData)
}

// Updated Setup handler
func (h *UserHandler) Setup(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		types.SendError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
		return
	}

	// Check if system is already initialized
	users, err := h.UserService.GetAllUsers()
	if err != nil {
		types.SendError(w, http.StatusInternalServerError, "SYSTEM_CHECK_FAILED", "Failed to check system initialization status", nil)
		return
	}

	if len(users) > 0 {
		responseData := map[string]interface{}{
			"initialized": true,
		}
		types.SendSuccess(w, http.StatusOK, "System is already initialized", responseData)
		return
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		types.SendError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format", nil)
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
		types.SendError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed", validationErrors)
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
		types.SendError(w, http.StatusInternalServerError, "ADMIN_CREATION_FAILED", "Failed to create admin user", nil)
		return
	}

	// Generate token for the new admin
	token, err := generateJWT(user.ID, true)
	if err != nil {
		types.SendError(w, http.StatusInternalServerError, "TOKEN_GENERATION_FAILED", "Failed to generate authentication token", nil)
		return
	}

	// Return success response with user details and token
	responseData := map[string]interface{}{
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
		"token": token,
	}

	types.SendSuccess(w, http.StatusCreated, "System initialized successfully", responseData)
}
