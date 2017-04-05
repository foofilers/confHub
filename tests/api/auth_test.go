package api

import (
	"testing"
	"gopkg.in/resty.v0"
	"gopkg.in/kataras/iris.v6"
)

func Login(t *testing.T, username, password string) string {
	resp, err := resty.R().SetFormData(map[string]string{
		"username":username,
		"password":password,
	}).Post(ServerUrl + "/api/auth/login")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t,resp,200)
	return string(resp.Body())
}

func TestLogin(t *testing.T) {
	Login(t, "root", RootPwd)
}

func TestUnauthorized(t *testing.T) {
	resp, err := resty.R().SetFormData(map[string]string{
		"username":"fakeUser",
		"password":"fakePassword",
	}).Post(ServerUrl + "/api/auth/login")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t,resp,iris.StatusForbidden)
}


