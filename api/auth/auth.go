package auth

import (
	"gopkg.in/kataras/iris.v6"
	"github.com/foofilers/confHub/etcd"
	"github.com/foofilers/confHub/auth"
	"github.com/Sirupsen/logrus"
	"github.com/coreos/pkg/cryptoutil"
	cnf "github.com/spf13/viper"
	"encoding/base64"
	"golang.org/x/net/context"
	"github.com/dgrijalva/jwt-go"
)



func InitAuth(router *iris.Router) {
	logrus.Info("initializing /auth api resources")
	authParty := router.Party("/auth")
	authParty.Post("/login", login)
	authParty.Get("/logout", logout)
}

func login(ctx *iris.Context) {
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")

	cypheredPwd, err := cryptoutil.AESEncrypt([]byte(password), []byte(cnf.GetString("pwdSecretKey")))
	user := &auth.LoggedUser{Username:username, CryptedPassword:base64.StdEncoding.EncodeToString(cypheredPwd)}

	etcdCl, err := etcd.New(username, password)
	if err != nil {
		logrus.Error(err)
		ctx.EmitError(iris.StatusForbidden)
		return
	}
	// retrieve user information
	userInfo, err := etcdCl.Client.UserGet(context.TODO(), username)
	if err != nil {
		logrus.Error(err)
		ctx.EmitError(iris.StatusForbidden)
		return
	}
	user.Roles = userInfo.Roles
	// jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)
	signedToken, err := token.SignedString([]byte(cnf.GetString("jwtSecretKey")))
	if err != nil {
		logrus.Error(err)
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}
	ctx.WriteString(signedToken)
}

func logout(ctx *iris.Context) {
	ctx.Writef("logout")
}