/*
 * Copyright 2017-2018 IBM Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
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

package rest_impl

import (
	//"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	//"webank/DI/restapi"
)

func TestLoadManifestV1(t *testing.T) {
	data, err := ioutil.ReadFile("../../tests/testdata/broken-manifest/manifest.yml")
	assert.NoError(t, err)
	m, err := LoadManifestV1(data)
	assert.Nil(t, m)
	assert.Error(t, err)

	data, err = ioutil.ReadFile("../../tests/testdata/caffe-mnist-model/manifest.yml")
	assert.NoError(t, err)
	m, err = LoadManifestV1(data)
	assert.NoError(t, err)
	assert.NotNil(t, m)
}
//
//func TestConvertMemoryFromManifest(t *testing.T) {
//
//	f, unit, _ := server.convertMemoryFromManifest("256GB")
//	fmt.Printf("f:%v , %v ", f, unit.String())
//}
