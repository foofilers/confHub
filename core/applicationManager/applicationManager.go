package applicationManager

import (
	"github.com/foofilers/cfhd/models"
	"github.com/foofilers/cfhd/db/applicationDao"
)

func AddApplication(application *models.Application) error {
	return applicationDao.Persist(application)
}

func GetApplication(name, version string) (*models.Application, error) {
	return applicationDao.Get(name, version)
}

func ListApplication(search *models.Application, page, count int32, order string) ([]models.Application, error) {
	return applicationDao.List(search, page, count, order)
}

func WatchApplication(name string, applications chan *models.Application, stopCh chan bool) error {
	go func() {
		applicationDao.Watch(name, applications,stopCh)
	}()
	return nil
}

