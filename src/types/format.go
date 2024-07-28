package types

import (
	"strconv"
)

func StrToInt(i string) int {
	id, err := strconv.Atoi(i)
	if err != nil {
		return 0
	}
	return id
}
