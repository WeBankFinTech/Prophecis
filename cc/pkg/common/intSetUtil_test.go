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
package common

import (
	"net/url"
	"testing"
)

func TestCas(t *testing.T) {
	u := &url.URL{
		Scheme:   "http",
		Host:     "aha",
		Path:     "/",
		RawQuery: "ticket=ST-136129-qMuhTiA2Ney27MJKblEc-sso.webank.com",
		Fragment: "/",
	}
	ticket := "ST-136129-qMuhTiA2Ney27MJKblEc-sso.webank.com"
	Auth(u, ticket)
}
