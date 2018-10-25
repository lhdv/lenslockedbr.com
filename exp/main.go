package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	host     = "192.168.56.101"
	port     = 5432
	user     = "developer"
	password = "1234qwer"
	dbname   = "lenslockedbr_dev"
)

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	fmt.Println("Successfully connected!")

	db.LogMode(true)
	db.AutoMigrate(&User{})

	// name, email := getInfo()
	// u := &User{Name: name, Email: email}
	// err = db.Create(u).Error
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("+v\n", u)

	var u = &User{}
	maxID := 3
	db.Where("id <= ?", maxID).First(&u)
	if db.Error != nil {
		panic(db.Error)
	}
	fmt.Println(u)

	var users []User
	db.Find(&users)
	if db.Error != nil {
		panic(db.Error)
	}
	fmt.Println("Retrieved", len(users), " users")
	fmt.Println(users)
}

func getInfo() (name, email string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(" What is your name?")
	name, _ = reader.ReadString('\n')
	name = strings.TrimSpace(name)
	fmt.Println(" What is your email?")
	email, _ = reader.ReadString('\n')
	email = strings.TrimSpace(email)

	return name, email
}
