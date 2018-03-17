3package core

import (
	"strconv"
)

type Id uint64

func GetIdFromString(idStr string) (id Id, err error) {

	id, err = strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return
	}

	return
}
