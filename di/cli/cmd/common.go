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

package cmd

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	oapiRuntime "github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"net/http"
	"time"
	dlaasClient "webank/DI/restapi/api_v1/client"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/terminal"
	log "github.com/sirupsen/logrus"
)

var (
	basicAuth oapiRuntime.ClientAuthInfoWriter
	//watsonUserInfo string
	authType string
	username string
	password string
	appid    string
	appts    string
	appToken string
)

const (
	//watsonUserInfoHeader = "X-Watson-Userinfo"
	CcAuthType       = "MLSS-Auth-Type"
	CcAuthSSOTicket  = "MLSS-Ticket"
	CcAuthUser       = "MLSS-UserID"
	CcAuthPWD        = "MLSS-Passwd"
	CcAuthAppID      = "MLSS-APPId"
	CcAuthAppTS      = "MLSS-APPTimestamp"
	CcAuthAppToken   = "MLSS-AppSignature"
	badUsernameOrPWD = "Bad username or password."
	defaultOpTimeout = 10 * time.Second
	dateFormat       = "2006-01-02 15:04:05.999999999 -0700 MST"
	cliConfigPath    = "/etc/climagic/climagic.json"
)

// DecorateRuntimeContext appends line, file and function context to the logger
func lflog() *log.Entry {
	if pc, file, line, ok := runtime.Caller(1); ok {
		fName := runtime.FuncForPC(pc).Name()
		clipDir := "dlaas-platform-cli/"
		trimmedFile := file[strings.LastIndex(file, clipDir)+len(clipDir):]
		trimmedFName := fName[strings.LastIndex(fName, clipDir)+len(clipDir):]
		return log.StandardLogger().WithField("file", trimmedFile).WithField("line", line).WithField("func", trimmedFName)
	}
	return log.NewEntry(log.StandardLogger())
}

// NewDlaaSClient is a helper for creating a new DLaaS REST API client with the
// right endpoint.
func NewDlaaSClient() (*dlaasClient.Dlaas, error) {
	envModeFlag := true
	if len(os.Getenv("DLAAS_URL")) <= 0 && len(os.Getenv("MLSS_AUTH_USER")) <= 0 && len(os.Getenv("MLSS_AUTH_PASSWD")) <= 0 && len(os.Getenv("MLSS_AUTH_TYPE")) <= 0 {
		envModeFlag = false
	}
	if envModeFlag {
		dlaasURL := os.Getenv("DLAAS_URL")
		if dlaasURL == "" {
			dlaasURL = dlaasClient.DefaultHost
		}
		lflog().Debugf("DLAAS_URL: %s", dlaasURL)

		u, _ := url.Parse(dlaasURL)

		lflog().Debugf("parsed DLAAS_URL: %+v", u)

		schemes := []string{u.Scheme}
		if u.Scheme == "" {
			schemes = []string{"http"}
		}

		if u.Path == "" {
			u.Path = dlaasClient.DefaultBasePath
		}

		var _username, _password, _authType string
		if u.User != nil {
			_username = u.User.Username()
			_password, _ = u.User.Password()
		} else {
			authType = os.Getenv("MLSS_AUTH_TYPE")
			username = os.Getenv("MLSS_AUTH_USER")
			password = os.Getenv("MLSS_AUTH_PASSWD")
			appid = os.Getenv("MLSS_APPID")
			appts = os.Getenv("MLSS_APPTimestamp")
			appToken = os.Getenv("MLSS_APPSignature")
			_authType = authType
			_username = username
			_password = password
			if _username == "" || _password == "" || _authType == "" {
				return nil, errors.New("Username or password or auth type not set")
			}
		}

		basicAuth = client.BasicAuth(username, password)
		lflog().Debugf("basicAuth: %+v", BasicAuth)

		transport := client.New(u.Host, u.Path, schemes)

		transport.Transport = createRoundTripper()
		return dlaasClient.New(transport, strfmt.Default), nil
	} else {
		configMap := map[string]interface{}{}
		configBytes, err := ioutil.ReadFile(cliConfigPath)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(configBytes, &configMap)
		if err != nil {
			return nil, err
		}
		dlaasURL := configMap["dlaas_url"].(string)
		if dlaasURL == "" {
			dlaasURL = dlaasClient.DefaultHost
		}
		lflog().Debugf("DLAAS_URL: %s", dlaasURL)
		u, _ := url.Parse(dlaasURL)
		lflog().Debugf("parsed DLAAS_URL: %+v", u)
		schemes := []string{u.Scheme}
		if u.Scheme == "" {
			schemes = []string{"http"}
		}
		if u.Path == "" {
			u.Path = dlaasClient.DefaultBasePath
		}
		var _username, _password, _authType string
		if u.User != nil {
			_username = u.User.Username()
			_password, _ = u.User.Password()
		} else {
			authType = configMap["mlss_auth_type"].(string)
			username = configMap["mlss_auth_user"].(string)
			password = configMap["mlss_auth_passwd"].(string)
			_authType = authType
			_username = username
			_password = password
			if _username == "" || _password == "" || _authType == "" {
				return nil, errors.New("Username or password or auth type not set")
			}
		}
		basicAuth = client.BasicAuth(username, password)
		lflog().Debugf("basicAuth: %+v", BasicAuth)
		transport := client.New(u.Host, u.Path, schemes)
		transport.Transport = createRoundTripper()
		return dlaasClient.New(transport, strfmt.Default), nil
	}
}

