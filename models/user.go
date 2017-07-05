package models

import "gopkg.in/mgo.v2/bson"

type Permission struct {
	Application string
	Perm        []string
}

type User struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	Username    string `bson:"username" json:"username"`
	Password    string `bson:"password" json:"password"`
	Email       string `bson:"email" json:"email"`
	Admin       bool `bson:"admin" json:"admin"`
	Permissions []Permission `bson:"permissions" json:"permissions"`
}
