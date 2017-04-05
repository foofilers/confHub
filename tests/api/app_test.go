package api

import (
	"testing"
	"gopkg.in/resty.v0"
	"gopkg.in/kataras/iris.v6"
)

func TestAddApplication(t *testing.T) {
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).SetFormData(map[string]string{
		"name": "testApp",
	}).Post(ServerUrl + "/api/apps")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, 201)
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
	}).Put(ServerUrl + "/api/apps/nopresentApplication")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, iris.StatusNotFound)
}

