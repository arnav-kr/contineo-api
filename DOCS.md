# Contineo API Reference

REST API for retrieving student academic data, attendance records, CIE scores, and examination history from Contineo portals.

---

## Base URL

```
/api
```

Live documentation & OpenAPI specification:
- Swagger UI: `/docs`
- OpenAPI Schema: `/openapi.json`
- Health check: `/health` (or `/api/health`)

---

## Authentication & Parameters

Requests authenticate using student credentials passed as query parameters:

| Parameter | Type | Required | Description | Example |
| :--- | :--- | :--- | :--- | :--- |
| `usn` | string | Yes | Student University Seat Number (`1NT24CS006` for mock test mode) | `1NT24CS006` |
| `dob` | string | Yes | Date of birth (`YYYY-MM-DD` format; any value accepted for test USN) | `2000-01-01` |
| `course` | string | Endpoint-specific | Course code (returned in `/api/dashboard`) | `21CS41` |
| `portal` | string | No | Target portal: `even` (default) or `odd` (`parentsodd`). Auto-detects if omitted. | `odd` |

> **Mock Test Student**: Use `usn=1NT24CS006` with any `dob` to test response schemas across all endpoints without real credentials. Demo courses: `21CS41`, `21CS42`, `21CS43`.

Sessions are managed server-side.

---

## Endpoints

### 1. Overview & Courses
```http
GET /api/dashboard?usn={usn}&dob={dob}&portal={portal}
```

Returns profile metadata, enrolled courses, and attendance summary.

#### Response (`200 OK`)
```json
{
  "portal": "odd",
  "student": {
    "name": "Jane Doe",
    "usn": "1XX00XX000",
    "program": "Computer Science & Engineering",
    "semester": "5",
    "section": "A",
    "fatherName": "John Doe",
    "yearOfAdmission": "2022"
  },
  "courses": [
    {
      "code": "21CS51",
      "name": "Management and Entrepreneurship",
      "ids": {
        "courseId": "6939",
        "secId": "26",
        "semId": "96"
      },
      "attendance": {
        "percentage": 92.5,
        "detailUrl": "https://nmitparents.contineo.in/parentsodd/index.php?..."
      },
      "cie": {
        "marks": 45.0,
        "detailUrl": "https://nmitparents.contineo.in/parentsodd/index.php?..."
      },
      "lessonNotesUrl": null
    }
  ],
  "summary": {
    "totalCourses": 6,
    "averageAttendance": 88.4,
    "belowThreshold": {
      "85": ["21CS53"],
      "75": []
    }
  }
}
```

---

### 2. Course Attendance Details
```http
GET /api/attendance?usn={usn}&dob={dob}&course={course}&portal={portal}
```

Returns session-by-session attendance logs and attendance metrics for a specific course.

#### Response (`200 OK`)
```json
{
  "stats": {
    "totalClasses": 40,
    "present": 37,
    "absent": 3,
    "percentage": 92.5
  },
  "records": [
    {
      "date": "2024-09-02",
      "startTime": "09:00",
      "endTime": "10:00",
      "isPresent": true,
      "status": "Present"
    }
  ]
}
```

---

### 3. Continuous Internal Evaluation (CIE)
```http
GET /api/cie?usn={usn}&dob={dob}&course={course}&portal={portal}
```

Returns individual assessment marks (internal tests, quizzes, assignments) and calculated CIE totals.

#### Response (`200 OK`)
```json
{
  "components": [
    {
      "name": "Internal Assessment 1",
      "marks": {
        "obtained": 28.0,
        "max": 30.0,
        "raw": "28 / 30"
      }
    }
  ],
  "finalCIE": {
    "obtained": 46.0,
    "max": 50.0,
    "raw": "46"
  },
  "attendancePercentage": 92.5
}
```

---

### 4. Examination History
```http
GET /api/exam-history?usn={usn}&dob={dob}&portal={portal}
```

Returns semester-by-semester grades, credit completion, SGPA, and cumulative CGPA.

#### Response (`200 OK`)
```json
{
  "lastUpdated": "2024-08-15",
  "student": {
    "name": "Jane Doe",
    "usn": "1XX00XX000",
    "program": "Computer Science & Engineering",
    "semester": "5",
    "section": "A",
    "fatherName": "John Doe",
    "yearOfAdmission": "2022"
  },
  "degreeStats": {
    "creditsRequired": 160.0,
    "creditsEarned": 88.0,
    "creditsToEarn": 72.0,
    "cgpa": 8.74
  },
  "semesters": [
    {
      "semesterName": "Semester 4",
      "creditsRegistered": 22.0,
      "creditsEarned": 22.0,
      "sgpa": 8.91,
      "cgpa": 8.74,
      "courses": [
        {
          "code": "21CS41",
          "name": "Design & Analysis of Algorithms",
          "creditRegistered": "4.0",
          "creditEarned": "4.0",
          "gpa": 9.0,
          "grade": "A+"
        }
      ]
    }
  ]
}
```

---

## Errors & Status Codes

All errors return a standard JSON payload:

```json
{
  "error": "description of failure"
}
```

| HTTP Status | Reason |
| :--- | :--- |
| `400 Bad Request` | Missing required parameters (`usn`, `dob`, or `course`). |
| `401 Unauthorized` | Invalid student credentials or portal session rejected. |
| `404 Not Found` | Requested course or detail record not found for the student. |
| `500 Internal Server Error` | Portal connection failure or parsing error. |
