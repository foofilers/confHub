package api

import (
	"gopkg.in/kataras/iris.v6"
	"github.com/foofilers/confHub/api/config"
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	cnf "github.com/spf13/viper"
	auth_api "github.com/foofilers/confHub/api/auth"
	"github.com/foofilers/confHub/api/users"
	"github.com/Sirupsen/logrus"
	"github.com/foofilers/confHub/auth"
	"github.com/foofilers/confHub/api/roles"
)

var jwtMiddleware *jwtmiddleware.Middleware

func loggedUserMiddleware(ctx *iris.Context) {
	token := jwtMiddleware.Get(ctx)
	if token == nil {
		logrus.Error("cannot find jwt field on context")
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}
	loggedUser, err := auth.FromClaims(token.Claims.(jwt.MapClaims))
	if err != nil {
		logrus.Error("error on loggedUser parser from jwt token:", err)
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}
	ctx.Set("LoggedUser", loggedUser)
	ctx.Next()
}

func InitApi(router *iris.Router) {
	jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(cnf.GetString("jwtSecretKey")), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	auth_api.InitAuth(router)
	users.InitAPI(router, jwtMiddleware.Serve, loggedUserMiddleware)
	config.InitAPI(router, jwtMiddleware.Serve, loggedUserMiddleware)
	roles.InitAPI(router, jwtMiddleware.Serve, loggedUserMiddleware)
}