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

	// Just trigger a dashboard fetch to ensure login
	_, err = services.GetDashboard(usn, dob)
	if err != nil {
		panic(err)
	}

	// Fetch Exam History manually
	targetUrl := fmt.Sprintf("https://nmitparents.contineo.in/index.php?option=com_history&task=getResult&usn=%s", usn)
	req, _ := http.NewRequest("GET", targetUrl, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", "https://nmitparents.contineo.in/index.php?option=com_studentdashboard")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	fmt.Println("Final URL after fetch:", res.Request.URL.String())

	bodyBytes, _ := io.ReadAll(res.Body)
	os.WriteFile("exam_history.html", bodyBytes, 0644)
	fmt.Println("Saved exam_history.html")
}
