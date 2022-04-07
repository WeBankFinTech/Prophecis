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

package plugin

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"testing"
)

func TestExc(t *testing.T) {
	//pathRegx := strings.ReplaceAll(strings.TrimSpace("/cc/v1/*"), "*", "(.*)")
	//
	//matched, err := regexp.MatchString(pathRegx, "/cc/v1/inter/auth")
	//if err != nil {
	//	fmt.Println(err.Error())
	//
	//}
	//fmt.Printf("bool: %+v", matched)

	re := regexp2.MustCompile(`^(hduser).*(?<!_c)$`, 0)
	matched, err := re.MatchString("hduseraha")
	if err != nil {
		fmt.Println(err.Error())

	}
	fmt.Printf("bool: %+v", matched)

}
