package db

import (
	"database/sql"
	"fmt"

	_ "github.com/sijms/go-ora/v2"
)

type Oracle struct {
	DbProperties
}

func (p Oracle) OpenConn() *sql.DB {
	url := "oracle://" + p.Username + ":" + p.Password + "@" + p.Hostname + ":" + p.Port + "/" + p.Dbname
	db, err := sql.Open("oracle", url)
	if err != nil {
		panic(fmt.Errorf("error in sql.Open: %w", err))
	}
	fmt.Println("Connection Open to Oracle Success!")
	return db
}

func (p Oracle) CloseConn(c *sql.DB) {
	err := c.Close()
	if err != nil {
		fmt.Println("Can't close connection: ", err)
	}
	fmt.Println("Connection to Oracle Closed!")
}
