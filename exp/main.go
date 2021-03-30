package main

import (
	"fmt"
	"go-web-dev-2/accounts"
)

const (
	dbHost = "localhost"
	dbPort = 5432
	dbUser = "robert"
	dbPassword = "password"
	dbName = "gallerio"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	us, err := accounts.NewService(psqlInfo)
	if err != nil {
		panic(err)
	}

	defer us.Close()

	// us.DestructiveReset()
	user, err := us.ByID(12)
	if err != nil {
		panic(err)
	}
	fmt.Println(user)
}
