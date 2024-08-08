package roles

type Roles struct {
	Admin   UserRole
	Service UserRole
	Member  UserRole
}

type UserRole struct {
	RoleName  string
	AccessLvl int
}
