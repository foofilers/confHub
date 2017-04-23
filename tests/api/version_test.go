package api

import (
	"testing"
	"gopkg.in/resty.v0"
	"encoding/json"
	"reflect"
	"gopkg.in/kataras/iris.v6"
	"github.com/foofilers/confHub/models"
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

func CopyVersion(t *testing.T, appName, srcVersion, dstVersion string, expStatus int) {
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).SetFormData(map[string]string{
		"version":dstVersion,
	}).Put(ServerUrl + "/api/versions/" + appName + "/" + srcVersion + "/copy")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, expStatus)
}

func GetVersion(t *testing.T, appName string, expStatus int) *models.ApplicationVersion {
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).Get(ServerUrl + "/api/versions/" + appName)
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, expStatus)
	appVersion := &models.ApplicationVersion{}
	if err := json.Unmarshal(resp.Body(), &appVersion); err != nil {
		t.Fatal(err)
	}
	return appVersion
}

func GetVersionList(t *testing.T, appName string, expStatus int) map[string]bool {
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).Get(ServerUrl + "/api/versions/" + appName)
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, expStatus)
	appVersion := &models.ApplicationVersion{}
	if err := json.Unmarshal(resp.Body(), &appVersion); err != nil {
		t.Fatal(err)
	}
	res := make(map[string]bool)
	for _, v := range appVersion.Versions {
		res[v] = true
	}
	return res
}

func DeleteVersion(t *testing.T, appName, appVersion string, codeExp int) {
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).Delete(ServerUrl + "/api/versions/" + appName + "/" + appVersion)
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, codeExp)
}

func TestAIOVersion(t *testing.T) {
	appName := RandString(8)
	CreateApp(t, appName)
	CreateVersion(t, appName, "1.0.0")
	CreateVersion(t, appName, "2.0.0")
	if _, found := GetVersionList(t, appName, 200)["1.0.0"]; !found {
		t.Fatal("Version 1.0.0 should be present")
	}
	DeleteVersion(t, appName, "2.0.0", 204)
	if _, found := GetVersionList(t, appName, 200)["2.0.0"]; found {
		t.Fatal("Version 2.0.0 should NOT be present")
	}
}

func TestSetDefaultVersion(t *testing.T) {
	appName := RandString(8)
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

func TestCopyVersion(t *testing.T) {
	appName := RandString(8)
	srcVer := "1.0.0"
	dstVer := "2.0.0"
	CreateApp(t, appName)
	CreateVersion(t, appName, srcVer)
	config := map[string]string{
		"prop1":"val1",
		"prop2":"val2",
	}
	PutConfig(t, appName, srcVer, config)
	CopyVersion(t, appName, srcVer, dstVer, iris.StatusOK)

	getConfig := GetConfig(t, appName, dstVer)
	if !reflect.DeepEqual(config, getConfig) {
		t.Fatalf("map %+v are not equals to %+v", config, getConfig)
	}
}

func TestCopyVersionAlreadyExist(t *testing.T) {
	appName := RandString(8)
	srcVer := "1.0.0"
	dstVer := "2.0.0"
	CreateApp(t, appName)
	CreateVersion(t, appName, srcVer)
	config := map[string]string{
		"prop1":"val1",
		"prop2":"val2",
	}
	PutConfig(t, appName, srcVer, config)
	CopyVersion(t, appName, srcVer, dstVer, iris.StatusOK)
	CopyVersion(t, appName, dstVer, srcVer, iris.StatusConflict)
}

func TestCopyVersionSrcNotFound(t *testing.T) {
	appName := RandString(8)
	srcVer := "1.0.0"
	CreateApp(t, appName)
	CopyVersion(t, appName, srcVer, "2.0.0", iris.StatusNotFound)
}

func TestFirstVersionIsCurrentVersion(t *testing.T) {
	appName := RandString(8)
	srcVer := "1.0.0"
	CreateApp(t, appName)
	CreateVersion(t, appName, srcVer)
	version := GetVersion(t, appName, 200)
	if version.DefaultVersion != srcVer {
		t.Fatalf("the default version should be %s but it's %s", srcVer, version.DefaultVersion)
	}
}

func TestCannotDeleteCurrentVersion(t *testing.T) {
	appName := RandString(8)
	srcVer := "1.0.0"
	CreateApp(t, appName)
	CreateVersion(t, appName, srcVer)
	DeleteVersion(t, appName, srcVer, iris.StatusPreconditionFailed)
}

