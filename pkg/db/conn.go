package db

import (
	"database/sql"

	"github.com/couchbase/gocb/v2"
)

type Connection interface {
	OpenConn() *sql.DB
	CloseConn(*sql.DB)
}

type ConnectionCouchbaseSDK interface {
	OpenConn() *gocb.Cluster
}

type DbProperties struct {
	Hostname string
	Port     string
	Dbname   string
	Username string
	Password string
}
