package main

import (
	"fmt"
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
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func (pointerToPage *Page) save() error {
	fileName := pointerToPage.Title + ".txt"
	return os.WriteFile(fileName, pointerToPage.Body, syscall.O_RDWR)
}

func loadPage(title string) (*Page, error) {
	fileName := title + ".txt"
	body, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return &Page{title, body}, nil
}

func handleRoot(rsWriter http.ResponseWriter, rq *http.Request) {
	// rq.URL.Path[1:] creates a sub-slice of Path from the 1st character to the end, this drops the leading "/" from the path name
	fmt.Fprintf(rsWriter, "Hi, %s!", rq.URL.Path[1:])
}

func handleView(rsWriter http.ResponseWriter, rq *http.Request) {
	title := rq.URL.Path[len("/view/"):]
	page, err := loadPage(title)
	if err != nil {
		rsWriter.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(rsWriter, "<h1>Warning!</h1><div>Requested page %s not found</div>", title)
	} else {
		rsWriter.WriteHeader(http.StatusOK)
		fmt.Fprintf(rsWriter, "<h1>%s</h1><div>%s</div>", page.Title, page.Body)
	}
}

type Page struct {
	Title string
	Body  []byte
}
