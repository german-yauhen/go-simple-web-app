package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"
)

func main() {
	// storedPage := &Page{Title: "firstpage", Body: []byte("The content of the first page")}
	// storedPage.save()
	
	// loadedPage, _ := loadPage("firstpage")
	// fmt.Println(string(loadedPage.Body))
	http.HandleFunc("/", handleRoot)
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
	fmt.Fprintf(rsWriter, "Hi there! I love %s!", rq.URL.Path[1:])
}

type Page struct {
	Title string
	Body  []byte
}
