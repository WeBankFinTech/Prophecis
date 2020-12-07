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
package middleware

import (
	"context"
	"mlss-controlcenter-go/pkg/authcache"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/constants"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"net/http"
)

func newContextWithRequestID(ctx context.Context, req *http.Request) context.Context {
	token := req.Header.Get(constants.AUTH_HEADER_TOKEN)
	logger.Logger().Debugf("middleware to get token: %v", token)
	var sessionUser *models.SessionUser
	if "" != token {
		get, b := authcache.TokenCache.Get(token)
		logger.Logger().Debugf("middleware to get user by token success: %v", b)
		if b {
			sessionUser, _ = common.FromInterfaceToSessionUser(get)
			logger.Logger().Debugf("middleware to get user by token: %v", sessionUser)
		}
	}

	return context.WithValue(ctx, token, sessionUser)
}

func NewContextMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := newContextWithRequestID(req.Context(), req)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
