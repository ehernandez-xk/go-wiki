/*
The tutorial https://golang.org/doc/articles/wiki/
*/
package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

// To render the templates
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	/*
		The http.Error function sends a specified HTTP response code (in this case
		"Internal Server Error") and error message. Already the decision to put this
		in a separate function is paying off.
	*/
}

//My own view handler
func viewHandler(w http.ResponseWriter, r *http.Request) {
	//removes /view/ from the path
	title := r.URL.Path[len("/view/"):]
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
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	//Using template
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
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

	// when xxx doesn't not exist
	// run the program http://localhost:8080/view/xxxxx
	// to see: The redirect is called to the /edit/ page and you can save it.

}
