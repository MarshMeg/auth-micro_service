package roles

import (
	"fmt"
	"reflect"
)

func List() Roles {
	return Roles{
		Admin:   Admin(),
		Service: Service(),
		Member:  Member(),
	}
}

func Admin() UserRole {
	return UserRole{
		RoleName:  "Admin",
		AccessLvl: 10,
	}
}

func Service() UserRole {
	return UserRole{
		RoleName:  "Service",
		AccessLvl: 5,
	}
}

func Member() UserRole {
	return UserRole{
		RoleName:  "Member",
		AccessLvl: 1,
	}
}

func (roles *Roles) GetRole(filters *UserRole) UserRole {
	var r = reflect.TypeOf(roles)
	for i := 0; i < r.NumField(); i++ {
		filed := r.Field(i)
		fmt.Println(filed)
		if filters.RoleName == r.Field(i).Name {
			return UserRole{}
		}
	}
	return UserRole{}
}
