package models

type Permission struct {
	Application string
	Perm        []string
}

type User struct {
	Id          string `bson:"_id,omitempty"`
	Username    string
	Email       string
	Permissions []Permission
}
