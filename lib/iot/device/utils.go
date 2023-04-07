package device

import (
	"encoding/json"
)

func GetConfig(data []byte, config interface{}) error {
	if err := json.Unmarshal(data, config); err != nil {
		return err
	}

	if err := validate.Struct(config); err != nil {
		return err
	}

	return nil
}
