package model

import (
	"encoding/json"
	"fmt"
)

func StructToMap(s interface{}) map[string]interface{} {
	if s == nil {
		return nil
	}
	var s_map map[string]interface{}
	fmt.Println(s)
	switch s.(type) {
	case User:
		s = s.(User)
	case Article:
		s = s.(Article)
	case Stat:
		s = s.(Stat)
	case Comment:
		s = s.(Comment)
	default:
		return nil
	}
	data, _ := json.Marshal(&s)
	json.Unmarshal(data, &s_map)
	return s_map
}
