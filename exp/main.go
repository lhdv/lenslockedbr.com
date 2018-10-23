package main

import (
	"html/template"
	"os"
	"time"
)

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	data := struct {
		Name     string
		City     string
		Today    time.Time
                Age      int
                Weight   float64
		Children map[string]int
	} {"John Smith", 
           "Bay Area", 
           time.Now(),
           33,
           88.25,
           map[string]int{"Jhonny":5, "Anne":10}}

	err = t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}

	data1 := struct {
		Name     string
		City     string
		Today    time.Time
                Age      int
                Weight   float64
		Children map[string]int
	} {"Mary Anne", 
           "Texas", 
           time.Now(),
           23,
           58.25,
	   nil}

	err = t.Execute(os.Stdout, data1)
	if err != nil {
		panic(err)
	}
}
