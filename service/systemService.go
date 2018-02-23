package service

import (
  "encoding/json"
)

func ParseConfig(b []byte, dataStruct interface{}) error {
	return json.Unmarshal(b, dataStruct)
}
