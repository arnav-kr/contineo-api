package parser

import (
	"fmt"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/arnav-kr/contineo-api/models"
)

var colMatchRegex = regexp.MustCompile(`columns\s*:\s*\[(\[.*?\](?:\s*,\s*\[.*?\])*)\s*,?\s*\]`)
var itemMatchRegex = regexp.MustCompile(`\["([^"]+)"\s*,\s*(-?\d+(?:\.\d+)?)\]`)

func ExtractBBColumns(html, bindtoId string) map[string]float64 {
	result := make(map[string]float64)
	targetStr := `bindto: "#` + bindtoId + `"`
	targetIdx := strings.Index(html, targetStr)
	if targetIdx == -1 {
		return result
	}

	prefix := html[:targetIdx]
	colsIdx := strings.LastIndex(prefix, "columns:")
	if colsIdx == -1 {
		return result
	}

	colSlice := prefix[colsIdx:]
	match := colMatchRegex.FindStringSubmatch(colSlice)
	if len(match) < 2 {
		return result
	}

	items := itemMatchRegex.FindAllStringSubmatch(match[1], -1)
	for _, item := range items {
		if len(item) == 3 {
			code := item[1]
			val, err := strconv.ParseFloat(item[2], 64)
			if err == nil {
				result[code] = val
			}
		}
	}

	return result
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}

func ParseOverview(html string, base ...string) (*models.Overview, error) {
	baseURL := "https://nmitparents.contineo.in"
	if len(base) > 0 && base[0] != "" {
		baseURL = strings.TrimRight(base[0], "/")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	// Student Data
	metaPartsRaw := doc.Find(".cn-stu-data1 p").First().Text()
	metaParts := strings.Split(metaPartsRaw, ",")
	for i := range metaParts {
		metaParts[i] = strings.TrimSpace(metaParts[i])
	}

	student := models.Student{
		Name: strings.TrimSpace(doc.Find(".cn-stu-data h3").First().Text()),
		USN:  strings.TrimSpace(doc.Find(".cn-stu-data1 h2").First().Text()),
	}
	if len(metaParts) > 0 && metaParts[0] != "" {
		student.Program = strPtr(metaParts[0])
	}
	if len(metaParts) > 1 && metaParts[1] != "" {
		student.Semester = strPtr(metaParts[1])
	}
	if len(metaParts) > 2 && metaParts[2] != "" {
		student.Section = strPtr(metaParts[2])
	}

	// Chart Data
	attendance := ExtractBBColumns(html, "gaugeTypeMulti")
	cie := ExtractBBColumns(html, "barPadding")

	// Course Rows
	var courses []models.Course
	doc.Find("table.cn-pay-table tbody tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() < 2 {
			return
		}

		code := strings.TrimSpace(cells.Eq(0).Text())
		name := strings.TrimSpace(cells.Eq(1).Text())
		if code == "" || name == "" {
			return
		}

		attHref, _ := row.Find(`a[href*="attendencelist"]`).Attr("href")
		cieHref, _ := row.Find(`a[href*="ciedetails"]`).Attr("href")
		notesHref, _ := row.Find(`a[href*="getuploadlession"]`).Attr("href")

		var courseId, secId, semId *string
		if attHref != "" {
			u, err := url.Parse(attHref)
			if err == nil {
				q := u.Query()
				if q.Has("courseId") {
					courseId = strPtr(q.Get("courseId"))
				}
				if q.Has("secId") {
					secId = strPtr(q.Get("secId"))
				}
				if q.Has("semId") {
					semId = strPtr(q.Get("semId"))
				}
			}
		}

		course := models.Course{
			Code: code,
			Name: name,
			IDs: models.IDs{
				CourseID: courseId,
				SecID:    secId,
				SemID:    semId,
			},
		}

		if val, exists := attendance[code]; exists {
			course.Attendance.Percentage = floatPtr(val)
		}
		if attHref != "" {
			fullUrl := baseURL + "/" + strings.TrimLeft(attHref, "/")
			course.Attendance.DetailURL = strPtr(fullUrl)
		}

		if val, exists := cie[code]; exists {
			course.CIE.Marks = floatPtr(val)
		}
		if cieHref != "" {
			fullUrl := baseURL + "/" + strings.TrimLeft(cieHref, "/")
			course.CIE.DetailURL = strPtr(fullUrl)
		}

		if notesHref != "" {
			fullUrl := baseURL + "/" + strings.TrimLeft(notesHref, "/")
			course.LessonNotesURL = strPtr(fullUrl)
		}

		courses = append(courses, course)
	})

	// Summary
	withAttCount := 0
	sumAtt := 0.0
	under85 := []string{}
	under75 := []string{}

	for _, c := range courses {
		if c.Attendance.Percentage != nil {
			withAttCount++
			sumAtt += *c.Attendance.Percentage
			if *c.Attendance.Percentage < 85 {
				under85 = append(under85, c.Code)
			}
			if *c.Attendance.Percentage < 75 {
				under75 = append(under75, c.Code)
			}
		}
	}

	summary := models.Summary{
		TotalCourses: len(courses),
		BelowThreshold: models.Thresholds{
			Under85: under85,
			Under75: under75,
		},
	}
	
	// ensure they encode to [] instead of null in JSON
	if summary.BelowThreshold.Under85 == nil {
		summary.BelowThreshold.Under85 = []string{}
	}
	if summary.BelowThreshold.Under75 == nil {
		summary.BelowThreshold.Under75 = []string{}
	}

	if withAttCount > 0 {
		avg := sumAtt / float64(withAttCount)
		avg = math.Round(avg*100) / 100
		summary.AverageAttendance = floatPtr(avg)
	}

	overview := &models.Overview{
		Student: student,
		Courses: courses,
		Summary: summary,
	}

	return overview, nil
}

