package model

type User struct {
	Id        int
	Name      string
	Addresses []Address
}

type Address struct {
	Id      int
	Type    string
	Address string
	City    string
	Country string
}
