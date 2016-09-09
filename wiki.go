/*
The tutorial https://golang.org/doc/articles/wiki/
*/
package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

//To organize the templates and data
var tmplDir = "templates/"
var dataDir = "data/"

//regular expresion to validate
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

//Initial tamplates, Then we can use the ExecuteTemplate method to render a
//specific template.
var templates = template.Must(template.ParseFiles(tmplDir+"edit.html", tmplDir+"view.html"))

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

/*
Catching the error condition in each handler introduces a lot of repeated code.
this function warp each of the handlers does the validation and error checking.
The returned function is called a closure because it encloses values defined
outside of it. In this case, the variable fn (the single argument to makeHandler)
is enclosed by the closure. The variable fn will be one of our save, edit, or view handlers.*/
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		// call the function handler view, save or edit
		fn(w, r, m[2]) // the title is the second subexpression.
	}
}

//To manage when the request comes localhost:8080/ go to home
func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/view/home", http.StatusFound)
	return
}

//My own view handler
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		// If the page doesn¡t exist redirect to /edit/ page
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	//Using template
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	//Using template
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
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
	filename := dataDir + p.Title + ".txt"
	// writes to the file
	return ioutil.WriteFile(filename, p.Body, 0600)
}

//To load a page from the file
func loadPage(title string) (*Page, error) {
	filename := dataDir + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	//ListenAndServe starts an HTTP server with a given address and handler. The handler is usually nil,
	//which means to use DefaultServeMux. Handle and HandleFunc add handlers to DefaultServeMux:
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error detected", err)
	}

	//http://localhost:8080/view/TestPage

	// invalid regexp
	// run the program http://localhost:8080/view/test-test
	// to see: 404 page not found

}
