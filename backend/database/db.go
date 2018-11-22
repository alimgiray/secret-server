package database

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	. "../secret"
)

type Database struct {
	Server string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "secrets"
)

func (m *Database) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

func (m *Database) Insert(s Secret) error {
	s.ID = bson.NewObjectId()
	err := db.C(COLLECTION).Insert(&s)
	return err
}

func (m *Database) Find(hash string) (Secret, error) {
	var secret Secret
	err := db.C(COLLECTION).Find(bson.M{"hash":hash}).One(&secret)
	return secret, err
}

func (m *Database) Update(secret Secret) error {
	err := db.C(COLLECTION).UpdateId(secret.ID, &secret)
	return err
}