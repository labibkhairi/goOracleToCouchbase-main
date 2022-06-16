package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/couchbase/gocb/v2"
	_ "github.com/joho/godotenv/autoload"
	"prima-integrasi.com/fendiya/goOracleToCouchbase/pkg/db"
	"prima-integrasi.com/fendiya/goOracleToCouchbase/pkg/model"
)

var myEnv map[string]string
var couchDB db.ConnectionCouchbaseSDK
var oracleDB db.Connection

func init() {
	//DB Declaration
	couchDB = db.Couchbase{
		db.DbProperties{
			Hostname: os.Getenv("DB_HOSTNAME"),
			Port:     os.Getenv("DB_PORT"),
			Dbname:   os.Getenv("DB_NAME"),
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD")}}
	oracleDB = db.Oracle{
		db.DbProperties{
			Hostname: os.Getenv("DB_HOSTNAME1"),
			Port:     os.Getenv("DB_PORT1"),
			Dbname:   os.Getenv("DB_NAME1"),
			Username: os.Getenv("DB_USERNAME1"),
			Password: os.Getenv("DB_PASSWORD1")}}

}

type hotel struct {
	Name string
}

func main() {
	log.Println("Welcome to Go Oracle to Couchbase")
	// gocb.SetLogger(gocb.VerboseStdioLogger())
	users := getUsersFromOracle()
	insertUsersToCouchbase(users)
	log.Println(users)

}

func selectFromCouchbase() {
	log.Println("Test Select to Couchbase")
	//OPEN DB
	c := couchDB.OpenConn()

	//select data from collection
	// Perform a N1QL Query
	queryResult, err := c.Query("SELECT name FROM `travel-sample`.inventory.hotel LIMIT 1;", &gocb.QueryOptions{})
	if err != nil {
		log.Fatal(err)
	}
	var hotels []hotel

	for queryResult.Next() {
		var h hotel // this could also just be an interface{} type
		err := queryResult.Row(&h)
		if err != nil {
			panic(err)
		}
		hotels = append(hotels, h)
	}
	fmt.Println(len(hotels))

	// always check for errors after iterating
	err = queryResult.Err()
	if err != nil {
		panic(err)
	}

	log.Println(hotels)
}

func insertToCouchbase() {
	log.Println("Test insert to Couchbase")
	//OPEN DB
	c := couchDB.OpenConn()
	//OPEN bucket, scope and collection
	bucket := c.Bucket(os.Getenv("DB_NAME"))

	err := bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		log.Fatal(err)
	}

	col := bucket.Scope("tenant_agent_00").Collection("users")

	//user struct
	type User struct {
		Name      string   `json:"name"`
		Email     string   `json:"email"`
		Interests []string `json:"interests"`
	}
	// Create and store a Document
	_, err = col.Upsert("3",
		User{
			Name:      "Fendi",
			Email:     "fendiya@couchbase.com",
			Interests: []string{"Holy Grail", "African Swallows"},
		}, nil)

	if err != nil {
		log.Fatal(err)
	}
}

func getUsersFromOracle() []model.User {
	log.Println("Test Select to Oracle")
	//OPEN DB
	c := oracleDB.OpenConn()
	defer oracleDB.CloseConn(c)

	//select data from Oracle
	queryResult, err := c.Query("SELECT ID,NAME FROM USERS")
	if err != nil {
		log.Fatal(err)
	}

	//print data
	var users []model.User
	for queryResult.Next() {
		var user model.User
		queryResult.Scan(&user.Id, &user.Name)

		//fetch address
		user.Addresses = getAddressofUsers(user.Id, c)
		//compile users array
		users = append(users, user)
	}
	log.Println(users)
	return users
}

func getAddressofUsers(userId int, c *sql.DB) []model.Address {
	log.Println("Address of users")
	query := "SELECT ID,Type,Address,City,Country FROM Address where Users_id = " + strconv.Itoa(userId)
	queryResult, err := c.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	var addresses []model.Address
	for queryResult.Next() {
		var address model.Address
		queryResult.Scan(&address.Id, &address.Type, &address.Address, &address.City, &address.Country)
		addresses = append(addresses, address)
	}

	return addresses
}

func insertUsersToCouchbase(users []model.User) {
	log.Println("Test insert to Couchbase")

	//OPEN DB
	c := couchDB.OpenConn()
	//OPEN bucket, scope and collection
	bucket := c.Bucket(os.Getenv("DB_NAME"))

	err := bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		log.Fatal(err)
	}

	col := bucket.Scope("tenant_agent_00").Collection("users")

	//user struct for couchbase
	type User struct {
		Name    string
		Address []model.Address
	}

	// Create and store a Document
	for _, user := range users {
		_, err = col.Upsert(strconv.Itoa(user.Id),
			User{
				Name:    user.Name,
				Address: user.Addresses,
			}, nil)

		if err != nil {
			log.Fatal(err)
		}

	}

}
