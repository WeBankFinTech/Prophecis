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
	"context"
	"fmt"
	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddyhttp/httpserver"
	"io/ioutil"
	"mlss-controlcenter-go/pkg/apigateway/caddy-plugin/internal"
	"mlss-controlcenter-go/pkg/apigateway/caddy-plugin/util"
	"mlss-controlcenter-go/pkg/logger"
	"net/http"
	"regexp"
	"strings"
)

const (
	AUTH_RSP_DEPTCODE     = "MLSS-DeptCode"
	AUTH_RSP_USERID       = "MLSS-UserID"
	AUTH_RSP_ORGCODE      = "MLSS-OrgCode"
	AUTH_RSP_LOGINID      = "MLSS-LoginID"
	AUTH_RSP_SUPERADMIN   = "MLSS-Superadmin"
	AUTH_RSP_TOKEN        = "MLSS-Token"
	AUTH_HEADER_APIKEY    = "MLSS-APPID"
	AUTH_HEADER_TIMESTAMP = "MLSS-APPTimestamp"
	AUTH_HEADER_AUTH_TYPE = "MLSS-Auth-Type"
	AUTH_HEADER_SIGNATURE = "MLSS-APPSignature"
	AUTH_HEADER_TICKET    = "MLSS-Ticket"
	AUTH_HEADER_UIURL     = "MLSS-UIURL"
	AUTH_HEADER_TOKEN     = "MLSS-Token"
	AUTH_HEADER_USERID    = "MLSS-UserID"
	AUTH_HEADER_COOKIE    = "MLSS-Cookie-key"
	AUTH_REAL_PATH        = "MLSS-RealPath"
	AUTH_REAL_METHOD      = "MLSS-RealMethod"
	AUTH_HEADER_PWD       = "MLSS-Passwd"
	AUTH_HEADER_REAL_IP   = "MLSS-RealIP"
	X_REAL_IP             = "X-Real-IP"
	SAMPLE_REQUEST        = "/cc/v1/sample"
	SAMPLE_REQUEST_USER   = "MLSS-Sample-User"
)

func init() {
	caddy.RegisterPlugin("auth_request", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	fmt.Println("auth_request!!!!!!!")

	rule, err := parse(c)
	if err != nil {
		return err
	}

	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return &AuthRequest{Next: next, Rule: rule}
	})

	return nil
}

func parse(c *caddy.Controller) (*Rule, error) {

	rule := &Rule{}
	rule.Checkers = make([]*Check, 0)

	if c.Next() {
		fmt.Printf("first Next: %v\n", c.Val())
		args := c.RemainingArgs()

		var check *Check

		switch len(args) {
		case 0:
			for c.Next() {
				fmt.Printf("Next: %v\n", c.Val())
				switch c.Val() {
				case "check":
					check = &Check{}
					rule.Checkers = append(rule.Checkers, check)
					check.ExclusionRules = make([]internal.ExclusionRule, 0)

					for c.NextBlock() {
						fmt.Printf("NextBlock: %v\n", c.Val())
						switch c.Val() {
						case "api":
							for c.NextArg() {
								fmt.Printf(c.Val())
								check.Api = c.Val()
							}
							break
						case "path":
							for c.NextArg() {
								fmt.Printf(c.Val())
								check.Path = c.Val()
							}
							break
						case "except":
							if !c.NextArg() {
								return nil, c.ArgErr()
							}

							method := c.Val()

							if !HasString(internal.HttpMethods, method) {
								return nil, c.ArgErr()
							}

							for c.NextArg() {
								path := c.Val()
								check.ExclusionRules = append(check.ExclusionRules, internal.ExclusionRule{Method: method, Path: path})
							}

							break
						}
					}

					break
				}
			}
			return rule, nil
		default:
			return rule, c.ArgErr()
		}
	}

	if c.Next() {
		return rule, c.ArgErr()
	}

	return rule, nil
}

func HasString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

type AuthRequest struct {
	Rule    *Rule
	Context context.Context
	Next    httpserver.Handler
}

type Rule struct {
	Checkers []*Check
}

type Check struct {
	Api            string
	Path           string
	ExclusionRules []internal.ExclusionRule
}

func (c AuthRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	realPath := r.URL.Path
	realMethod := r.Method

	logger.Logger().Debugf("AuthRequest ServeHTTP, realPath: %v, realMethod: %v", realPath, realMethod)

I:
	for _, check := range c.Rule.Checkers {
		pathAuthUrl := check.Path

		PathRegx := strings.ReplaceAll(pathAuthUrl, "*", "(.*)")
		pathMatched, _ := regexp.MatchString(PathRegx, realPath)
		if pathMatched {
			for _, rule := range check.ExclusionRules {
				exclusionRegx := strings.ReplaceAll(strings.TrimSpace(rule.Path), "*", "(.*)")
				matched, _ := regexp.MatchString(exclusionRegx, realPath)
				if matched && (rule.Method == internal.AllMethod || rule.Method == realMethod) {
					logger.Logger().Debugf("AuthRequest ServeHTTP match ExclusionRules, now jump out.exclusionRegx: %v, rule.Method: %v", exclusionRegx, rule.Method)

					continue I
				}
			}
			permitted, err := permissionValidateForUserAndAuth(r, w)
			if !permitted {
				return handleForbidden(w, err), nil
			}

		}
	}

	return c.Next.ServeHTTP(w, r)
}

