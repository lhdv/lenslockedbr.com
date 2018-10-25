package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "192.168.56.101"
	port     = 5432
	user     = "developer"
	password = "1234qwer"
	dbname   = "lenslockedbr_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	var id int
	row := db.QueryRow(`INSERT INTO users(name, email) 
                            VALUES($1, $2) RETURNING id`,
		"Jon Calhoun",
		"jon@calhoun.io")
	err = row.Scan(&id)
	if err != nil {
		panic(err)
	}

	fmt.Println("Inserted record. ID = ", id)

	db.Close()
}
