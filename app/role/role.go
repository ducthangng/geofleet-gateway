package role

type Role int

const (
	RoleUnknown  Role = iota
	RoleCustomer      // 1
	RoleDriver        // 2
	RoleAdmin         // 3
)

func (r Role) String() string {
	switch r {
	case RoleCustomer:
		return "USER_ROLE_PASSENGER"
	case RoleDriver:
		return "USER_ROLE_DRIVER"
	case RoleAdmin:
		return "USER_ROLE_UNSPECIFIED"
	default:
		return "ADMIN"
	}
}

func ConvertRole(role int) Role {
	switch role {
	case 1:
		return RoleCustomer
	case 2:
		return RoleDriver
	case 3:
		return RoleAdmin
	}

	return RoleUnknown
}
