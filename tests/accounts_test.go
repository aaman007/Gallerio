package tests

import (
	"fmt"
	"gallerio/accounts"
	"testing"
	"time"
)

func testingUserService() (accounts.UserService, error) {
	const (
		dbHost = "localhost"
		dbPort = 5432
		dbUser = "robert"
		dbPassword = "password"
		dbName = "gallerio_test"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	us, err := accounts.NewUserService(psqlInfo)
	if err != nil {
		return nil, err
	}
	// us.DB.LogMode(false)
	us.DestructiveReset()
	return us, nil
}

func TestCreateUser(t *testing.T) {
	us, err := testingUserService()
	if err != nil {
		t.Fatal(err)
	}

	user := accounts.User{
		Name: "Amanur Rahman",
		Email: "aaman007.liilab@gmail.com",
	}
	err = us.Create(&user)
	if err != nil {
		t.Fatal(err)
	}

	if user.ID <= 0 {
		t.Errorf("Expected ID to be > 0; Received %d", user.ID)
	}
	if time.Since(user.CreatedAt) > 5 * time.Second {
		// Fatalf will terminate, so we are using Errorf
		t.Errorf("Expected CreatedAt to be recent; Received %s", user.CreatedAt)
	}
	if time.Since(user.UpdatedAt) > 5 * time.Second {
		t.Errorf("Expected UpdatedAt to be recent; Received %s", user.UpdatedAt)
	}
}

// go test go-web-dev-2/tests/
// test file name should start with what we are testing followed by _test.go
// eg, accounts_test.go, order_test.go etc
//
// Function name also needs to in specific format