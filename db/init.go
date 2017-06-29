package db

import (
	"gopkg.in/mgo.v2"
	_ "github.com/foofilers/cfhd/conf"
	"github.com/spf13/viper"
	"github.com/foofilers/cfhd/util"
	"github.com/sirupsen/logrus"
)

var Session *mgo.Session

const COLLECTION_USERS string = "users"
const COLLECTION_APPLICATIONS string = "applications"
const COLLECTION_APPLICATIONS_CAPPED string = "applications_capped"

func Init() {
	var err error
	Session, err = mgo.Dial(viper.GetString("db.host"))
	if err != nil {
		panic(err)
	}
	Session.SetMode(mgo.Monotonic, true)
	createUsers()
	createApplications()
}

func createUsers() {
	ses := Session.Clone()
	defer ses.Close()
	err := ses.DB(viper.GetString("db.name")).C(COLLECTION_USERS).EnsureIndex(mgo.Index{
		Name:"UsersPK",
		Key:[]string{"username"},
		Unique:true,
		DropDups:false,
	})
	if err != nil {
		panic(err)
	}
}

func createApplications() {
	ses := Session.Clone()
	defer ses.Close()
	if names, err := ses.DB(viper.GetString("db.name")).CollectionNames(); err == nil && util.Contains(names, COLLECTION_APPLICATIONS_CAPPED) {
		logrus.Debugf("Collection %s already created", COLLECTION_APPLICATIONS_CAPPED)
		return
	}
	err := ses.DB(viper.GetString("db.name")).C(COLLECTION_APPLICATIONS_CAPPED).Create(&mgo.CollectionInfo{
		Capped:true,
		MaxDocs:10000,
		MaxBytes:104857600, //100Mb
	})
	if err != nil {
		panic(err)
	}
}

func Close() {
	Session.Close()
}