func ParseAttendanceDetail(html string) (*models.AttendanceData, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	stats := models.AttendanceStats{}
	var records []models.AttendanceRecord

	doc.Find("table.uk-table tbody tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() < 4 {
			return
		}

		date := strings.TrimSpace(cells.Eq(1).Text())
		// Convert DD-MM-YYYY to YYYY-MM-DD
		dateParts := strings.Split(date, "-")
		if len(dateParts) == 3 {
			date = fmt.Sprintf("%s-%s-%s", dateParts[2], dateParts[1], dateParts[0])
		}
		
		timeRaw := strings.TrimSpace(cells.Eq(2).Text())
		timeStr := strings.Join(strings.Fields(timeRaw), " ")
		
		startTime := ""
		endTime := ""
		timeParts := strings.Split(timeStr, " TO ")
		if len(timeParts) == 2 {
			startTime = strings.TrimSpace(timeParts[0])
			endTime = strings.TrimSpace(timeParts[1])
		} else {
			startTime = timeStr
		}
		
		status := strings.ToLower(strings.TrimSpace(cells.Eq(3).Text()))
		isPresent := status == "present"

		records = append(records, models.AttendanceRecord{
			Date:      date,
			StartTime: startTime,
			EndTime:   endTime,
			IsPresent: isPresent,
			Status:    status,
		})

		stats.TotalClasses++
		if status == "present" {
			stats.Present++
		}
		if status == "absent" {
			stats.Absent++
		}
	})
	
	if records == nil {
		records = []models.AttendanceRecord{}
	}

	if stats.TotalClasses > 0 {
		stats.Percentage = math.Round((float64(stats.Present)/float64(stats.TotalClasses))*10000) / 100
	}

	return &models.AttendanceData{
		Stats:   stats,
		Records: records,
	}, nil
}

func ParseMarks(raw string) models.Marks {
	m := models.Marks{Raw: raw}
	if raw == "-" || raw == "" {
		return m
	}
	parts := strings.Split(raw, "/")
	if len(parts) > 0 {
		if val, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64); err == nil {
			m.Obtained = floatPtr(val)
		}
	}
	if len(parts) > 1 {
		if val, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64); err == nil {
			m.Max = floatPtr(val)
		}
	}
	return m
}

func ParseCIEDetail(html string) (*models.CIEData, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	data := &models.CIEData{
		Components: []models.CIEComponent{},
	}

	doc.Find("table").Each(func(i int, table *goquery.Selection) {
		headers := []string{}
		table.Find("tr").Eq(0).Find("th, td").Each(func(j int, cell *goquery.Selection) {
			headers = append(headers, strings.TrimSpace(cell.Text()))
		})

		// Check if this table has CIE headers
		hasCIE := false
		for _, h := range headers {
			if strings.Contains(h, "MSE") || strings.Contains(h, "Final IA") || strings.Contains(h, "CIE") {
				hasCIE = true
				break
			}
		}

		if hasCIE && len(headers) > 0 {
			table.Find("tr").Eq(1).Find("th, td").Each(func(j int, cell *goquery.Selection) {
				if j < len(headers) {
					header := headers[j]
					val := strings.TrimSpace(cell.Text())
					// Remove inline scripts and extra spaces
					if idx := strings.Index(val, "document.getElementById"); idx != -1 {
						val = val[:idx]
					}
					val = strings.Join(strings.Fields(val), " ")

					if header != "" {
						if header == "Final IA" || header == "Final CIE" {
							data.FinalCIE = ParseMarks(val)
						} else if header == "Attendance" {
							attRaw := strings.ReplaceAll(val, "%", "")
							if p, err := strconv.ParseFloat(attRaw, 64); err == nil {
								data.Attendance = floatPtr(p)
							}
						} else {
							data.Components = append(data.Components, models.CIEComponent{
								Name:  header,
								Marks: ParseMarks(val),
							})
						}
					}
				}
			})
		}
	})

	return data, nil
}

