/*
The tutorial https://golang.org/doc/articles/wiki/
*/
package main

import (
	"fmt"
	"io/ioutil"
)

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
	//Create a page and save it
	p1 := &Page{Title: "TestPage", Body: []byte("this is a sample page, cool ")}
	p1.save()
	//Load the page
	p2, _ := loadPage(p1.Title)
	fmt.Println(string(p2.Body))
}
