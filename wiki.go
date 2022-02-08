package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"syscall"
)

func main() {
	testPage := &Page{Title: "test", Body: []byte("Welcome to the test page!")}
	testPage.save()

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/view/", handleView)
	http.HandleFunc("/edit/", handleEditView)
	http.HandleFunc("/save/", handleSaveView)

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

func handleView(rsWriter http.ResponseWriter, rq *http.Request) {
	title := rq.URL.Path[len("/view/"):]
	page, err := loadPage(title)
	if err != nil {
		renderTemplate(rsWriter, "notFoundPage", &Page{Title: title, Body: []byte(err.Error())})
	} else {
		renderTemplate(rsWriter, "view", page)
	}
}

func handleEditView(rsWriter http.ResponseWriter, rq *http.Request) {
	title := rq.URL.Path[len("/edit/"):]
	page, err := loadPage(title)
	if err != nil {
		renderTemplate(rsWriter, "notFoundPage", &Page{Title: title, Body: []byte(err.Error())})
	} else {
		renderTemplate(rsWriter, "edit", page)
	}
}

func handleSaveView(rsWriter http.ResponseWriter, rq *http.Request) {
	title := rq.URL.Path[len("/save/"):]
	bodyStr := rq.FormValue("body")
	createdPage := &Page{Title: title, Body: []byte(bodyStr)}
	renderTemplate(rsWriter, "saved", createdPage)
	// page.save()
	// http.Redirect(rsWriter, rq, fmt.Sprintf("/view/%s", title), http.StatusCreated)
}

func loadPage(title string) (*Page, error) {
	fileName := title + ".txt"
	body, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return &Page{title, body}, nil
}

func renderTemplate(rsWriter http.ResponseWriter, viewName string, page *Page) {
	viewPath := fmt.Sprintf("./template/%s.html", viewName)
	template, err := template.ParseFiles(viewPath)
	if err != nil {
		http.Error(rsWriter, err.Error(), http.StatusInternalServerError)
	}
	err = template.Execute(rsWriter, page)
	if err != nil {
		http.Error(rsWriter, err.Error(), http.StatusInternalServerError)
	}
}

type Page struct {
	Title string
	Body  []byte
}
