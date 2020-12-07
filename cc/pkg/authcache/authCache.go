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
package authcache

import (
	"github.com/patrickmn/go-cache"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/logger"
	"time"
)


var TokenCache *cache.Cache

func Init() {
	cacheDuration, err := time.ParseDuration(common.GetAppConfig().Core.Cache.CacheTime)
	if nil != err {
		logger.Logger().Errorf("failed to init cache")
	}

	TokenCache = cache.New(cacheDuration, cacheDuration)

}
