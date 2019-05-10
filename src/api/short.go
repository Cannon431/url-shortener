package api

import (
	"crypto/md5"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

type Short struct {
	Session    *mgo.Session
	DB         *mgo.Database
	Collection *mgo.Collection
}

func New(url, db, collection string) (Short, error) {
	var short Short

	session, err := mgo.Dial(url)
	if err != nil {
		return short, err
	}

	d := session.DB(db)
	c := d.C(collection)

	return Short{
		Session:    session,
		DB:         d,
		Collection: c,
	}, nil
}

func (s *Short) Has(q interface{}) (bool, error) {
	hasDocument := false
	count, err := s.Collection.Find(q).Count()

	if err != nil {
		return false, err
	}

	if count > 0 {
		hasDocument = true
	}

	return hasDocument, nil
}

func (s *Short) Get(q interface{}, isNew bool, scheme, host, port string) (URL, error) {
	var url URL
	err := s.Collection.Find(q).One(&url)

	if err != nil {
		return url, err
	}

	url.IsNew = isNew
	url.Address = s.makeAddress(scheme, host, port, url.Hash)
	return url, nil
}

func (s *Short) makeAddress(scheme, host, port, hash string) string {
	return fmt.Sprintf("%s://%s:%s/%s", scheme, host, port, hash)
}

func (s *Short) Insert(link, scheme, host, port string) (URL, error) {
	url := URL{
		Link:        link,
		Hash:        s.createHash(link),
		Transitions: 0,
		IsNew:       true,
	}

	err := s.Collection.Insert(&url)

	if err != nil {
		return url, err
	}

	url.Address = s.makeAddress(scheme, host, port, url.Hash)

	return url, nil
}

func (s *Short) IncTransitions(hash string) error {
	return s.Collection.Update(bson.M{"hash": hash}, bson.M{"$inc": bson.M{"transitions": 1}})
}

func (s *Short) createHash(link string) string {
	link += strconv.FormatInt(time.Now().Unix(), 10)
	return fmt.Sprintf("%x", md5.Sum([]byte(link)))[:8]
}
