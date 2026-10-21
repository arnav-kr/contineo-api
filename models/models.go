package models

type Student struct {
	Name            string  `json:"name"`
	USN             string  `json:"usn"`
	Program         *string `json:"program"`
	Semester        *string `json:"semester"`
	Section         *string `json:"section"`
	FatherName      *string `json:"fatherName,omitempty"`
	YearOfAdmission *string `json:"yearOfAdmission,omitempty"`
}

type IDs struct {
	CourseID *string `json:"courseId"`
	SecID    *string `json:"secId"`
	SemID    *string `json:"semId"`
}

type AttendanceDetail struct {
	Percentage *float64 `json:"percentage"`
	DetailURL  *string  `json:"detailUrl"`
}

type CIEDetail struct {
	Marks     *float64 `json:"marks"`
	DetailURL *string  `json:"detailUrl"`
}

type Course struct {
	Code           string           `json:"code"`
	Name           string           `json:"name"`
	IDs            IDs              `json:"ids"`
	Attendance     AttendanceDetail `json:"attendance"`
	CIE            CIEDetail        `json:"cie"`
	LessonNotesURL *string          `json:"lessonNotesUrl"`
}

type Thresholds struct {
	Under85 []string `json:"85"`
	Under75 []string `json:"75"`
}

type Summary struct {
	TotalCourses      int         `json:"totalCourses"`
	AverageAttendance *float64    `json:"averageAttendance"`
	BelowThreshold    Thresholds  `json:"belowThreshold"`
}

type Overview struct {
	Portal  string   `json:"portal,omitempty"`
	Student Student  `json:"student"`
	Courses []Course `json:"courses"`
	Summary Summary  `json:"summary"`
}

type AttendanceRecord struct {
	Date      string `json:"date"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	IsPresent bool   `json:"isPresent"`
	Status    string `json:"status"`
}

type AttendanceStats struct {
	TotalClasses int     `json:"totalClasses"`
	Present      int     `json:"present"`
	Absent       int     `json:"absent"`
	Percentage   float64 `json:"percentage"`
}

type AttendanceData struct {
	Stats   AttendanceStats    `json:"stats"`
	Records []AttendanceRecord `json:"records"`
}

type Marks struct {
	Obtained *float64 `json:"obtained"`
	Max      *float64 `json:"max"`
	Raw      string   `json:"raw"`
}

type CIEComponent struct {
	Name  string `json:"name"`
	Marks Marks  `json:"marks"`
}

type CIEData struct {
	Components []CIEComponent `json:"components"`
	FinalCIE   Marks          `json:"finalCIE"`
	Attendance *float64       `json:"attendancePercentage"`
}

type ExamCourse struct {
	Code             string   `json:"code"`
	Name             string   `json:"name"`
	CreditRegistered string   `json:"creditRegistered"`
	CreditEarned     string   `json:"creditEarned"`
	GPA              *float64 `json:"gpa"`
	Grade            string   `json:"grade"`
}

type SemesterResult struct {
	SemesterName      string       `json:"semesterName"`
	CreditsRegistered *float64     `json:"creditsRegistered"`
	CreditsEarned     *float64     `json:"creditsEarned"`
	SGPA              *float64     `json:"sgpa"`
	CGPA              *float64     `json:"cgpa"`
	Courses           []ExamCourse `json:"courses"`
}

type DegreeStats struct {
	CreditsRequired *float64 `json:"creditsRequired,omitempty"`
	CreditsEarned   *float64 `json:"creditsEarned,omitempty"`
	CreditsToEarn   *float64 `json:"creditsToEarn,omitempty"`
	CGPA            *float64 `json:"cgpa,omitempty"`
}

type ExamHistory struct {
	LastUpdated string           `json:"lastUpdated,omitempty"`
	Student     Student          `json:"student"`
	DegreeStats DegreeStats      `json:"degreeStats"`
	Semesters   []SemesterResult `json:"semesters"`
}
