/*
 * Copyright 2020 WeBank
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package controlcenter

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
