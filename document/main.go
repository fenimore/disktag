package main

import (
	"fmt"
	"rsc.io/pdf"
)

func main() {
	reader, err := pdf.Open("zen.pdf")
	if err != nil {
		fmt.Println(err)
	}

	type Page struct {
		page    pdf.Page
		content pdf.Content
		font    []pdf.Font
	}

	pages := make([]Page, reader.NumPage())

	for i := 1; i <= len(pages); i++ {
		page := reader.Page(i)
		pages[i-1].page = page
		pages[i-1].content = page.Content()
		fmt.Print(page.Resources())
		// fonts := page.Fonts()
		// for f := range fonts {
		//	fmt.Print(page.Font(fonts[f]).BaseFont())
		//	fmt.Print(page.Font(fonts[f]).Widths())
		// }
		fmt.Println()
	}

}
