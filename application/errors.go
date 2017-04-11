package application

import "github.com/foofilers/confHub/utils"


var AppAlreadyExistError = utils.NewConfHubError("Application %s already exist")
var AppNotFoundError = utils.NewConfHubError("Application %s not found")
var VersionNotFound = utils.NewConfHubError("Version %s not found in %s application")
var CurrentVersionNotSetted = utils.NewConfHubError("No current version setted for %s application")
var VersionAlreadyExistError = utils.NewConfHubError("Version %s already exist in %s application")
