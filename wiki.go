package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"syscall"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html", "saved.html", "notFoundPage.html"))
var titleRegexp = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func main() {
	testPage := &Page{Title: "test", Body: []byte("Welcome to the test page!")}
	testPage.save()

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/view/", makeHandler(handleView))
	http.HandleFunc("/edit/", makeHandler(handleEditView))
	http.HandleFunc("/save/", makeHandler(handleSaveView))

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func (pointerToPage *Page) save() error {
	fileName := pointerToPage.Title + ".txt"
	return os.WriteFile(fileName, pointerToPage.Body, syscall.O_RDWR)
}

func handleRoot(rsWriter http.ResponseWriter, rq *http.Request) {
	// rq.URL.Path[1:] creates a sub-slice of Path from the 1st character to the end, this drops the leading "/" from the path name
	fmt.Fprintf(rsWriter, "Hi, %s!", rq.URL.Path[1:])
}

func handleView(rsWriter http.ResponseWriter, rq *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		renderTemplate(rsWriter, "notFoundPage", &Page{Title: title, Body: []byte(err.Error())})
	} else {
		renderTemplate(rsWriter, "view", page)
	}
}

func handleEditView(rsWriter http.ResponseWriter, rq *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		renderTemplate(rsWriter, "notFoundPage", &Page{Title: title, Body: []byte(err.Error())})
	} else {
		renderTemplate(rsWriter, "edit", page)
	}
}

func handleSaveView(rsWriter http.ResponseWriter, rq *http.Request, title string) {
	bodyStr := rq.FormValue("body")
	createdPage := &Page{Title: title, Body: []byte(bodyStr)}
	err := createdPage.save()
	if err != nil {
		http.Error(rsWriter, err.Error(), http.StatusInternalServerError)
	} else {
		http.Redirect(rsWriter, rq, fmt.Sprintf("/view/%s", title), http.StatusCreated)
	}
}

// func handleView(rsWriter http.ResponseWriter, rq *http.Request) {
// 	title, err := validateAndGetTitle(rsWriter, rq)
// 	if err != nil {
// 		return
// 	}
// 	page, err := loadPage(title)
// 	if err != nil {
// 		renderTemplate(rsWriter, "notFoundPage", &Page{Title: title, Body: []byte(err.Error())})
// 	} else {
// 		renderTemplate(rsWriter, "view", page)
// 	}
// }

// func handleEditView(rsWriter http.ResponseWriter, rq *http.Request) {
// 	title, err := validateAndGetTitle(rsWriter, rq)
// 	if err != nil {
// 		return
// 	}
// 	page, err := loadPage(title)
// 	if err != nil {
// 		renderTemplate(rsWriter, "notFoundPage", &Page{Title: title, Body: []byte(err.Error())})
// 	} else {
// 		renderTemplate(rsWriter, "edit", page)
// 	}
// }

// func handleSaveView(rsWriter http.ResponseWriter, rq *http.Request) {
// 	title, err := validateAndGetTitle(rsWriter, rq)
// 	if err != nil {
// 		return
// 	}
// 	bodyStr := rq.FormValue("body")
// 	createdPage := &Page{Title: title, Body: []byte(bodyStr)}
// 	err = createdPage.save()
// 	if err != nil {
// 		http.Error(rsWriter, err.Error(), http.StatusInternalServerError)
// 	} else {
// 		http.Redirect(rsWriter, rq, fmt.Sprintf("/view/%s", title), http.StatusCreated)
// 	}
// }

func loadPage(title string) (*Page, error) {
	fileName := title + ".txt"
	body, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return &Page{title, body}, nil
}

func renderTemplate(rsWriter http.ResponseWriter, viewName string, page *Page) {
	fullViewName := fmt.Sprintf("%s.html", viewName)
	err := templates.ExecuteTemplate(rsWriter, fullViewName, page)
	if err != nil {
		http.Error(rsWriter, err.Error(), http.StatusInternalServerError)
	}
}

// func validateAndGetTitle(rsWriter http.ResponseWriter, rq *http.Request) (string, error) {
// 	matches := titleRegexp.FindStringSubmatch(rq.URL.Path)
// 	if matches == nil {
// 		http.NotFound(rsWriter, rq)
// 		return "", errors.New("Invalid page title provided")
// 	}
// 	return matches[2], nil // The title is the second subexpression.
// }

func makeHandler(function func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(rsWriter http.ResponseWriter, rq *http.Request) {
		matches := titleRegexp.FindStringSubmatch(rq.URL.Path)
		if matches == nil {
			http.NotFound(rsWriter, rq)
			return
		}
		function(rsWriter, rq, matches[2])
	}
}

type Page struct {
	Title string
	Body  []byte
}
