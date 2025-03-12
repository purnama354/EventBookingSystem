# API Documentation

## Users

- `POST /api/users/register`: Register a new user.
  - Request body:
    ```json
    {
      "username": "string",
      "email": "string",
      "password": "string"
    }
    ```
  - Response body:
    ```json
    {
      "id": "string",
      "username": "string",
      "email": "string"
    }
    ```
- `POST /api/users/login`: Login and get a JWT token.
  - Request body:
    ```json
    {
      "email": "string",
      "password": "string"
    }
    ```
  - Response body:
    ```json
    {
      "token": "string"
    }
    ```
- `GET /api/users/profile`: Get the user profile (requires authentication).
  - Request header:
    ```
    Authorization: Bearer <JWT token>
    ```
  - Response body:
    ```json
    {
      "id": "string",
      "username": "string",
      "email": "string"
    }
    ```

## Events

- `GET /api/events`: Get a list of events.
- `POST /api/events`: Create a new event (requires authentication and admin role).
  - Request header:
    ```
    Authorization: Bearer <JWT token>
    ```
  - Request body:
    ```json
    {
      "title": "string",
      "description": "string",
      "date": "string (RFC3339)",
      "location": "string",
      "capacity": "integer"
    }
    ```
- `GET /api/events/{eventID}`: Get event details.
- `PUT /api/events/{eventID}`: Update an event (requires authentication and admin role).
  - Request header:
    ```
    Authorization: Bearer <JWT token>
    ```
  - Request body:
    ```json
    {
      "title": "string",
      "description": "string",
      "date": "string (RFC3339)",
      "location": "string",
      "capacity": "integer"
    }
    ```
- `DELETE /api/events/{eventID}`: Delete an event (requires authentication and admin role).

## Bookings

- `POST /api/bookings`: Create a new booking (requires authentication).
  - Request header:
    ```
    Authorization: Bearer <JWT token>
    ```
  - Request body:
    ```json
    {
      "eventID": "string",
      "seats": "integer"
    }
    ```
- `GET /api/bookings/{bookingID}`: Get booking details (requires authentication).
  - Request header:
    ```
    Authorization: Bearer <JWT token>
    ```
- `GET /api/bookings/users/{userID}`: Get all bookings for a user (requires authentication).
  - Request header:
    ```
    Authorization: Bearer <JWT token>
    ```
- `DELETE /api/bookings/{bookingID}`: Cancel a booking (requires authentication).
  - Request header:
    ```
    Authorization: Bearer <JWT token>
    ```
