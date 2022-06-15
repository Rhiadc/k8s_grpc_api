package main

import (
	"fmt"

	"github.com/rhiadc/grpc_api/client/domain"
)

func main() {

	user := &domain.User{
		FirstName: "Badger",
		LastName:  "Smith",
		Email:     "Badger.Smith@gmail.com",
	}

	fmt.Println(user.Validate())
}
