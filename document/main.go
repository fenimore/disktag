package main

import (
	"fmt"
	"rsc.io/pdf"
)

var pages []pdf.Page
var content []pdf.Content

func main() {
	reader, err := pdf.Open("trust.pdf")
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

	for _, val := range content[110].Text {
		fmt.Println(val.S, val.W)
	}
	// fmt.Println("content")
	// fmt.Println(len(content[1].Text))
}
