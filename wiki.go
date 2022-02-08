package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	storedPage := &Page{Title: "firstpage", Body: []byte("The content of the first page")}
	storedPage.save()
	
	loadedPage, _ := loadPage("firstpage")
	fmt.Println(string(loadedPage.Body))
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

type Page struct {
	Title string
	Body  []byte
}
