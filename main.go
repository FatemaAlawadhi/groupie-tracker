package main

import (
	"fmt"
	"log"
	"net/http"
	"groupie-tracker/Handlers"
)

var err error 
func main() {
	http.HandleFunc("/", Handlers.HandleHomePage)
	http.HandleFunc("/stars", Handlers.HandleStarsPage)
	http.HandleFunc("/about", Handlers.HandleAboutPage)
	http.HandleFunc("/stardetails/", Handlers.HandleStarDetailsPage)
	fs := http.FileServer(http.Dir("style"))
	http.Handle("/style/", http.StripPrefix("/style/", fs))
	fmt.Println("starting server at port 8122\n")
	err = http.ListenAndServe(":8122", nil)
	if err != nil {
		log.Fatal(err)
	}
}
