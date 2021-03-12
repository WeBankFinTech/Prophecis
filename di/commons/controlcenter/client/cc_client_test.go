package client

import (
	"fmt"
	"testing"
)

func TestCcClient_AuthAccessCheck(t *testing.T) {
	CcAuthType := "MLSS-Auth-Type"
	CcAuthSSOTicket := "MLSS-Ticket"
	CcAuthUser := "MLSS-UserID"
	CcAuthPWD := "MLSS-Passwd"
	CcAuthAppID := "MLSS-APPID"
	CcAuthAppTS := "MLSS-APPTimestamp"
	CcAuthAppToken := "MLSS-APPSignature"

	ccCleint := GetCcClient("")

	sSOAuthOptions := make(map[string]string)
	sSOAuthOptions[CcAuthType] = "SSO"
	sSOAuthOptions[CcAuthSSOTicket] = "CASTGC=TGT-170952-EU4GRXKRhXqqv2xappncPRH7ogcXFPn62iFIW1QBQAwaOSThuH-sso.webank.com"
	sSOAuthOptions[CcAuthUser] = "kirkzhou"
	stateCode, token, isSA, _ := ccCleint.AuthAccessCheck(&sSOAuthOptions)

	uMAuthOptions := make(map[string]string)
	uMAuthOptions[CcAuthType] = "UM"
	uMAuthOptions[CcAuthUser] = "kirkzhou"
	uMAuthOptions[CcAuthPWD] = "QWJjZDEyMzQK"
	stateCode, token, isSA, _ = ccCleint.AuthAccessCheck(&uMAuthOptions)
	fmt.Println(token)

	sysAuthOptions := make(map[string]string)
	sysAuthOptions[CcAuthType] = "SYSTEM"
	sysAuthOptions[CcAuthAppID] = "xxx"
	sysAuthOptions[CcAuthAppTS] = "1543490005818"
	sysAuthOptions[CcAuthAppToken] = "bdGG99inSviq%2FVRk9LtW%2BKNm2vg%3D"
	sysAuthOptions[CcAuthUser] = "kirkzhou"
	stateCode, token, isSA, _ = ccCleint.AuthAccessCheck(&sysAuthOptions)
	fmt.Println(token)
	fmt.Println(isSA)
	fmt.Println(stateCode)
}

func TestCcClient_UserNamespaceCheck(t *testing.T) {
	ccCleint := GetCcClient("")
	ccCleint.UserNamespaceCheck("3670ed6d-18f6-4aac-ae01-fe807ac6ba42", "kirkzhou", "test")
}
