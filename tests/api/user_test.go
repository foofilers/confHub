package api

import (
	"testing"
	"gopkg.in/resty.v0"
	"github.com/foofilers/confHub/models"
	"gopkg.in/kataras/iris.v6"
	"encoding/json"
)

func GetUserList(t *testing.T) map[string]bool {
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).Get(ServerUrl + "/api/users")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, 200)
	users := make([]string, 0)
	if err := json.Unmarshal(resp.Body(), &users); err != nil {
		t.Fatal(err)
	}
	res := make(map[string]bool)
	for _, u := range users {
		res[u] = true
	}
	return res
}

func CreateUser(t *testing.T, username, password string, roles []string) {
	user := &models.User{Username:username, Password:password, Roles:roles}
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).SetBody(user).Post(ServerUrl + "/api/users")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, 201)
}

func UpdateUser(t *testing.T, currUsername,newUsername, password string, roles []string) {
	user := &models.User{Username:newUsername, Password:password, Roles:roles}
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).SetBody(user).Put(ServerUrl + "/api/users/"+currUsername)
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, 204)
}

func TestAddUser(t *testing.T) {
	CreateUser(t, "user1", "password1", nil)

	if _, ok := GetUserList(t)["user1"]; !ok {
		t.Fatal("User not found on user lists")
	}
	Login(t,"user1","password1")
}

func TestUserRole(t *testing.T) {
	appName := "testPermApp1"
	CreateApp(t, appName)
	CreateUser(t, "permUser1", "password1", []string{appName + "RW", appName + "R"})
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
	user.Roles = []string{appName2 + "RW", appName2 + "R"}
	resp, err = resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).SetBody(user).Put(ServerUrl + "/api/users/permUser2")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, iris.StatusNoContent)

}

func TestUserDelete(t *testing.T) {
	user := &models.User{Username:"toDelUser", Password:"password1"}
	resp, err := resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).SetBody(user).Post(ServerUrl + "/api/users")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, 201)
	if _, ok := GetUserList(t)["toDelUser"]; !ok {
		t.Fatal("User not found on user lists")
	}

	resp, err = resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).Delete(ServerUrl + "/api/users/toDelUser")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, 204)

	if _, ok := GetUserList(t)["toDelUser"]; ok {
		t.Fatal("User found on user lists after the deletion")
	}
	resp, err = resty.R().SetHeader("Authorization", "Bearer " + Login(t, "root", RootPwd)).Delete(ServerUrl + "/api/users/toDelUser")
	if err != nil {
		t.Fatal(err)
	}
	checkHttpStatus(t, resp, 404)

}
