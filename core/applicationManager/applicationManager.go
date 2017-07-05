package applicationManager

import (
	"github.com/foofilers/cfhd/models"
	"github.com/foofilers/cfhd/db/applicationDao"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/fatih/set.v0"
	"regexp"
	"github.com/foofilers/cfhd/auth"
	"errors"
)

func AddApplication(application *models.Application) error {
	return applicationDao.Persist(application)
}

func GetApplication(name, version string) (*models.Application, error) {
	return applicationDao.Get(name, version)
}

func ListApplication(search *models.Application, page, count int32, order string, user *models.User) ([]models.Application, error) {
	srcParam := bson.M{}
	if !user.Admin {
		apps := set.New()
		for _, userPerm := range user.Permissions {
			if search != nil && len(search.Name) > 0 {
				if match, err := regexp.MatchString(".*" + search.Name + ".*", userPerm.Application); match && err != nil {
					apps.Add(userPerm.Application)
				}
			} else {
				apps.Add(userPerm.Application)
			}
		}
		srcParam["name"] = bson.M{
			"$in":set.StringSlice(apps),
		}
	} else {
		if search != nil && len(search.Name) > 0 {
			srcParam["name"] = bson.M{"$regex":".*" + search.Name + ".*"}
		}
	}
	return applicationDao.List(srcParam, page, count, order)
}

func WatchApplication(name string, applications chan *models.Application, stopCh chan bool, user *models.User) error {
	if !auth.IsApplicationGranted(name, user) {
		return errors.New("Permission denied")
	}
	go func() {
		applicationDao.Watch(name, applications, stopCh)
	}()
	return nil
}

