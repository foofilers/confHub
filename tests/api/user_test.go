package api

import (
	"testing"
	"gopkg.in/resty.v0"
	"github.com/foofilers/confHub/models"
	"gopkg.in/kataras/iris.v6"
)

func TestAddUser(t *testing.T) {
	user := &models.User{Username:"user1", Password:"password1"}
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).SetBody(user).Post(ServerUrl + "/api/users")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, 201)
}

func TestUserRole(t *testing.T) {
	appName := "testPermApp1"
	CreateApp(t, appName)

	user := &models.User{Username:"permUser1", Password:"password1", Roles:[]string{appName + "RW", appName + "R"}}
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).SetBody(user).Post(ServerUrl + "/api/users")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, 201)
}

func TestNotExistRole(t *testing.T) {
	user := &models.User{Username:"fakeRoleUser1", Password:"password1", Roles:[]string{"fakeRole"}}
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).SetBody(user).Post(ServerUrl + "/api/users")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, iris.StatusPreconditionFailed)
}


func TestUserUpdate(t *testing.T) {
	appName := "testPermApp2"
	CreateApp(t, appName)

	appName2 := "testPermApp3"
	CreateApp(t, appName2)

	user := &models.User{Username:"permUser2", Password:"password1", Roles:[]string{appName + "RW", appName + "R"}}
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).SetBody(user).Post(ServerUrl + "/api/users")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, 201)
	// update
	user.Roles=[]string{appName2 + "RW", appName2 + "R"}
	resp, err = resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).SetBody(user).Put(ServerUrl + "/api/users/permUser2")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, iris.StatusNoContent)

}