// Create roundTripper to inject X-Watson-Userinfo header into every request if not present
type roundTripper struct {
}

func (r *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// we set the dummy header if it is not set so we can use the CLI internally too (without going through Datapower)
	//lflog().Debugf("X-Watson-Userinfo: %s", req.Header.Get(watsonUserInfoHeader))
	if req.Header.Get(CcAuthType) == "" {
		lflog().Debugf("Adding %s: %s", CcAuthType, authType)
		req.Header.Add(CcAuthType, authType)
	}
	if req.Header.Get(CcAuthUser) == "" {
		lflog().Debugf("Adding %s: %s", CcAuthUser, username)
		req.Header.Add(CcAuthUser, username)
	}
	if req.Header.Get(CcAuthPWD) == "" {
		req.Header.Add(CcAuthPWD, password)
	}
	req.Header.Add(CcAuthAppID, appid)
	req.Header.Add(CcAuthAppTS, appts)
	req.Header.Add(CcAuthAppToken, appToken)
	return http.DefaultTransport.RoundTrip(req)
}

func createRoundTripper() http.RoundTripper {
	return &roundTripper{}
}

// BasicAuth returns the basic auth credentials.
func BasicAuth() oapiRuntime.ClientAuthInfoWriter {
	return basicAuth
}

// LocationToID return the ID component of a DLaaS Location header.
// For example:
//   /dlaas/api/v1/models/training-gQYXhh2gg -> training-gQYXhh2gg
func LocationToID(location string) string {
	// This implementation is intentionally restricted to only work with known
	// Location header values.
	prefixes := []string{"/di/v1/models/"}
	id := location
	for _, prefix := range prefixes {
		if strings.HasPrefix(location, prefix) {
			id = strings.TrimPrefix(location, prefix)
			break
		}
	}
	return id
}

// IsValidManifest returns true if the data is a valid manifest file.
func IsValidManifest(manifest []byte) bool {
	return true
}

// IsValidZip returns true of the path points to a valid zip file.
func IsValidZip(zipFile string) bool {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return false
	}
	defer r.Close()
	return true
}

func stringOrDefault(val string, def string) string {
	if val == "" {
		return def
	}
	return val
}

//zipit zips contents of source dir into target file
func zipit(source string) (*os.File, error) {
	// on windows we are seeing the source path in the ZIP, so we change the current
	// dir as a workaround and reset is back after
	if runtime.GOOS == "windows" {
		currentDir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		defer os.Chdir(currentDir)
		source, _ = filepath.Abs(source)
		if err := os.Chdir(source); err != nil {
			return nil, err
		}
	}
	zipfile, err := ioutil.TempFile("", "dlaas")
	if err != nil {
		return nil, err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil, err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if source != path {
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = strings.TrimPrefix(path, source)
			}

			if info.IsDir() {
				header.Name += "/"
			} else {
				header.Method = zip.Deflate
			}

			writer, err := archive.CreateHeader(header)

			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			return err
		}
		return err
	})

	return zipfile, nil
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func openManifestFile(ui terminal.UI, manifestFile string) os.File {
	f, err := os.Open(manifestFile)
	if err != nil {
		ui.Failed("Error opening manifest file %s: %v", manifestFile, err)
	}
	data, err := ioutil.ReadFile(manifestFile)
	if err != nil {
		ui.Failed("Error reading manifest file.")
	}
	if !IsValidManifest(data) {
		ui.Failed("Bad manifest file.")
	}
	return *f
}

func openModelDefinitionFile(ui terminal.UI, mdFile string) *os.File {
	f, err := os.Open(mdFile)
	if err != nil {
		ui.Failed("Error opening slug file %s: %v", mdFile, err)
	}
	if !IsValidZip(mdFile) {
		ui.Failed("Model definition file is not a valid zip file.")
	}
	return f
}

func responseError(s string, err error, ui terminal.UI) {
	if s != "" {
		ui.Failed(s)
	} else {
		if apiErr, ok := err.(*oapiRuntime.APIError); ok {
			ui.Failed(fmt.Sprintf("Error code: %d\nDetails: %s", apiErr.Code, apiErr.Response))
		}
	}

	// TODO we may want to strip everything after the base URL
	resp, _ := http.DefaultClient.Get(os.Getenv("DLAAS_URL") + "/health")
	if resp != nil {
		if resp.StatusCode == http.StatusServiceUnavailable {
			ui.Failed("Error 503: service unavailable")
		} else {
			ui.Failed("Error %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		}
	} else {
		ui.Failed("Error: env var 'DLAAS_URL' is invalid \n       %s", os.Getenv("DLAAS_URL"))
	}

	return
}

func formatTimestamp(dateTime string) string {
	// FIXME MLSS Change: bug fixed for time parse
	//t, err := time.Parse(dateFormat, dateTime)
	ts, err := strconv.ParseInt(dateTime, 10, 64)
	if err != nil {
		return "N/A"
	}
	t := time.Unix(ts/1000, 0)
	if err != nil {
		return "N/A"
	}
	return t.Format(time.RFC1123)
}
