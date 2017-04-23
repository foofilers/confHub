package api

import (
	"testing"
	"gopkg.in/resty.v0"
	"encoding/json"
	"gopkg.in/kataras/iris.v6"
	"reflect"
	"fmt"
	"github.com/Sirupsen/logrus"
)

func GetConfig(t *testing.T, appName, appVersion string) map[string]string {
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).Get(ServerUrl + "/api/configs/" + appName + "/" + appVersion)
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, iris.StatusOK)
	fmt.Printf("%+v", string(resp.Body()))
	configs := make(map[string]string)
	if err := json.Unmarshal(resp.Body(), &configs); err != nil {
		t.Fatal(err)
	}
	return configs
}

func PutConfig(t *testing.T, appName, appVersion string, configs map[string]string) {
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).SetBody(configs).Put(ServerUrl + "/api/configs/" + appName + "/" + appVersion)
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, iris.StatusNoContent)
}

func TestGetConfig(t *testing.T) {
	appName := "ConfigApp1"
	appVersion := "1.0.0"
	CreateApp(t, appName)
	CreateVersion(t, appName, appVersion)
	config := map[string]string{
		"prop1":"val1",
		"prop2":"val2",
	}
	PutConfig(t, appName, appVersion, config)
	getConfig := GetConfig(t, appName, appVersion)
	if !reflect.DeepEqual(config, getConfig) {
		t.Fatalf("map %+v are not equals to %+v", config, getConfig)
	}
}

func TestReference(t *testing.T) {
	app1Name := RandString(8)
	app1Version := "1.0.0"
	app2Name := RandString(8)
	app2Version := "2.0.0"

	prop1Key := "prop1Key"
	prop1Value := "prop1Value"

	refProp1Key := "refProp1Key"
	refProp1Value := fmt.Sprintf("${%s/%s/%s}", app1Name, app1Version, prop1Key)

	CreateApp(t, app1Name)
	CreateApp(t, app2Name)
	CreateVersion(t, app1Name, app1Version)
	CreateVersion(t, app2Name, app2Version)
	SetValue(t, app1Name, app1Version, prop1Key, prop1Value)
	SetValue(t, app2Name, app2Version, refProp1Key, refProp1Value)

	cnf := GetConfig(t, app2Name, app2Version)
	if cnf[refProp1Key] != prop1Value {
		t.Fatalf("The value should be %s but was %s", prop1Value, cnf[refProp1Key])
	}

	value := GetValue(t, app2Name, app2Version, refProp1Key, 200)
	if value != refProp1Value {
		t.Fatalf("the value should be %s but was %s", refProp1Value, value)
	}

	value = GetValueFollowingReference(t, app2Name, app2Version, refProp1Key, 200)
	if value != prop1Value {
		t.Fatalf("the value should be %s but was %s", prop1Value, value)
	}

}

func TestLoopReference(t *testing.T) {
	app1Name := RandString(8)
	app1Version := "1.0.0"
	app2Name := RandString(8)
	app2Version := "2.0.0"

	prop1Key := "prop1Key"
	refProp1Key := "refProp1Key"

	prop1Value := fmt.Sprintf("${%s.%s.%s}", app2Name, app2Version, refProp1Key)
	refProp1Value := fmt.Sprintf("${%s.%s.%s}", app1Name, app1Version, prop1Key)

	CreateApp(t, app1Name)
	CreateApp(t, app2Name)
	CreateVersion(t, app1Name, app1Version)
	CreateVersion(t, app2Name, app2Version)
	SetValue(t, app1Name, app1Version, prop1Key, prop1Value)
	SetValue(t, app2Name, app2Version, refProp1Key, refProp1Value)

	GetValue(t, app2Name, app2Version, refProp1Key, iris.StatusOK)
	GetValueFollowingReference(t, app2Name, app2Version, refProp1Key, iris.StatusPreconditionFailed)

}

func TestUnauthorizedReference(t *testing.T) {
	app1Name := RandString(8)
	app1Version := "1.0.0"
	app2Name := RandString(8)
	app2Version := "2.0.0"

	prop1Key := "prop1Key"
	prop1Value := "prop1Value"

	refProp1Key := "refProp1Key"
	refProp1Value := fmt.Sprintf("${%s/%s/%s}", app1Name, app1Version, prop1Key)

	CreateApp(t, app1Name)
	CreateApp(t, app2Name)

	CreateVersion(t, app1Name, app1Version)
	CreateVersion(t, app2Name, app2Version)
	SetValue(t, app1Name, app1Version, prop1Key, prop1Value)
	SetValue(t, app2Name, app2Version, refProp1Key, refProp1Value)

	CreateUser(t, "app1User", "app1Pwd", []string{app1Name + "RW"})
	CreateUser(t, "app2User", "app2Pwd", []string{app2Name + "RW"})
	logrus.Debug("getting values of app2")
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "app2User", "app2Pwd")).SetQueryParam("reference", "true").Get(ServerUrl + "/api/values/" + app2Name + "/" + app2Version + "/" + refProp1Key)
	if err != nil {
		t.Fatal("getting value", err)
	}
	checkHttpStatus(t, resp, iris.StatusForbidden)

	UpdateUser(t, "app2User","app2User", "app2Pwd", []string{app2Name + "RW",app1Name+"R"})

	resp, err = resty.R().SetHeader("Authorization", "Bearer " + Login(t, "app2User", "app2Pwd")).SetQueryParam("reference", "true").Get(ServerUrl + "/api/values/" + app2Name + "/" + app2Version + "/" + refProp1Key)
	if err != nil {
		t.Fatal("getting value", err)
	}
	checkHttpStatus(t, resp, 200)

	value := string(resp.Body())
	if value != prop1Value {
		t.Fatalf("the value shoult be %s but was %s", prop1Value, value)
	}
}