func ParseExamHistory(html string) (*models.ExamHistory, error) {
        doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
        if err != nil {
                return nil, err
        }

        history := &models.ExamHistory{}

	// Extract top level details
	doc.Find(".cn-last-update").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if strings.Contains(text, "Last Updated On:") {
			rawDate := strings.TrimSpace(strings.Replace(text, "Last Updated On:", "", 1))
			parts := strings.Split(rawDate, "/")
			if len(parts) == 3 {
				history.LastUpdated = fmt.Sprintf("%s-%s-%s", parts[2], parts[1], parts[0])
			} else {
				history.LastUpdated = rawDate
			}
		}
	})

	// Extract student details
	doc.Find("span.cn-table-heading").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		val := strings.TrimSpace(s.Next().Text())
		valDirect := strings.TrimSpace(s.Parent().Text())

		if strings.Contains(text, "Student Name") {
			history.Student.Name = val
		} else if strings.Contains(text, "USN") {
			history.Student.USN = val
		} else if strings.Contains(text, "Stream & Department") {
			parts := strings.Split(val, " - ")
			if len(parts) > 0 {
				prog := strings.TrimSpace(parts[0])
				history.Student.Program = &prog
			}
			if len(parts) > 1 {
				dept := strings.TrimSpace(parts[1])
				// store in section or somewhere if dept needed, using Section for now
				history.Student.Section = &dept
			}
		} else if strings.Contains(text, "Father Name") {
			history.Student.FatherName = &val
		} else if strings.Contains(text, "Year of Admission") {
			history.Student.YearOfAdmission = &val
		} else if strings.Contains(text, "Credits required to earn degree") {
			actualVal := strings.Replace(valDirect, text, "", 1)
			actualVal = strings.TrimSpace(actualVal)
			reqVal, err := strconv.ParseFloat(actualVal, 64)
			if err == nil {
				history.DegreeStats.CreditsRequired = &reqVal
			}
		}
	})

	// Extract degree stats boxes
	doc.Find(".credits-sec1").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h3").Text())
		titleLower := strings.ToLower(strings.Join(strings.Fields(title), " "))
		valStr := strings.TrimSpace(s.Find("p").Text())
		val, err := strconv.ParseFloat(valStr, 64)
		if err == nil {
			if strings.Contains(titleLower, "earned so far") {
				history.DegreeStats.CreditsEarned = &val
			} else if strings.Contains(titleLower, "to be earned") {
				history.DegreeStats.CreditsToEarn = &val
			} else if titleLower == "cgpa" {
				history.DegreeStats.CGPA = &val
			}
		}
	})

	// Extract semesters
	doc.Find("ul[uk-accordion] > li").Each(func(i int, s *goquery.Selection) {
                sem := models.SemesterResult{}

                // Parse accordion title
                titleLink := s.Find("a.cn-acc-title")
                
                titleClone := titleLink.Clone()
                titleClone.Find("span").Remove()
                sem.SemesterName = strings.TrimSpace(titleClone.Text())

                // extract labels
                titleLink.Find("span.uk-label").Each(func(j int, span *goquery.Selection) {
                        text := strings.TrimSpace(span.Text())
                        parts := strings.Split(text, ":")
                        if len(parts) == 2 {
                                key := strings.ToLower(strings.TrimSpace(parts[0]))
                                valStr := strings.TrimSpace(parts[1])
                                val, _ := strconv.ParseFloat(valStr, 64)

                                if strings.Contains(key, "registered") {
                                        sem.CreditsRegistered = &val
                                } else if strings.Contains(key, "earned") {
                                        sem.CreditsEarned = &val
                                } else if strings.Contains(key, "sgpa") {
                                        sem.SGPA = &val
                                } else if strings.Contains(key, "cgpa") {
                                        sem.CGPA = &val
                                }
                        }
                })

                // parse courses
                s.Find("table.res-table tbody tr").Each(func(k int, tr *goquery.Selection) {
                        tds := tr.Find("td")
                        if tds.Length() >= 6 {
                                crStr := strings.TrimSpace(tds.Eq(2).Text())
                                ceStr := strings.TrimSpace(tds.Eq(3).Text())
                                
                                var gpaPtr *float64
                                gpaStr := strings.TrimSpace(tds.Eq(4).Text())
                                gpa, err := strconv.ParseFloat(gpaStr, 64)
                                if err == nil {
                                        gpaPtr = &gpa
                                }

                                course := models.ExamCourse{
                                        Code:             strings.TrimSpace(tds.Eq(0).Text()),
                                        Name:             strings.TrimSpace(tds.Eq(1).Text()),
                                        CreditRegistered: crStr,
                                        CreditEarned:     ceStr,
                                        GPA:              gpaPtr,
                                        Grade:            strings.TrimSpace(tds.Eq(5).Text()),
                                }
                                sem.Courses = append(sem.Courses, course)
                        }
                })

                if sem.SemesterName != "" {
                        history.Semesters = append(history.Semesters, sem)
                }
        })

        return history, nil
}
