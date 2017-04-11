package api

import (
	"testing"
	"gopkg.in/resty.v0"
	"gopkg.in/kataras/iris.v6"
)

func SetValue(t *testing.T, app, version, key, value string) {
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).SetFormData(map[string]string{
		"value": value,
	}).Put(ServerUrl + "/api/values/" + app + "/" + version + "/" + key)
	if err != nil {
		t.Fatal("setting value", err)
	}
	checkHttpStatus(t, resp, iris.StatusNoContent)
}

func GetValue(t *testing.T, app, version, key string, expStatus int) []byte {
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).Get(ServerUrl + "/api/values/" + app + "/" + version + "/" + key)
	if err != nil {
		t.Fatal("getting value", err)
	}
	checkHttpStatus(t, resp, expStatus)
	return resp.Body()
}

func DeleteValue(t *testing.T, app, version, key string) {
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).Delete(ServerUrl + "/api/values/" + app + "/" + version + "/" + key)
	if err != nil {
		t.Fatal("deleting value", err)
	}
	checkHttpStatus(t, resp, iris.StatusNoContent)
}

func TestSetGetDeleteValue(t *testing.T) {
	appName := "valueApp1"
	version := "1.0.0"
	CreateApp(t, appName)
	CreateVersion(t, appName, version)
	SetValue(t, appName, version, "prop1", "val1")
	prop1Value := string(GetValue(t, appName, version, "prop1", iris.StatusOK))
	if prop1Value != "val1" {
		t.Fatalf("the value [%s] doesn't match with [val1]", prop1Value)
	}
	DeleteValue(t, appName, version, "prop1")

	GetValue(t, appName, version, "prop1", iris.StatusNotFound)
}
