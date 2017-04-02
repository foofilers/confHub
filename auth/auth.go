package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/Sirupsen/logrus"

	"reflect"
	"encoding/json"
)

type LoggedUser struct {
	Username        string `json:"username"`
	CryptedPassword string `json:"crypted_password"`
	Roles           []string `json:"roles"`
}

func (this *LoggedUser) Valid() error {
	return nil
}

func FromClaims(claims jwt.MapClaims) (LoggedUser, error) {
	u := LoggedUser{}
	logrus.Debugf("parse claims %+v to LoggedUser %s", claims, reflect.TypeOf(claims))
	jsonClaims, err := json.Marshal(claims)
	if err != nil {
		return u, err
	}
	json.Unmarshal(jsonClaims, &u);
	return u, err
}

