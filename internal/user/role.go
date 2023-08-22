package user

import (
	"fmt"
)

const (
	roleUser  = "ROLE_USER"
	roleAdmin = "ROLE_ADMIN"
)

type Role struct {
	role string
}

func TryFrom(role string) (*Role, error) {
	allowedRoles := map[string]string{
		"ROLE_USER":  "ROLE_USER",
		"ROLE_ADMIN": "ROLE_ADMIN",
	}

	if value, ok := allowedRoles[role]; !ok {
		return nil, fmt.Errorf("role %s is not allowed", role)
	} else {
		return &Role{role: value}, nil
	}
}

func CreateUserRole() *Role {
	return &Role{
		role: roleUser,
	}
}

func (r *Role) CreateAdminRole() *Role {
	return &Role{
		role: roleAdmin,
	}
}

func (r *Role) GetRole() string {
	return r.role
}

func (r Role) String() string {
	return r.role
}

func RoleFromString(roleStr string) (Role, error) {
	role, err := TryFrom(roleStr)
	if err != nil {
		return Role{}, err
	}
	return *role, nil
}
