package services

import (
	"strings"

	"github.com/arnav-kr/contineo-api/models"
)

// TestUSN is the reserved demo USN used for mock/test responses.
const TestUSN = "1NT24CS006"

// IsTestStudent checks whether the USN matches the test student account.
func IsTestStudent(usn string) bool {
	return strings.EqualFold(strings.TrimSpace(usn), TestUSN)
}

func strPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}

// GetTestDashboard returns synthetic mock dashboard data for TestUSN.
func GetTestDashboard(portal string) *models.Overview {
	p := NormalizePortal(portal)

	return &models.Overview{
		Portal: p,
		Student: models.Student{
			Name:            "Demo Student",
			USN:             TestUSN,
			Program:         strPtr("Computer Science & Engineering"),
			Semester:        strPtr("4"),
			Section:         strPtr("A"),
			FatherName:      strPtr("Demo Guardian"),
			YearOfAdmission: strPtr("2024"),
		},
		Courses: []models.Course{
			{
				Code: "21CS41",
				Name: "Design and Analysis of Algorithms",
				IDs: models.IDs{
					CourseID: strPtr("4101"),
					SecID:    strPtr("11"),
					SemID:    strPtr("4"),
				},
				Attendance: models.AttendanceDetail{
					Percentage: floatPtr(91.5),
					DetailURL:  strPtr("test://attendance/21CS41"),
				},
				CIE: models.CIEDetail{
					Marks:     floatPtr(44.0),
					DetailURL: strPtr("test://cie/21CS41"),
				},
			},
			{
				Code: "21CS42",
				Name: "Operating Systems",
				IDs: models.IDs{
					CourseID: strPtr("4102"),
					SecID:    strPtr("11"),
					SemID:    strPtr("4"),
				},
				Attendance: models.AttendanceDetail{
					Percentage: floatPtr(84.0),
					DetailURL:  strPtr("test://attendance/21CS42"),
				},
				CIE: models.CIEDetail{
					Marks:     floatPtr(41.5),
					DetailURL: strPtr("test://cie/21CS42"),
				},
			},
			{
				Code: "21CS43",
				Name: "Microcontrollers and Embedded Systems",
				IDs: models.IDs{
					CourseID: strPtr("4103"),
					SecID:    strPtr("11"),
					SemID:    strPtr("4"),
				},
				Attendance: models.AttendanceDetail{
					Percentage: floatPtr(72.0),
					DetailURL:  strPtr("test://attendance/21CS43"),
				},
				CIE: models.CIEDetail{
					Marks:     floatPtr(38.0),
					DetailURL: strPtr("test://cie/21CS43"),
				},
			},
		},
		Summary: models.Summary{
			TotalCourses:      3,
			AverageAttendance: floatPtr(82.5),
			BelowThreshold: models.Thresholds{
				Under85: []string{"21CS42", "21CS43"},
				Under75: []string{"21CS43"},
			},
		},
	}
}

// GetTestAttendanceDetail returns synthetic mock attendance data for TestUSN.
func GetTestAttendanceDetail(courseCode string) *models.AttendanceData {
	code := strings.ToUpper(strings.TrimSpace(courseCode))
	switch code {
	case "21CS42":
		return &models.AttendanceData{
			Stats: models.AttendanceStats{
				TotalClasses: 38,
				Present:      32,
				Absent:       6,
				Percentage:   84.2,
			},
			Records: []models.AttendanceRecord{
				{Date: "2025-02-03", StartTime: "09:00", EndTime: "10:00", IsPresent: true, Status: "Present"},
				{Date: "2025-02-04", StartTime: "10:00", EndTime: "11:00", IsPresent: false, Status: "Absent"},
				{Date: "2025-02-05", StartTime: "11:15", EndTime: "12:15", IsPresent: true, Status: "Present"},
			},
		}
	case "21CS43":
		return &models.AttendanceData{
			Stats: models.AttendanceStats{
				TotalClasses: 36,
				Present:      26,
				Absent:       10,
				Percentage:   72.2,
			},
			Records: []models.AttendanceRecord{
				{Date: "2025-02-03", StartTime: "11:15", EndTime: "12:15", IsPresent: true, Status: "Present"},
				{Date: "2025-02-04", StartTime: "14:00", EndTime: "15:00", IsPresent: false, Status: "Absent"},
				{Date: "2025-02-06", StartTime: "09:00", EndTime: "10:00", IsPresent: false, Status: "Absent"},
			},
		}
	default: // 21CS41 or fallback
		return &models.AttendanceData{
			Stats: models.AttendanceStats{
				TotalClasses: 40,
				Present:      37,
				Absent:       3,
				Percentage:   92.5,
			},
			Records: []models.AttendanceRecord{
				{Date: "2025-02-03", StartTime: "10:00", EndTime: "11:00", IsPresent: true, Status: "Present"},
				{Date: "2025-02-04", StartTime: "09:00", EndTime: "10:00", IsPresent: true, Status: "Present"},
				{Date: "2025-02-05", StartTime: "14:00", EndTime: "15:00", IsPresent: true, Status: "Present"},
				{Date: "2025-02-06", StartTime: "10:00", EndTime: "11:00", IsPresent: false, Status: "Absent"},
			},
		}
	}
}

