package application

import "github.com/foofilers/confHub/utils"

var AppAlreadyExistError = utils.NewConfHubError("Application %s already exist")
var AppNotFoundError = utils.NewConfHubError("Application %s not found")
var VersionNotFound = utils.NewConfHubError("Version %s not found in %s application")
var CurrentVersionNotSetted = utils.NewConfHubError("No current version setted for %s application")
var VersionAlreadyExistError = utils.NewConfHubError("Version %s already exist in %s application")
var ReferenceNotFoundError = utils.NewConfHubError("Reference %s not found for %s application")
var TooManyReferenceLinksError = utils.NewConfHubError("Too many reference link for %s application version %s")
var ValueNotFoundError = utils.NewConfHubError("Value %s not found for %s application version %s")
