package main

import (
	"fmt"

	"github.com/arnav-kr/contineo-api/services"
)

func main() {
	usn := "1NT24CS049"
	dob := "2007-09-17"

	overview, err := services.GetDashboard(usn, dob)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Parsed %d courses\n", len(overview.Courses))

	var cieUrl string
	for _, c := range overview.Courses {
		if c.Code == "22MAT41S" {
			if c.CIE.DetailURL != nil {
				cieUrl = *c.CIE.DetailURL
			}
		}
	}

	fmt.Println("Fetching CIE URL:", cieUrl)
	detail, err := services.GetCIEDetail(usn, dob, cieUrl)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Final CIE: %s, Attendance: %v%%\n", detail.FinalCIE.Raw, func() string {
		if detail.Attendance != nil {
			return fmt.Sprintf("%.0f", *detail.Attendance)
		}
		return "N/A"
	}())
	for _, comp := range detail.Components {
		fmt.Printf("  %s: %s\n", comp.Name, comp.Marks.Raw)
	}
}
