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