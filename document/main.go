package main

import (
	"fmt"
	"rsc.io/pdf"
)

type Page struct {
	page    pdf.Page
	content pdf.Content
	font    []pdf.Font
}

func main() {
	reader, err := pdf.Open("insurrection.pdf")
	if err != nil {
		fmt.Println(err)
	}

	pages := make([]Page, reader.NumPage())

	for i := 1; i <= len(pages); i++ {
		page := reader.Page(i)
		pages[i-1].page = page
		pages[i-1].content = page.Content()
		//fmt.Print(page.Resources())
		// fonts := page.Fonts()
		// for f := range fonts {
		//	fmt.Print(page.Font(fonts[f]).BaseFont())
		//	fmt.Print(page.Font(fonts[f]).Widths())
		// }
		//fmt.Println()
	}

	for _, val := range pages[10].content.Text {
		//fmt.Println(val)
		fmt.Print(val.S)
	}

	//fmt.Println(pages[10].page.V.Key("Contents"))

}
