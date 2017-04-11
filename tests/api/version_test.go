package api

import (
	"testing"
	"gopkg.in/resty.v0"
	"encoding/json"
)

func CreateVersion(t *testing.T, appName, appVersion string) {
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).SetFormData(map[string]string{
		"version": appVersion,
	}).Post(ServerUrl + "/api/versions/" + appName)
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, 201)
}

func SetDefaultVersion(t *testing.T, appName, appVersion string) {
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).Put(ServerUrl + "/api/versions/" + appName + "/" + appVersion)
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, 204)
}

func GetVersion(t *testing.T, appName string, expStatus int) map[string]bool {
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).Get(ServerUrl + "/api/versions/" + appName)
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, expStatus)
	versions := make([]string, 0)
	if err := json.Unmarshal(resp.Body(), &versions); err != nil {
		t.Fatal(err)
	}
	res := make(map[string]bool)
	for _, v := range versions {
		res[v] = true
	}
	return res
}

func DeleteVersion(t *testing.T, appName, appVersion string) {
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).Delete(ServerUrl + "/api/versions/" + appName + "/" + appVersion)
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, 204)
}

func TestAIOVersion(t *testing.T) {
	appName := "verApp1"
	CreateApp(t, appName)
	CreateVersion(t, appName, "1.0.0")
	if _, found := GetVersion(t, appName, 200)["1.0.0"]; !found {
		t.Fatal("Version 1.0.0 should be present")
	}
	DeleteVersion(t, appName, "1.0.0")
	if _, found := GetVersion(t, appName, 200)["1.0.0"]; found {
		t.Fatal("Version 1.0.0 should NOT be present")
	}
}

func TestSetDefaultVersion(t *testing.T) {
	appName := "currverapp1"
	const ver1 = "1.0.0"
	const ver2 = "2.0.0"
	const prop1V1Value = "val1"
	const prop1V2Value = "val2"
	CreateApp(t, appName)
	CreateVersion(t, appName, ver1)
	CreateVersion(t, appName, ver2)
	SetValue(t, appName, ver1, "prop1", prop1V1Value)
	SetValue(t, appName, ver2, "prop1", prop1V2Value)

	SetDefaultVersion(t, appName, ver1)
	confs := GetConfig(t, appName, "")
	if cnfProp1, found := confs["prop1"]; !found || cnfProp1 != prop1V1Value {
		t.Fatalf("wrong property value should be %s but is %s", prop1V1Value, cnfProp1)
	}

	SetDefaultVersion(t, appName, ver2)
	confs = GetConfig(t, appName, "")
	if cnfProp1, found := confs["prop1"]; !found || cnfProp1 != prop1V2Value {
		t.Fatalf("wrong property value should be %s but is %s", prop1V2Value, cnfProp1)
	}

}


