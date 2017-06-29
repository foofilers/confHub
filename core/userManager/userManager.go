package userManager

import (
	"github.com/foofilers/cfhd/models"
	"github.com/foofilers/cfhd/db/userDao"
)

func AddUser(user *models.User) error {
	return userDao.Persist(user)
}

func ListUsers(query string, page, count int32, order string) ([]models.User, error) {
	return userDao.List(query, page, count, order)
}

func DeleteById(id string) error{
	return userDao.DeleteById(id)
}