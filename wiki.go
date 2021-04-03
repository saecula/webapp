package main

import (
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
)

type Page struct {
	Title string
	Body  []byte
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html", "main.html", "notfound.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-0-9]+)$")
var extraPath = regexp.MustCompile("^/([a-zA-Z0-0-9]+)$")

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil
}

func validateHasNoExtraPath(w http.ResponseWriter, r *http.Request) error {
	m := extraPath.FindStringSubmatch(r.URL.Path)
	if m != nil {
		return errors.New("Tried to go somewhere")
	}
	return nil
}

func (p *Page) save() error {
	filename := "pages/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "pages/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	err := validateHasNoExtraPath(w, r)
	if err != nil {
		renderEmptyTemplate(w, "notfound")
		return
	}
	files, err := ioutil.ReadDir("./pages")
	if err != nil {
		log.Fatal(err)
	}
	var titles []string
	for _, f := range files {
		rawTitle := f.Name()
		title := strings.TrimSuffix(rawTitle, filepath.Ext(rawTitle))
		titles = append(titles, title)
	}
	renderTemplateFromStrSlice(w, "main", titles)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	} else {
		renderPageTemplate(w, "view", p)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderPageTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			renderEmptyTemplate(w, "notfound")
			return
		}
		fn(w, r, m[2])
	}
}

func renderEmptyTemplate(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", &Page{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func renderPageTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func renderTemplateFromStrSlice(w http.ResponseWriter, tmpl string, l []string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", l)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/", mainPageHandler)
	log.Fatal(http.ListenAndServe(":4000", nil))
}
