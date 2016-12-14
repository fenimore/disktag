package main

import (
	"fmt"
	"rsc.io/pdf"
)

var pages []pdf.Page
var content []pdf.Content

func main() {
	reader, err := pdf.Open("observations.pdf")
	if err != nil {
		fmt.Println(err)
	}

	pages = make([]pdf.Page, reader.NumPage())
	content = make([]pdf.Content, reader.NumPage())

	for i := 1; i <= len(pages); i++ {
		pages[i-1] = reader.Page(i)
	}

	for i, page := range pages {
		content[i] = page.Content()
	}
	return

	for _, val := range content[20].Text {
		fmt.Print(val.S) //, val.W)
	}
	// fmt.Println("content")
	// fmt.Println(len(content[1].Text))
}