func handleForbidden(w http.ResponseWriter, err error) int {
	return http.StatusUnauthorized
}

func permissionValidateForUserAndAuth(r *http.Request, w http.ResponseWriter) (bool, error) {
	var userInterceptor = util.GetAppConfig().AuthAddress.User

	client := &http.Client{}
	req, err := http.NewRequest("GET", userInterceptor, nil)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	realPath := r.URL.Path
	realMethod := r.Method
	remoteAddr := r.RemoteAddr
	remoteAddr = remoteAddr[0:strings.Index(remoteAddr, ":")]

	logger.Logger().Debugf("auth_request  r.header: %v", r.Header)

	addHeader(req, r, realPath, realMethod, remoteAddr)

	if SAMPLE_REQUEST == realPath {
		logger.Logger().Debugf("auth_request cookie r.header: %v, req headers: %v", r.Header.Get(AUTH_HEADER_TOKEN), req.Header)
		req.Header.Set(AUTH_HEADER_TOKEN, "")
		r.Header.Set(AUTH_HEADER_TOKEN, "")
		logger.Logger().Debugf("auth_request cookie r.header: %v, req headers: %v", r.Header.Get(AUTH_HEADER_TOKEN), req.Header)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	headers := resp.Header
	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if nil != readErr {

	}
	bodyMsg := fmt.Sprintf("auth_request realPath: %v user interceptor res body: %v", realPath, string(body))
	fmt.Println(bodyMsg)
	fmt.Println(resp.StatusCode)
	if resp.StatusCode != 200 || "" == headers.Get(AUTH_HEADER_TOKEN) {
		fmt.Println("user failed")
		return false, nil
	}

	var authInterceptor = util.GetAppConfig().AuthAddress.Auth

	reqAuth, err := http.NewRequest("GET", authInterceptor, nil)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	reqAuth.Header.Add(AUTH_HEADER_TOKEN, headers.Get(AUTH_HEADER_TOKEN))
	if SAMPLE_REQUEST == realPath {
		r.Header.Add(SAMPLE_REQUEST_USER, headers.Get(AUTH_HEADER_USERID))
		r.Header.Set(AUTH_HEADER_TOKEN, headers.Get(AUTH_HEADER_TOKEN))
		logger.Logger().Debugf("auth_request cookie r token after auth: %v", r.Header.Get(AUTH_HEADER_TOKEN))
	}
	reqAuth.Header.Add(AUTH_REAL_PATH, realPath)
	reqAuth.Header.Add(AUTH_REAL_METHOD, realMethod)
	respAuth, err := client.Do(reqAuth)

	auBody, auReadErr := ioutil.ReadAll(respAuth.Body)
	if nil != auReadErr {
		fmt.Println("auReadErr")
	}

	authBodyMsg := fmt.Sprintf("auth_request realPath: %v auth interceptor res body: %v", realPath, string(auBody))
	fmt.Println(authBodyMsg)

	fmt.Println(respAuth.StatusCode)
	if respAuth.StatusCode != 200 || "" == respAuth.Header.Get(AUTH_HEADER_TOKEN) {
		fmt.Println("auth failed")
		return false, nil
	}

	for k, v := range headers {
		headerMsg := fmt.Sprintf("response headers key: %v and value: %v", k, v)
		fmt.Println(headerMsg)
		w.Header().Add(k, v[0])
	}

	return true, nil
}

func addHeader(req *http.Request, r *http.Request, realPath string, realMethod string, remoteAddr string) {
	req.Header.Add(AUTH_HEADER_USERID, r.Header.Get(AUTH_HEADER_USERID))
	req.Header.Add(AUTH_HEADER_TOKEN, r.Header.Get(AUTH_HEADER_TOKEN))
	req.Header.Add(AUTH_HEADER_PWD, r.Header.Get(AUTH_HEADER_PWD))
	req.Header.Add(AUTH_HEADER_AUTH_TYPE, r.Header.Get(AUTH_HEADER_AUTH_TYPE))
	req.Header.Add(AUTH_HEADER_TICKET, r.Header.Get(AUTH_HEADER_TICKET))
	req.Header.Add(AUTH_HEADER_UIURL, r.Header.Get(AUTH_HEADER_UIURL))
	req.Header.Add(AUTH_REAL_PATH, realPath)
	req.Header.Add(AUTH_REAL_METHOD, realMethod)
	req.Header.Add(AUTH_HEADER_REAL_IP, remoteAddr)
	req.Header.Add(AUTH_HEADER_APIKEY, r.Header.Get(AUTH_HEADER_APIKEY))
	req.Header.Add(AUTH_HEADER_TIMESTAMP, r.Header.Get(AUTH_HEADER_TIMESTAMP))
	req.Header.Add(AUTH_HEADER_SIGNATURE, r.Header.Get(AUTH_HEADER_SIGNATURE))
	req.Header.Add("Content-Type", r.Header.Get("Content-Type"))
}
