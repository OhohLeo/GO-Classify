package core

import (
	"strconv"
)

type Id uint64

func GetIdFromString(idStr string) (id Id, err error) {

	idConv, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return
	}

	id = Id(idConv)
	return
}
