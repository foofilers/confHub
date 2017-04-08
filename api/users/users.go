package users

import (
	iris "gopkg.in/kataras/iris.v6"
	"github.com/foofilers/confHub/models"
	"github.com/Sirupsen/logrus"
	"github.com/foofilers/confHub/etcd"
	"golang.org/x/net/context"
	"github.com/foofilers/confHub/auth"
	"github.com/foofilers/confHub/api/utils"
)

func InitAPI(router *iris.Router, handlersFn ...iris.HandlerFunc) *iris.Router {
	logrus.Info("initializing /users api resources")
	usersParty := router.Party("/users", handlersFn...)
	usersParty.Post("/", addUser)
	usersParty.Get("/", listUsers)
	usersParty.Put("/:username", updateUser)
	usersParty.Delete("/:username", deleteUser)
	return usersParty
}

func createUser(etcdCl *etcd.EtcdClient, user *models.User) error {
	// adding user
	if _, err := etcdCl.Client.UserAdd(context.TODO(), user.Username, user.Password); err != nil {
		return err
	}

	// setting user roles
	for _, role := range user.Roles {
		if _, err := etcdCl.Client.UserGrantRole(context.TODO(), user.Username, role); err != nil {
			return err
		}
	}
	return nil
}

func addUser(ctx *iris.Context) {
	user := &models.User{}
	if err := ctx.ReadJSON(user); err != nil {
		logrus.Error(err)
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}
	logrus.Debugf("adding user: %+v", user)

	etcdCl, err := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	if utils.HandleError(ctx, err) {
		return
	}
	defer etcdCl.Client.Close()
	if utils.HandleError(ctx, createUser(etcdCl, user)) {
		return
	}
	ctx.SetStatusCode(iris.StatusCreated)

}

func updateUser(ctx *iris.Context) {
	username := ctx.Param("username")
	logrus.Infof("Updating user %s", username)

	userUpd := &models.User{}
	if err := ctx.ReadJSON(userUpd); err != nil {
		logrus.Error(err)
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}
	logrus.Debugf("updating user %s with data: %+v", username, userUpd)

	etcdCl, err := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	if utils.HandleError(ctx, err) {
		return
	}
	defer etcdCl.Client.Close()
	userResp, err := etcdCl.Client.UserGet(context.TODO(), username)
	if utils.HandleError(ctx, err) {
		return
	}

	if username != userUpd.Username {
		//renaming
		if utils.HandleError(ctx, createUser(etcdCl, userUpd)) {
			return
		}
		etcdCl.Client.UserDelete(context.TODO(), username)
	} else {
		if len(userUpd.Password) > 0 {
			if _, err := etcdCl.Client.UserChangePassword(context.TODO(), username, userUpd.Password); utils.HandleError(ctx, err) {
				return
			}
		}
		roleMap := make(map[string]bool)
		for _, role := range userUpd.Roles {
			roleMap[role] = true
		}
		//remove role
		currRoleMap := make(map[string]bool)
		for _, currRole := range userResp.Roles {
			currRoleMap[currRole] = true
			if _, ok := roleMap[currRole]; !ok {
				//revoke role
				logrus.Infof("Revoke role %s from user %s", currRole, username)
				if _, err := etcdCl.Client.UserRevokeRole(context.TODO(), username, currRole); utils.HandleError(ctx, err) {
					//TODO fatalError or warning??
					return
				}
			}
		}
		//adding role
		for role := range roleMap {
			if _, ok := currRoleMap[role]; !ok {
				//add grant
				logrus.Infof("Grant role %s from user %s", role, username)
				if _, err := etcdCl.Client.UserGrantRole(context.TODO(), username, role); utils.HandleError(ctx, err) {
					//TODO fatalError or warning??
					return
				}
			}
		}
	}
	ctx.SetStatusCode(iris.StatusNoContent)
}

func listUsers(ctx *iris.Context) {
	etcdCl, err := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	if utils.HandleError(ctx, err) {
		return
	}
	defer etcdCl.Client.Close()
	usersResp, err := etcdCl.Client.UserList(context.TODO())
	if utils.HandleError(ctx, err) {
		return
	}
	ctx.JSON(iris.StatusOK,usersResp.Users)
}

func deleteUser(ctx *iris.Context) {
	username := ctx.Param("username")
	logrus.Infof("Deleting user [%s]", username)

	etcdCl, err := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	if utils.HandleError(ctx, err) {
		return
	}
	defer etcdCl.Client.Close()
	if _, err := etcdCl.Client.UserDelete(context.TODO(), username); utils.HandleError(ctx, err) {
		logrus.Errorf("mmm")
		return
	}
	ctx.SetStatusCode(iris.StatusNoContent)
}
