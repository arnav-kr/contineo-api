# Contineo API Agent Integration Guide (agent.md)

Instructions and schema specifications for LLMs, autonomous agents, and tool-calling systems integrating with the Contineo API.

---

## 1. System Overview

The Contineo API wraps student portals into structured JSON endpoints.
- **Hosted Production Base URL:** `https://contineo-api.vercel.app`
- **API Base Path:** `/api` (e.g. `https://contineo-api.vercel.app/api/dashboard`)
- **Interactive Swagger UI:** `https://contineo-api.vercel.app/docs`
- **OpenAPI 3.0 Specification:** `https://contineo-api.vercel.app/openapi.json`
- **Health Checks:** `https://contineo-api.vercel.app/health` and `/api/health`
- **Response Format:** `application/json; charset=utf-8`
- **Authentication:** Stateless query parameters per request (`usn`, `dob`).

---

## 2. Parameter Contract

| Parameter | Type | Required | Allowed Values / Format | Description |
| :--- | :--- | :--- | :--- | :--- |
| `usn` | `string` | Yes | Match student USN format | University Seat Number. Use `1NT24CS006` for mock test data. |
| `dob` | `string` | Yes | `YYYY-MM-DD` | Date of Birth. Any string accepted when `usn=1NT24CS006`. |
| `portal` | `string` | No | `even`, `odd`, `parentsodd` | Portal to query. Auto-falls back to the alternate portal if omitted. |
| `course` | `string` | Endpoint-dependent | Extracted course `code` | Required for `/api/attendance` and `/api/cie` (e.g. `21CS41` for test student). |

### Test / Mock Student Mode
To test API response contracts without real student credentials:
- **USN:** `1NT24CS006`
- **DOB:** Any value (e.g. `2000-01-01`)
- **Courses available:** `21CS41`, `21CS42`, `21CS43`
All endpoints (`/api/dashboard`, `/api/attendance`, `/api/cie`, `/api/exam-history`) return realistic mock responses for this account.

---

## 3. Tool-Calling Agent Workflow

When answering user queries about student status, attendance, marks, or degree progress, follow this execution flow:

```
                  ┌──────────────────────┐
                  │ Receive User Request │
                  └──────────┬───────────┘
                             │
                             ▼
                 [Query /api/dashboard]
                             │
          ┌──────────────────┴──────────────────┐
          │                                     │
          ▼                                     ▼
 [High-level summary needed?]           [Specific course detail needed?]
          │                                     │
          ▼                                     ▼
   Return dashboard data          Extract `course` code from `courses[]`
                                                │
                                  ┌─────────────┴─────────────┐
                                  ▼                           ▼
                       [/api/attendance]                 [/api/cie]
```

1. **Always fetch `/api/dashboard` first**:
   - Validates student credentials (`usn`, `dob`).
   - Retrieves active enrolled courses and their unique `code` values (e.g. `21CS51`).
   - Returns overall attendance averages and warning lists under `summary.belowThreshold`.
   - Returns `portal` indicating which portal resolved successfully.
2. **Retrieve Drill-Down Data**:
   - For daily/session attendance breakdown, call `GET /api/attendance` passing the `course` code obtained from the dashboard.
   - For detailed test/quiz/assignment splits, call `GET /api/cie` passing the `course` code.
3. **Retrieve Historical & Cumulative Stats**:
   - For previous semesters, credits tally, SGPA, or CGPA, call `GET /api/exam-history`.

---

## 4. Endpoints & Schemas

### `GET /api/dashboard`
Retrieves student profile and enrolled courses overview.

**Query Parameters:**
- `usn` (string, required)
- `dob` (string, required)
- `portal` (string, optional)

**Response Shape:**
```json
{
  "portal": "even",
  "student": {
    "name": "string",
    "usn": "string",
    "program": "string | null",
    "semester": "string | null",
    "section": "string | null",
    "fatherName": "string | null",
    "yearOfAdmission": "string | null"
  },
  "courses": [
    {
      "code": "string",
      "name": "string",
      "ids": { "courseId": "string", "secId": "string", "semId": "string" },
      "attendance": { "percentage": 0.0, "detailUrl": "string | null" },
      "cie": { "marks": 0.0, "detailUrl": "string | null" },
      "lessonNotesUrl": "string | null"
    }
  ],
  "summary": {
    "totalCourses": 0,
    "averageAttendance": 0.0,
    "belowThreshold": {
      "85": ["string"],
      "75": ["string"]
    }
  }
}
```

---

### `GET /api/attendance`
Retrieves per-class attendance log and aggregate counts for one course.

**Query Parameters:**
- `usn` (string, required)
- `dob` (string, required)
- `course` (string, required)
- `portal` (string, optional)

**Response Shape:**
```json
{
  "stats": {
    "totalClasses": 0,
    "present": 0,
    "absent": 0,
    "percentage": 0.0
  },
  "records": [
    {
      "date": "YYYY-MM-DD",
      "startTime": "HH:MM",
      "endTime": "HH:MM",
      "isPresent": true,
      "status": "Present | Absent"
    }
  ]
}
```

---

### `GET /api/cie`
Retrieves breakdown of internal assessments and components for a course.

**Query Parameters:**
- `usn` (string, required)
- `dob` (string, required)
- `course` (string, required)
- `portal` (string, optional)

**Response Shape:**
```json
{
  "components": [
    {
      "name": "string",
      "marks": {
        "obtained": 0.0,
        "max": 0.0,
        "raw": "string"
      }
    }
  ],
  "finalCIE": {
    "obtained": 0.0,
    "max": 0.0,
    "raw": "string"
  },
  "attendancePercentage": 0.0
}
```

---

### `GET /api/exam-history`
Retrieves historical semester performance, SGPA, and cumulative degree credits.

**Query Parameters:**
- `usn` (string, required)
- `dob` (string, required)
- `portal` (string, optional)

**Response Shape:**
```json
{
  "lastUpdated": "YYYY-MM-DD",
  "student": {
    "name": "string",
    "usn": "string",
    "program": "string",
    "semester": "string",
    "section": "string",
    "fatherName": "string",
    "yearOfAdmission": "string"
  },
  "degreeStats": {
    "creditsRequired": 0.0,
    "creditsEarned": 0.0,
    "creditsToEarn": 0.0,
    "cgpa": 0.0
  },
  "semesters": [
    {
      "semesterName": "string",
      "creditsRegistered": 0.0,
      "creditsEarned": 0.0,
      "sgpa": 0.0,
      "cgpa": 0.0,
      "courses": [
        {
          "code": "string",
          "name": "string",
          "creditRegistered": "string",
          "creditEarned": "string",
          "gpa": 0.0,
          "grade": "string"
        }
      ]
    }
  ]
}
```

---

## 5. Error Handling & Edge Cases

All failure responses use standard HTTP status codes with body format:
```json
{
  "error": "message details"
}
```

| HTTP Status | Trigger | Agent Remediation |
| :--- | :--- | :--- |
| `400` | Missing query params | Prompt user for `usn`, `dob` (`YYYY-MM-DD`), or provide required `course` code. |
| `401` | Invalid credentials / portal unreachable | Inform user credentials failed or portal is down; verify `dob` format. |
| `404` | Course code not found in dashboard | Verify `course` code against `courses[].code` from `/api/dashboard`. |
| `500` | Upstream portal parse error | Retry or report temporary upstream issue. |
