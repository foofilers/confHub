package userDao

import (
	"github.com/foofilers/cfhd/models"
	"github.com/spf13/viper"
	"github.com/foofilers/cfhd/db"
	"gopkg.in/mgo.v2/bson"
	"github.com/sirupsen/logrus"
)

func Persist(user *models.User) error {
	sess := db.Session.Clone()
	defer sess.Close()
	user.Id = bson.NewObjectId()
	logrus.Debugf("Inserting user[%+v] on db ", user)
	return sess.DB(viper.GetString("db.name")).C(db.COLLECTION_USERS).Insert(user)
}

func List(search string, page, count int32, order string) ([]models.User, error) {
	sess := db.Session.Clone()
	defer sess.Close()
	var result []models.User
	var srcParam bson.M;
	//fixme move the search logic to userManager and change search parameters
	if len(search) > 0 {
		srcParam = bson.M{"username":bson.M{"$regex":".*" + search + ".*"}}
	}
	query := sess.DB(viper.GetString("db.name")).C(db.COLLECTION_USERS).Find(srcParam)
	if len(order) > 0 {
		query.Sort(order)
	}
	if count > 0 && page > 0 {
		query.Skip(int(page * count)).Limit(int(count))
	}
	return result, query.All(&result)
}

func GetById(id string) (*models.User, error) {
	sess := db.Session.Clone()
	defer sess.Close()
	var result models.User
	return &result, sess.DB(viper.GetString("db.name")).C(db.COLLECTION_USERS).FindId(bson.ObjectIdHex(id)).One(&result)
}

func DeleteById(id string) error {
	sess := db.Session.Clone()
	defer sess.Close()
	return sess.DB(viper.GetString("db.name")).C(db.COLLECTION_USERS).RemoveId(id)
}

func GetByUsernameAndPassword(username, md5password string) (*models.User, error) {
	sess := db.Session.Clone()
	defer sess.Close()
	query := bson.M{
		"username":username,
		"password":md5password,
	}
	logrus.Debugf("GetByUsernameAndPassword : %+v", query)
	var result models.User;
	return &result, sess.DB(viper.GetString("db.name")).C(db.COLLECTION_USERS).Find(query).One(&result)

}
