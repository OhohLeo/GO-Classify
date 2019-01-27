package params

import (
	"encoding/json"
)

type Param interface {
	GetName() string
	ExecuteParam(input json.RawMessage) (interface{}, error)
}

type HasParams interface {
	GetParams() []Param
}

func GetParams(input interface{}) map[string]Param {
	result := make(map[string]Param)

	params, ok := input.(HasParams)
	if ok == false {
		return result
	}

	for _, param := range params.GetParams() {
		result[param.GetName()] = param
	}
	return result
}
