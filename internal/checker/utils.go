package checker

import (
	"encoding/json"

	"gorm.io/datatypes"
)

func mustMarshalJSON(v interface{}) datatypes.JSON {
	data, err := json.Marshal(v)
	if err != nil {
		return datatypes.JSON("{}")
	}
	return datatypes.JSON(data)
}
