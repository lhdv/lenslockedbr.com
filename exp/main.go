package main

import (
	"fmt"

	"lenslockedbr.com/models"
	_ "lenslockedbr.com/hash"
	_ "lenslockedbr.com/rand"
)

const (
	host     = "192.168.56.101"
	port     = 5432
	user     = "developer"
	password = "1234qwer"
	dbname   = "lenslockedbr_dev"
)


func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s " +  
                                "dbname=%s sslmode=disable",
		                 host, port, user, password, dbname)

	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}

	defer us.Close()

	us.DestructiveReset()

	fmt.Println("Successfully connected!")

	var user models.User

	user = models.User{ Name: "Foobar", 
			    Age: 10,
                            Email: "foobar@example.com",
			    Password: "321123",
	}
	err = us.Create(&user)
	if err != nil {
		panic(err)
	}

	fmt.Println("User created:", user)

	// Verify that the user has a Remember and RememberHash
	fmt.Printf("%+v\n", user)
	if user.Remember == "" {
		panic("Invalid remember token")
	}

	user2, err := us.ByRemember(user.Remember)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", *user2)

	user = models.User{ Name: "Test", 
			Age: 15,
			Email: "test@example.com",
			Password: "abcdef",
	}
	err = us.Create(&user)
	if err != nil {
		panic(err)
	}

	fmt.Println("User created:", user)

	byId, err := us.ByID(2)
	if err != nil {
		panic(err)
	}

	fmt.Println("User find by ID:", byId)

	byEmail, err := us.ByEmail("foobar@example.com")
	if err != nil {
		panic(err)
	}

	fmt.Println("User find by Email:", byEmail)

	byId.Name = "Updated"
	err = us.Update(byId)
	if err != nil {
		panic(err)
	}

	fmt.Println("User updated:", byId)

	agesInRange, err := us.InAgeRange(1,20)
	if err != nil {
		panic(err)
	}

	fmt.Println("Users on age range of 1,20:", agesInRange)

	err = us.Delete(2)
	if err != nil {
		panic(err)
	}

	fetchById, err := us.ByID(2)
	if err != nil {
		panic(err)
	}
	fmt.Println("User find deleted by ID:", fetchById)
/*

	// Generating random string
	s, _ := rand.String(10)
	fmt.Println(s, len(s))
	s, _ = rand.RememberToken()
	fmt.Println(s, len(s))

	// Hashing a string
	hmac := hash.NewHMAC("my-secret-key")
	fmt.Println(hmac.Hash("this is my string to hash"))
*/	
}
