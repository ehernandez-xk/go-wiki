/*
The tutorial https://golang.org/doc/articles/wiki/
*/
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//My own handler
func myhandler(w http.ResponseWriter, r *http.Request) {
	//If I use r.URL.Path the path comes with /dogs
	//instead of I use [1:] to ignore the slash
	//Fprint helps to write to the response
	fmt.Fprintf(w, "Hi there, I love %s", r.URL.Path[1:])
}

//Page represents a page of the wiki
type Page struct {
	Title string
	Body  []byte
}

//To save a page in a file
func (p *Page) save() error {
	filename := p.Title + ".txt"
	// writes to the file
	return ioutil.WriteFile(filename, p.Body, 0600)
}

//To load a page from the file
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func main() {
	http.HandleFunc("/", myhandler)
	http.ListenAndServe(":8080", nil)

	// run the program http://localhost:8080/dogs
	// to see: Hi there, I love dogs

}
