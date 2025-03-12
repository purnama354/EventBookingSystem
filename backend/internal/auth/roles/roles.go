package roles

import "slices"

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

const (
	PermissionReadEvents     = "events:read"
	PermissionCreateEvents   = "events:create"
	PermissionUpdateEvents   = "events:update"
	PermissionDeleteEvents   = "events:delete"
	PermissionManageUsers    = "users:manage"
	PermissionCreateBookings = "bookings:create"
	PermissionReadBookings   = "bookings:read"
	PermissionCancelBookings = "bookings:cancel"
)

var RolePermissions = map[string][]string{
	RoleUser: {
		PermissionReadEvents,
		PermissionCreateBookings,
		PermissionReadBookings,
		PermissionCancelBookings,
	},
	RoleAdmin: {
		PermissionReadEvents,
		PermissionCreateEvents,
		PermissionUpdateEvents,
		PermissionDeleteEvents,
		PermissionManageUsers,
		PermissionReadBookings,
		PermissionCreateBookings,
		PermissionCancelBookings,
	},
}

func HasPermission(role, permission string) bool {
	permissions, exists := RolePermissions[role]
	if !exists {
		return false
	}

	return slices.Contains(permissions, permission)
}
