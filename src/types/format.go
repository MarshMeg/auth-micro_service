package types

import (
	"strconv"
)

func StrToInt(i string) (int, error) {
	id, err := strconv.Atoi(i)
	return id, err
}
