package application

import "github.com/foofilers/confHub/utils"


var AppAlreadyExistError = utils.NewConfHubError("Application %s already Exist")
var AppNotFoundError = utils.NewConfHubError("Application %s not found")
