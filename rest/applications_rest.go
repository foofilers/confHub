package rest

import (
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/foofilers/cfhd/core/applicationManager"
	"net/http"
	"github.com/sirupsen/logrus"
)

func InitApplication(server *gin.Engine) {
	server.GET("/application/:name/:version", get)
}

func get(ctx *gin.Context) {
	appName := ctx.Param("name")
	appVersion := ctx.Param("version")
	app, err := applicationManager.GetApplication(appName, appVersion)
	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	logrus.Debug(app)
	if app == nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, app)
}