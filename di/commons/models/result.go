package models

import (
	"encoding/json"
)

type Result struct {
	Code   string          `json:"code"`
	Msg    string          `json:"msg"`
	Result json.RawMessage `json:"result"`
}

func GetResultData(b []byte, t interface{}) error {
	var res Result
	err := json.Unmarshal(b, &res)
	if err != nil {
		return err
	}
	messages := res.Result
	err = json.Unmarshal([]byte(messages), &t)
	if err != nil {
		return err
	}
	return nil
}

func GetResultCode(b []byte) (code string, r error) {
	var res Result
	err := json.Unmarshal(b, &res)
	if err != nil {
		return "", err
	}
	return code, nil
}
