package main

import (
	"html/template"
	"os"
)

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	data := struct {
		Name string
		City string
	} {"John Smith", "Bay Area"}

	err = t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
