/*
The tutorial https://golang.org/doc/articles/wiki/
*/
package main

import (
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

//regular expresion to validate
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

//Initial tamplates, Then we can use the ExecuteTemplate method to render a
//specific template.
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

//validates and return the title
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Title")
	}
	return m[2], nil // the title is the second subexpression.
}

// To render the templates
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	/*
		The http.Error function sends a specified HTTP response code (in this case
		"Internal Server Error") and error message. Already the decision to put this
		in a separate function is paying off.

		The function template.Must is a convenience wrapper that panics when passed
		a non-nil error value, and otherwise returns the *Template unaltered. A panic
		is appropriate here; if the templates can't be loaded the only sensible thing
		to do is exit the program.

		The ParseFiles function takes any number of string arguments that identify
		our template files, and parses those files into templates that are named
		after the base file name. If we were to add more templates to our program,
		we would add their names to the ParseFiles call's arguments.

	*/
}

//My own view handler
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		// If the page doesn¡t exist redirect to /edit/ page
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	//Using template
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	//Using template
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err = p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)

	//Any errors that occur during p.save() will be reported to the user.
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

	// run the program http://localhost:8080/view/test-test
	// to see: 404 page not found

}
