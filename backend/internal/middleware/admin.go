package middleware

import (
	"eventBookingSystem/internal/auth/roles"
	"net/http"
)

// RequirePermission creates a middleware that checks for a specific permission
func RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the user role from the context
			role, ok := r.Context().Value(UserRoleKey).(string)
			if !ok {
				http.Error(w, "Unauthorized: Missing role", http.StatusUnauthorized)
				return
			}

			// Check if the role has the required permission
			if !roles.HasPermission(role, permission) {
				http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// AdminMiddleware checks if the user has admin role
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(UserRoleKey).(string)
		if !ok {
			http.Error(w, "Unauthorized: Missing role", http.StatusUnauthorized)
			return
		}

		if role != roles.RoleAdmin {
			http.Error(w, "Forbidden: Admin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
