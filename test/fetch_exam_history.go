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

	_, err = services.GetDashboard(usn, dob)
	if err != nil {
		panic(err)
	}

	targetUrl := fmt.Sprintf("https://nmitparents.contineo.in/index.php?option=com_history&task=getResult&usn=%s", usn)
	req, _ := http.NewRequest("GET", targetUrl, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Referer", "https://nmitparents.contineo.in/index.php?option=com_studentdashboard")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	fmt.Println("Final URL:", res.Request.URL.String())
	fmt.Println("Status Code:", res.StatusCode)

	bodyBytes, _ := io.ReadAll(res.Body)
	os.WriteFile("exam_history_real.html", bodyBytes, 0644)
	fmt.Println("Saved exam_history_real.html, length:", len(bodyBytes))
}
