package main

import (
	"fmt"
	"os"

	"github.com/arnav-kr/contineo-api/parser"
)

func main() {
	html, err := os.ReadFile("../attendance_detail.html")
	if err != nil {
		panic(err)
	}

	data, err := parser.ParseAttendanceDetail(string(html))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Total: %d, Present: %d, Absent: %d, Percentage: %.2f\n", data.Stats.TotalClasses, data.Stats.Present, data.Stats.Absent, data.Stats.Percentage)
	fmt.Printf("Records length: %d\n", len(data.Records))
}
