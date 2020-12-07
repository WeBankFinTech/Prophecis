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
	"strconv"
	"testing"
	"webank/AIDE/notebook-server/pkg/models"
)

func TestName(t *testing.T) {
	var l []models.Notebook

	for i := 0; i < 10; i++ {
		notebook := models.Notebook{
			Name: strconv.Itoa(i),
		}
		l = append(l, notebook)
	}

	//pages, _ := GetSubListByPageAndSize(l, 2, 3)

	//fmt.Printf("origin: %+v, sub: %+v", l, pages)

}

func TestSlicesplit(t *testing.T) {
	ints := []int{0, 1, 2, 3, 4}
	i := ints[4:5]
	println(i)

}
