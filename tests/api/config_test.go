package api

import (
	"testing"
	"gopkg.in/resty.v0"
	"encoding/json"
	"gopkg.in/kataras/iris.v6"
	"reflect"
	"fmt"
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

