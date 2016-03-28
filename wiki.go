/*
The tutorial https://golang.org/doc/articles/wiki/
*/
package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

//My own view handler
func viewHandler(w http.ResponseWriter, r *http.Request) {
	//removes /view/ from the path
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	//Using template
	//template.ParseFiles will read the contetns of edit.html and return a *template
	t, _ := template.ParseFiles("edit.html")
	//Executes the template, writing the generated HTML to the http.ResponseWriter
	// .Title and .Body
	t.Execute(w, p)
	//Template directives are enclosed in double curly braces. The printf "%s" .Body instruction
	//is a function call that outputs .Body as a string instead of a stream of bytes
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Put logic here to save")
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
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)

	// run the program http://localhost:8080/edit/TestPage
	// to see: a form to edit the page using /save/ handler

}
