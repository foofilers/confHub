package models

import "time"

type ApplicationConfiguration map[string]interface{}

type Application struct {
	Id            string `bson:"_id,omitempty" json:"id" structs:",omitempty"`
	Name          string `bson:"name" json:"name" structs:"name,omitempty"`
	Version       string `bson:"version" json:"version" structs:"version,omitempty"`
	UpdatedAt     time.Time `bson:"updatedAt" json:"updatedAt" structs:"updatedAt,omitempty"`
	Configuration ApplicationConfiguration `bson:"configuration" json:"configuration" structs:"configuration,omitempty"`
}
