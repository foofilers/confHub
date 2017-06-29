package userDao

import (
	"github.com/foofilers/cfhd/models"
	"github.com/spf13/viper"
	"github.com/foofilers/cfhd/db"
	"gopkg.in/mgo.v2/bson"
)

func Persist(user *models.User) error {
	sess := db.Session.Clone()
	defer sess.Close()
	user.Id = bson.NewObjectId().Hex()
	return sess.DB(viper.GetString("db.name")).C(db.COLLECTION_USERS).Insert(user)
}

func List(search string, page, count int32, order string) ([]models.User, error) {
	sess := db.Session.Clone()
	defer sess.Close()
	var result []models.User
	var srcParam bson.M;
	if len(search) > 0 {
		srcParam = bson.M{"username":"/" + search + "/"}
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

func DeleteById(id string) error {
	sess := db.Session.Clone()
	defer sess.Close()
	return sess.DB(viper.GetString("db.name")).C("users").RemoveId(id)
}
