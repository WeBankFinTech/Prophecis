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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"sort"
	"webank/AIDE/notebook-server/pkg/commons/config"
	cc "webank/AIDE/notebook-server/pkg/commons/controlcenter/client"
	"webank/AIDE/notebook-server/pkg/commons/logger"
)

const (
	// UserIDHeader is the name of the HTTP header used to identify the user
	UserIDHeader     = "X-DLaaS-UserID"
	CcAuthType       = "MLSS-Auth-Type"
	CcAuthSSOTicket  = "MLSS-Ticket"
	CcAuthUser       = "MLSS-UserID"
	CcAuthPWD        = "MLSS-Passwd"
	CcAuthAppID      = "MLSS-APPId"
	CcAuthAppTS      = "MLSS-APPTimestamp"
	CcAuthAppToken   = "MLSS-AppSignature"
	CcAuthSuperadmin = "MLSS-Superadmin"
)

// AuthOptions for the auth middleware.
type AuthOptions struct {
	ExcludedURLs []string
}

// NewAuthMiddleware creates a new http.Handler that adds authentication logic to a given Handler
func NewAuthMiddleware(opts *AuthOptions) func(h http.Handler) http.Handler {
	if opts == nil {
		opts = &AuthOptions{ExcludedURLs: []string{}}
	}

	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			logger.Logger().Debugf("Enter into auth handler")
			logger.Logger().Debugf("request: %+v", r)

			// check if the request URI matched excluded request URIs
			sort.Strings(opts.ExcludedURLs)
			i := sort.SearchStrings(opts.ExcludedURLs, r.RequestURI)
			if i < len(opts.ExcludedURLs) && opts.ExcludedURLs[i] == r.RequestURI {
				h.ServeHTTP(w, r)
				return // we are done
			}

			if r.Method == "OPTIONS" {
				if r.Header.Get("Access-Control-Request-Method") != "" {
					// TODO: Needs product stake-holder review

					logger.Logger().Debugf("cors preflight detected")

					// cors preflight request/response
					w.Header().Add("Access-Control-Allow-Origin", "*")
					w.Header().Add("Access-Control-Allow-Methods", "PUT, GET, POST, DELETE, OPTIONS")
					w.Header().Add("Access-Control-Allow-Headers", "origin, x-requested-with, content-type, authorization, x-watson-userinfo, x-watson-authorization-token")
					w.Header().Add("Access-Control-Max-Age", "86400")

					w.Header().Add("Content-Type", "text/html; charset=utf-8")

					w.WriteHeader(200)

					if flusher, ok := w.(http.Flusher); ok {
						flusher.Flush()
					}

					return
				}
			}

			logger.Logger().Debugf("Writing to header in callBefore \"Access-Control-Allow-Origin: *\"")

			w.Header().Add("Access-Control-Allow-Origin", "*")

			userID := r.Header.Get(CcAuthUser)
			if userID == "" {
				w.WriteHeader(401)
				w.Header().Add("Content-Type", "application/json; charset=utf-8")
				w.Write([]byte("{ \"message\" : \"Missing or malformed MLSS-UserID header.\"}"))
				return
			}

			authType := r.Header.Get(CcAuthType)
			ccToken := ""
			ccIsSA := ""
			authOptions := make(map[string]string)
			if authType == "" {
				token := r.Header.Get(cc.CcAuthToken)
				if token == "" {
					w.WriteHeader(403)
					w.Header().Add("Content-Type", "application/json; charset=utf-8")
					w.Write([]byte("{ \"message\" : \"Missing or malformed MLSS-Auth-Type or MLSS-Token header.\"}"))
					return
				}
				ccToken = token
				authOptions[cc.CcAuthToken] = ccToken
			} else {
				if authType == "SSO" {
					authOptions[CcAuthType] = "SSO"
					ssoTicket := r.Header.Get(CcAuthSSOTicket)
					if ssoTicket == "" {
						w.WriteHeader(403)
						w.Header().Add("Content-Type", "application/json; charset=utf-8")
						w.Write([]byte("{ \"message\" : \"Missing or malformed MLSS-Ticket header.\"}"))
						return
					}
					authOptions[CcAuthSSOTicket] = ssoTicket
				} else if authType == "UM" {
					authOptions[CcAuthType] = "UM"
					passWd := r.Header.Get(CcAuthPWD)
					if passWd == "" {
						w.WriteHeader(403)
						w.Header().Add("Content-Type", "application/json; charset=utf-8")
						w.Write([]byte("{ \"message\" : \"Missing or malformed MLSS-Passwd header.\"}"))
						return
					}
					authOptions[CcAuthPWD] = passWd
				} else if authType == "SYSTEM" {
					authOptions[CcAuthType] = "SYSTEM"
					appKey := r.Header.Get(CcAuthAppID)
					if appKey == "" {
						w.WriteHeader(403)
						w.Header().Add("Content-Type", "application/json; charset=utf-8")
						w.Write([]byte("{ \"message\" : \"Missing or malformed MLSS-AppID header.\"}"))
						return
					}
					authOptions[CcAuthAppID] = appKey
					appTimestamp := r.Header.Get(CcAuthAppTS)
					if appTimestamp == "" {
						w.WriteHeader(403)
						w.Header().Add("Content-Type", "application/json; charset=utf-8")
						w.Write([]byte("{ \"message\" : \"Missing or malformed MLSS-AppTimestamp header.\"}"))
						return
					}
					authOptions[CcAuthAppTS] = appTimestamp
					appToken := r.Header.Get(CcAuthAppToken)
					if appToken == "" {
						w.WriteHeader(403)
						w.Header().Add("Content-Type", "application/json; charset=utf-8")
						w.Write([]byte("{ \"message\" : \"Missing or malformed MLSS-AppToken header.\"}"))
						return
					}
					authOptions[CcAuthAppToken] = appToken
				}

			}

			authOptions[CcAuthUser] = userID
			// TODO: Get cc url from config.
			ccCleint := cc.GetCcClient(viper.GetString(config.CCAddress))
			stateCode, token, isSA, err := ccCleint.AuthAccessCheck(&authOptions)
			if err != nil || token == "" {
				w.WriteHeader(stateCode)
				w.Header().Add("Content-Type", "application/json; charset=utf-8")
				w.Write([]byte("{ \"message\" : \"Authorization check failed.\"}"))
				return
			}
			ccToken = token
			ccIsSA = isSA

			// set the header
			r.Header.Set(UserIDHeader, userID)
			r.Header.Set(cc.CcAuthToken, ccToken)
			r.Header.Set(cc.CcSuperadmin, ccIsSA)
			logger.Logger().Debugf("%s: %v, %v", UserIDHeader, userID, ccIsSA)

			if logger.Logger().GetLevel() == log.DebugLevel {
				entry := log.NewEntry(log.StandardLogger())
				for k, v := range r.Header {
					entry = entry.WithField(k, v)
				}
				logger.Logger().Debug("Request headers:")
			}

			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
