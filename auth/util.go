package auth

import "github.com/foofilers/cfhd/models"

func IsApplicationGranted(appName string, user *models.User) bool {
	if user.Admin {
		return true
	}
	for _, perm := range user.Permissions {
		if perm.Application == appName {
			return true
		}
	}
	return false
}
