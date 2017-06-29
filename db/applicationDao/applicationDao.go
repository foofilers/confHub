package applicationDao

import (
	"github.com/foofilers/cfhd/models"
	"github.com/foofilers/cfhd/db"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2/bson"
	"github.com/sirupsen/logrus"
	"github.com/fatih/structs"
	"gopkg.in/mgo.v2"
	"time"
)

func Persist(application *models.Application) error {
	sess := db.Session.Clone()
	defer sess.Close()
	application.Id = bson.NewObjectId().Hex()
	insErr := sess.DB(viper.GetString("db.name")).C(db.COLLECTION_APPLICATIONS).Insert(application)
	insCappErr := sess.DB(viper.GetString("db.name")).C(db.COLLECTION_APPLICATIONS_CAPPED).Insert(application)
	if insCappErr != nil {
		logrus.Error(insCappErr)
	}
	return insErr

}

func Get(name, version string) (*models.Application, error) {
	sess := db.Session.Clone()
	defer sess.Close()
	var result models.Application
	query := bson.M{
		"name":name,
	}
	if len(version) > 0 {
		query["version"] = version
	} else {
		query["latest"] = true
	}
	err := sess.DB(viper.GetString("db.name")).C(db.COLLECTION_APPLICATIONS).Find(query).One(&result)
	if err == mgo.ErrNotFound {
		return nil, nil
	}
	return &result, err
}

func List(search *models.Application, page, count int32, order string) ([]models.Application, error) {
	sess := db.Session.Clone()
	defer sess.Close()
	var result []models.Application
	srcParam := bson.M{}
	if search != nil {
		srcParam = structs.Map(search)
	}
	logrus.Debugf("applicationDao:list application query filter:%+v", srcParam)
	query := sess.DB(viper.GetString("db.name")).C(db.COLLECTION_APPLICATIONS).Find(srcParam)
	if len(order) > 0 {
		query.Sort(order)
	}
	if count > 0 && page > 0 {
		query.Skip(int(page * count)).Limit(int(count))
	}
	err := query.All(&result)
	logrus.Debugf("applicationDao:list application result:%+v", result)
	return result, err
}

func Watch(name string, appCh chan *models.Application, stopCh chan bool) error {
	sess := db.Session.Clone()
	defer sess.Close()
	query := bson.M{
		"name":name,
	}
	tail := sess.DB(viper.GetString("db.name")).C(db.COLLECTION_APPLICATIONS_CAPPED).Find(query).Tail(1 * time.Second)
	defer tail.Close()

	logrus.Debug("start tailing")
	running := true
	go func() {
		select {
		case <-stopCh:
			logrus.Debug("applicationDao: watch received stop signal")
			running = false
		}
	}()
	app := &models.Application{}
	for running {
		for tail.Next(app) {
			appCh <- app
			logrus.Debug(app)
			app = &models.Application{}
		}
		if tail.Err() != nil {
			logrus.Debugf("applicationDao:Tail Error %+v", tail.Err())
			return tail.Err()
		}
		if running && tail.Timeout() {
			logrus.Debug("applicationDao:Tail timeout")
			appCh <- nil
		}
	}
	logrus.Debug("applicationDao:watch finished")
	return nil
}