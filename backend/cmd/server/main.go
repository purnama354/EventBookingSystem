package main

import (
	"eventBookingSystem/configs"
	"eventBookingSystem/internal/auth/roles"
	"eventBookingSystem/internal/bookings"
	"eventBookingSystem/internal/events"
	"eventBookingSystem/internal/middleware"
	"eventBookingSystem/internal/users"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
		return
	}

	db, err := configs.ConnectDB(config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return
	}
	db.AutoMigrate(&users.User{}, &events.Event{}, &bookings.Booking{})

	userRepository := users.NewUserRepository(db)
	userService := users.NewUserService(userRepository)
	userHandler := users.NewUserHandler(userService)

	eventRepository := events.NewEventRepository(db)
	eventService := events.NewEventService(eventRepository)
	eventHandler := events.NewEventHandler(eventService)

	bookingRepository := bookings.NewBookingRepository(db)
	bookingService := bookings.NewBookingService(bookingRepository)
	bookingHandler := bookings.NewBookingHandler(bookingService)

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/api/setup", userHandler.Setup)
	mux.HandleFunc("/api/users/register", userHandler.Register)
	mux.HandleFunc("/api/users/login", userHandler.Login)

	// Protected routes with specific permissions
	mux.Handle("/api/users/profile",
		middleware.AuthMiddleware(
			middleware.RequirePermission(roles.PermissionReadBookings)(
				http.HandlerFunc(userHandler.GetProfile),
			),
		),
	)

	mux.Handle("/api/admin/users/create",
		middleware.AuthMiddleware(
			middleware.RequirePermission(roles.PermissionManageUsers)(
				http.HandlerFunc(userHandler.CreateAdmin),
			),
		),
	)

	mux.Handle("/api/events",
		middleware.AuthMiddleware(
			middleware.RequirePermission(roles.PermissionCreateEvents)(
				http.HandlerFunc(eventHandler.HandleEvents),
			),
		),
	)

	mux.Handle("/api/bookings/",
		middleware.AuthMiddleware(
			middleware.RequirePermission(roles.PermissionCreateBookings)(
				http.HandlerFunc(bookingHandler.HandleBookings),
			),
		),
	)

	// CORS configuration
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Allow requests from your React app
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	// Wrap the ServeMux with the CORS handler
	handler := corsHandler.Handler(mux)
	http.Handle("/", middleware.LoggingMiddleware(handler))

	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
