package api

import (
	"testing"
	"gopkg.in/resty.v0"
	"gopkg.in/kataras/iris.v6"
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
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).SetFormData(map[string]string{
		"name": "toDelApp",
	}).Post(ServerUrl + "/api/apps")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, 201)

	resp, err = resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).Delete(ServerUrl + "/api/apps/toDelApp")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, 204)
}

func TestDeleteNotPresentApplication(t *testing.T) {
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).Delete(ServerUrl + "/api/apps/notPresentApp")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, iris.StatusNotFound)
}


