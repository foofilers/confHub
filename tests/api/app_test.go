package api

import (
	"testing"
	"gopkg.in/resty.v0"
	"gopkg.in/kataras/iris.v6"
	"encoding/json"
)

func CreateApp(t *testing.T, name string) {
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).SetFormData(map[string]string{
		"name": name,
	}).Post(ServerUrl + "/api/apps")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, 201)
}

func TestAddApplication(t *testing.T) {
	CreateApp(t, "testApp")
}

func TestAddApplicationAlreadyExist(t *testing.T) {
	for _, expStatus := range ([]int{iris.StatusCreated, iris.StatusConflict}) {
		resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).SetFormData(map[string]string{
			"name": "duplApp",
		}).Post(ServerUrl + "/api/apps")
		if err != nil {
			t.Fatal(err)
		}
		checkHttpStatus(t, resp, expStatus)
	}
}

func TestRenameApplication(t *testing.T) {
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).SetFormData(map[string]string{
		"name": "origApp",
	}).Post(ServerUrl + "/api/apps")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, 201)

	resp, err = resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).SetFormData(map[string]string{
		"name": "newName",
	}).Put(ServerUrl + "/api/apps/origApp")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, 204)

	apps := GetListApplications(t)
	appFound := false
	for _, app := range apps {
		if app["name"] == "newName" {
			appFound = true
		}
	}
	if !appFound {
		t.Fatal("app not found")
	}

}

func TestRenameAlreadyExistApplication(t *testing.T) {
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).SetFormData(map[string]string{
		"name": "origApp2",
	}).Post(ServerUrl + "/api/apps")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, 201)

	resp, err = resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).SetFormData(map[string]string{
		"name": "origApp2",
	}).Put(ServerUrl + "/api/apps/origApp2")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, iris.StatusConflict)
}

func TestRenameNotPresentApplication(t *testing.T) {
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).SetFormData(map[string]string{
		"name": "newFakename",
	}).Put(ServerUrl + "/api/apps/notPresentApp")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, iris.StatusNotFound)
}

func TestDeleteApplication(t *testing.T) {
	CreateApp(t,"toDelApp")

	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).Delete(ServerUrl + "/api/apps/toDelApp")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, 204)

	apps := GetListApplications(t)
	appFound := false
	for _, app := range apps {
		if app["Name"] == "toDelApp" {
			appFound = true
		}
	}
	if appFound {
		t.Fatalf("apps should not exist")
	}
}

func TestDeleteNotPresentApplication(t *testing.T) {
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).Delete(ServerUrl + "/api/apps/notPresentApp")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, iris.StatusNotFound)
}

func GetListApplications(t *testing.T) []map[string]interface{} {
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).Get(ServerUrl + "/api/apps")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, iris.StatusOK)
	apps := make([]map[string]interface{}, 0)
	if err := json.Unmarshal(resp.Body(), &apps); err != nil {
		t.Fatal(err)
	}
	return apps
}

func TestListApplications(t *testing.T) {
	CreateApp(t, "app1")
	CreateApp(t, "app2")

	apps := GetListApplications(t)
	app1Found := false;
	app2Found := false
	for _, app := range apps {
		if app["name"] == "app1" {
			app1Found = true
		}
		if app["name"] == "app2" {
			app2Found = true
		}
	}
	if !app1Found || !app2Found {
		t.Fatalf("apps not found apps:%+v", apps)
	}
}


