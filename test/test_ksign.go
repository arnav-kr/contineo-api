package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/arnav-kr/contineo-api/services"
)

func main() {
	usn := "1NT24CS049"
	dob := "2007-09-17"

	client, err := services.GetClientExported(usn, dob)
	if err != nil {
		panic(err)
	}

	url := "https://nmitparents.contineo.in/index.php?ksign=Zou0%2F97gOI%2FQrHOTqcTCraE0Jfz9UfKaGYnFn%2BgjbArW%2F4VXxQ%2BAzguBuMXsyQpHtvEM1z02gwQ%3D"
	
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Referer", "https://nmitparents.contineo.in/index.php?option=com_studentdashboard")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	fmt.Println("Final URL:", res.Request.URL.String())

	bodyBytes, _ := io.ReadAll(res.Body)
	os.WriteFile("ksign_history.html", bodyBytes, 0644)
	fmt.Println("Saved ksign_history.html")
}
