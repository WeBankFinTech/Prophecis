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
package utils

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestGetCPUValueWithUnits(t *testing.T) {
	units, e := GetCPUValueWithUnits("7800m")
	if e != nil {
		fmt.Println(e.Error())
	}
	fmt.Println(*units)
}

func TestTimeNow(t *testing.T) {
	now := time.Now()
	unix := now.Unix()
	//fmt.Println(string(unix))
	i := strconv.FormatInt(unix, 10)
	fmt.Println(i)
	fmt.Println(unix)
	format := time.Unix(unix, 0).Format(time.RFC1123)
	fmt.Println(format)
}
