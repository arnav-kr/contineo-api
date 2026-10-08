package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	html, _ := os.ReadFile("cie_debug.txt")
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		panic(err)
	}

	doc.Find("table").Each(func(i int, table *goquery.Selection) {
		fmt.Printf("Table %d:\n", i)
		table.Find("tr").Each(func(j int, row *goquery.Selection) {
			var rowData []string
			row.Find("th, td").Each(func(k int, cell *goquery.Selection) {
				text := strings.TrimSpace(cell.Text())
				text = strings.ReplaceAll(text, "\n", " ")
				rowData = append(rowData, text)
			})
			fmt.Printf("  Row %d (%d cells): %v\n", j, len(rowData), rowData)
		})
	})
}
