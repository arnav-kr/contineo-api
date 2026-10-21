---
name: contineo-api-client
description: >-
  Interact with the Contineo student portal API to fetch student dashboard data,
  course enrollments, attendance records, continuous internal evaluation (CIE) marks,
  and semester examination history. Use whenever asked to check grades, attendance,
  exam results, or profile information from Contineo portals.
---

# Contineo API Skill

Guide for LLMs and AI agents to query the Contineo student portal API.

## 1. Host & Endpoints Overview

- **Hosted Production Base URL:** `https://contineo-api.vercel.app`
- **Interactive Swagger UI:** `https://contineo-api.vercel.app/docs`
- **OpenAPI 3.0 Schema:** `https://contineo-api.vercel.app/openapi.json`
- **Health Checks:** `https://contineo-api.vercel.app/health` and `/api/health`

Agents querying this service should prepend `https://contineo-api.vercel.app` to all request paths.

---

## 2. Parameters & Authentication

Endpoints authenticate via HTTP GET query parameters:

| Parameter | Type | Required | Description | Example |
| :--- | :--- | :--- | :--- | :--- |
| `usn` | string | Yes | University Seat Number | `1NT24CS006` (mock) or student USN |
| `dob` | string | Yes | Date of birth in `YYYY-MM-DD` | `2000-01-01` |
| `portal` | string | No | `even` or `odd` (`parentsodd`). Auto-detects if omitted. | `even` |
| `course` | string | Endpoint-specific | Course code from `/api/dashboard` | `21CS41` |

> **Mock Mode**: Use `usn=1NT24CS006` with any `dob` for testing response schemas without real credentials.

---

## 2. Standard Query Workflow

Always follow this sequence when interacting with the API:

```
1. Call https://contineo-api.vercel.app/api/dashboard?usn={usn}&dob={dob}
   ├── Inspect profile details and enrolled courses
   ├── Check `summary.belowThreshold` for attendance warnings (<75% and <85%)
   └── Note `portal` ("even" or "odd") returned in payload

2. If specific details requested:
   ├── Attendance breakdown: Call https://contineo-api.vercel.app/api/attendance?usn={usn}&dob={dob}&course={courseCode}&portal={portal}
   ├── CIE / internal marks: Call https://contineo-api.vercel.app/api/cie?usn={usn}&dob={dob}&course={courseCode}&portal={portal}
   └── Cumulative degree & SGPA: Call https://contineo-api.vercel.app/api/exam-history?usn={usn}&dob={dob}&portal={portal}
```

---

## 3. Endpoints

### 3.1 Overview (`/api/dashboard`)
```http
GET https://contineo-api.vercel.app/api/dashboard?usn={usn}&dob={dob}&portal={even|odd}
```
Returns student profile metadata, enrolled courses list, attendance percentages, CIE totals, and aggregate statistics.

### 3.2 Attendance Details (`/api/attendance`)
```http
GET https://contineo-api.vercel.app/api/attendance?usn={usn}&dob={dob}&course={courseCode}&portal={even|odd}
```
Returns class attendance breakdown:
- `stats`: `totalClasses`, `present`, `absent`, `percentage`
- `records`: List of class sessions with `date`, `startTime`, `endTime`, `isPresent`, `status`.

### 3.3 Internal Evaluation / CIE (`/api/cie`)
```http
GET https://contineo-api.vercel.app/api/cie?usn={usn}&dob={dob}&course={courseCode}&portal={even|odd}
```
Returns:
- `components`: Breakdown of tests, assignments, quizzes (`obtained`, `max`, `raw`).
- `finalCIE`: Final aggregate CIE score.
- `attendancePercentage`: Associated attendance percentage.

### 3.4 Examination History (`/api/exam-history`)
```http
GET https://contineo-api.vercel.app/api/exam-history?usn={usn}&dob={dob}&portal={even|odd}
```
Returns:
- `degreeStats`: Total credits required, earned, remaining, and cumulative CGPA.
- `semesters`: Semester results including SGPA, credits earned, and per-course grades and GPAs.

---

## 4. Error Handling

Errors return `{"error": "message"}`:
- `400 Bad Request`: Missing `usn`, `dob`, or `course`. Ask user for missing parameters.
- `401 Unauthorized`: Invalid credentials or portal blocked/down. Verify USN and `YYYY-MM-DD` DOB format.
- `404 Not Found`: Course code or records not found for the student.
- `500 Internal Error`: Upstream portal issue or parsing error.
