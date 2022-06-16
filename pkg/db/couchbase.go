package db

import (
	"crypto/x509"
	"io/ioutil"
	"log"

	"github.com/couchbase/gocb/v2"
)

type Couchbase struct {
	DbProperties
}

func (p Couchbase) OpenConn() *gocb.Cluster {
	log.Printf("Start Opening Connection to Couchbase ...")
	//initialize connection

	pc, err := ioutil.ReadFile("ca.pem")

	// log.Println(p)
	if err != nil {
		log.Fatal("error opening file : ", err)
	}

	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM(pc)

	cluster, err := gocb.Connect("couchbases://"+p.Hostname, gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: p.Username,
			Password: p.Password,
		},
		SecurityConfig: gocb.SecurityConfig{
			TLSRootCAs: roots,
			// WARNING: Do not set this to true in production, only use this for testing!
			TLSSkipVerify: true,
		},
	})
	if err != nil {
		log.Printf("error in sql.Open: %S", err)
	}

	//open Bucket
	// bucket := cluster.Bucket(p.Dbname)

	// err = bucket.WaitUntilReady(5*time.Second, nil)
	// if err != nil {
	// 	log.Printf("error in opening bucket: %S, with error : %w", p.Dbname, err)
	// }

	log.Printf("Connection Open to Couchbase Success!\n")
	return cluster
}
