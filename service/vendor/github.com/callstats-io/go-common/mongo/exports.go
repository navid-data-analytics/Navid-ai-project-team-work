package mongo

import (
	mgo "github.com/globalsign/mgo"
	bson "github.com/globalsign/mgo/bson"
)

// M aliased from bson.M
type M bson.M

// mgo Error aliases
var (
	ErrNotFound = mgo.ErrNotFound
	ErrCursor   = mgo.ErrCursor
)

// Session aliased from mgo.Session
type Session struct {
	*mgo.Session
}

// Status checks if the connection is alive and returns the status object
func (ses *Session) Status() (M, error) {
	status := M{}
	err := ses.Run("serverStatus", &status)
	return status, err
}

// struct alias helpers
func asAliasedSession(ses *mgo.Session) *Session {
	return &Session{
		Session: ses,
	}
}
