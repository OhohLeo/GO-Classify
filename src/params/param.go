package params

import (
	"encoding/json"
)

type HasParam interface {
	GetParam(string, json.RawMessage) (interface{}, error)
}
