package userManager

import (
	"github.com/foofilers/cfhd/models"
	"github.com/foofilers/cfhd/db/userDao"
	"crypto/md5"
	"fmt"
)

func AddUser(user *models.User, password string) error {
	user.Password = fmt.Sprintf("%x",md5.Sum([]byte(password)))
	return userDao.Persist(user)
}

func ListUsers(query string, page, count int32, order string) ([]models.User, error) {
	return userDao.List(query, page, count, order)
}

func DeleteById(id string) error {
	return userDao.DeleteById(id)
}

//todo cache it!
func GetById(userId string) (*models.User, error) {
	return userDao.GetById(userId)
}

func Login(username, password string) (*models.User, error) {
	md5Pwd := fmt.Sprintf("%x",md5.Sum([]byte(password)))
	return userDao.GetByUsernameAndPassword(username, md5Pwd)
}