package secret

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Secret struct {
	ID bson.ObjectId `bson:"_id" json:"-"`
	Hash string `bson:"hash" json:"hash"`
	SecretText string `bson:"secretText" json:"secretText"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	ExpiresAt time.Time `bson:"expiresAt" json:"expiresAt"`
	RemainingViews int `bson:"remainingViews" json:"remainingViews"`
}