// GetTestCIEDetail returns synthetic mock CIE component data for TestUSN.
func GetTestCIEDetail(courseCode string) *models.CIEData {
	code := strings.ToUpper(strings.TrimSpace(courseCode))
	switch code {
	case "21CS42":
		return &models.CIEData{
			Components: []models.CIEComponent{
				{Name: "IA-1", Marks: models.Marks{Obtained: floatPtr(24.0), Max: floatPtr(30.0), Raw: "24 / 30"}},
				{Name: "IA-2", Marks: models.Marks{Obtained: floatPtr(26.0), Max: floatPtr(30.0), Raw: "26 / 30"}},
				{Name: "Assignment", Marks: models.Marks{Obtained: floatPtr(18.0), Max: floatPtr(20.0), Raw: "18 / 20"}},
			},
			FinalCIE: models.Marks{
				Obtained: floatPtr(41.5),
				Max:      floatPtr(50.0),
				Raw:      "41.5 / 50",
			},
			Attendance: floatPtr(84.0),
		}
	case "21CS43":
		return &models.CIEData{
			Components: []models.CIEComponent{
				{Name: "IA-1", Marks: models.Marks{Obtained: floatPtr(21.0), Max: floatPtr(30.0), Raw: "21 / 30"}},
				{Name: "IA-2", Marks: models.Marks{Obtained: floatPtr(23.0), Max: floatPtr(30.0), Raw: "23 / 30"}},
				{Name: "Assignment", Marks: models.Marks{Obtained: floatPtr(16.0), Max: floatPtr(20.0), Raw: "16 / 20"}},
			},
			FinalCIE: models.Marks{
				Obtained: floatPtr(38.0),
				Max:      floatPtr(50.0),
				Raw:      "38.0 / 50",
			},
			Attendance: floatPtr(72.0),
		}
	default:
		return &models.CIEData{
			Components: []models.CIEComponent{
				{Name: "IA-1", Marks: models.Marks{Obtained: floatPtr(27.0), Max: floatPtr(30.0), Raw: "27 / 30"}},
				{Name: "IA-2", Marks: models.Marks{Obtained: floatPtr(28.5), Max: floatPtr(30.0), Raw: "28.5 / 30"}},
				{Name: "Assignment", Marks: models.Marks{Obtained: floatPtr(19.0), Max: floatPtr(20.0), Raw: "19 / 20"}},
			},
			FinalCIE: models.Marks{
				Obtained: floatPtr(44.0),
				Max:      floatPtr(50.0),
				Raw:      "44.0 / 50",
			},
			Attendance: floatPtr(91.5),
		}
	}
}

// GetTestExamHistory returns synthetic mock exam history for TestUSN.
func GetTestExamHistory() *models.ExamHistory {
	return &models.ExamHistory{
		LastUpdated: "2025-01-15",
		Student: models.Student{
			Name:            "Demo Student",
			USN:             TestUSN,
			Program:         strPtr("Computer Science & Engineering"),
			Semester:        strPtr("4"),
			Section:         strPtr("A"),
			FatherName:      strPtr("Demo Guardian"),
			YearOfAdmission: strPtr("2024"),
		},
		DegreeStats: models.DegreeStats{
			CreditsRequired: floatPtr(160.0),
			CreditsEarned:   floatPtr(64.0),
			CreditsToEarn:   floatPtr(96.0),
			CGPA:            floatPtr(8.85),
		},
		Semesters: []models.SemesterResult{
			{
				SemesterName:      "Semester 1",
				CreditsRegistered: floatPtr(20.0),
				CreditsEarned:     floatPtr(20.0),
				SGPA:              floatPtr(8.75),
				CGPA:              floatPtr(8.75),
				Courses: []models.ExamCourse{
					{Code: "21MAT11", Name: "Calculus and Linear Algebra", CreditRegistered: "4.0", CreditEarned: "4.0", GPA: floatPtr(9.0), Grade: "A+"},
					{Code: "21PHY12", Name: "Engineering Physics", CreditRegistered: "4.0", CreditEarned: "4.0", GPA: floatPtr(8.5), Grade: "A"},
				},
			},
			{
				SemesterName:      "Semester 2",
				CreditsRegistered: floatPtr(22.0),
				CreditsEarned:     floatPtr(22.0),
				SGPA:              floatPtr(8.90),
				CGPA:              floatPtr(8.83),
				Courses: []models.ExamCourse{
					{Code: "21MAT21", Name: "Advanced Calculus and Laplace Transforms", CreditRegistered: "4.0", CreditEarned: "4.0", GPA: floatPtr(9.0), Grade: "A+"},
					{Code: "21CHE22", Name: "Engineering Chemistry", CreditRegistered: "4.0", CreditEarned: "4.0", GPA: floatPtr(9.0), Grade: "A+"},
				},
			},
			{
				SemesterName:      "Semester 3",
				CreditsRegistered: floatPtr(22.0),
				CreditsEarned:     floatPtr(22.0),
				SGPA:              floatPtr(8.88),
				CGPA:              floatPtr(8.85),
				Courses: []models.ExamCourse{
					{Code: "21CS31", Name: "Data Structures and Applications", CreditRegistered: "4.0", CreditEarned: "4.0", GPA: floatPtr(9.0), Grade: "A+"},
					{Code: "21CS32", Name: "Analog and Digital Electronics", CreditRegistered: "4.0", CreditEarned: "4.0", GPA: floatPtr(8.5), Grade: "A"},
				},
			},
		},
	}
}
