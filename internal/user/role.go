package user

const (
	roleUser  = "ROLE_USER"
	roleAdmin = "ROLE_ADMIN"
)

type Role struct {
	role string
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
