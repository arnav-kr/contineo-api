# Contineo API Documentation

Overview
- Base path: `/api`
- Auth: most endpoints require `usn` and `dob` query parameters. `/api/attendance` and `/api/cie` also require `course`. The service uses the portal login flow and cookie jar to maintain session state; sessions expire after a few minutes.

Global behavior
- `usn` is the student USN.
- `dob` must use `YYYY-MM-DD`.
- Errors are returned as JSON: `{ "error": "message" }`.
- The root health check is `GET /` and returns `Contineo API is running`.

Endpoints

- `GET /`
  - Description: Health check / service status.
  - Response: plain text.

- `GET /api/dashboard?usn={usn}&dob={dob}`
  - Description: Returns the student's overview including personal details and course list.
  - Query params: `usn` (string), `dob` (YYYY-MM-DD string)
  - Response model: `models.Overview`
  - Response fields: `student`, `courses[]`, `summary`.
  - Errors: 400 if params missing, 401 on auth/login failure.

- `GET /api/attendance?usn={usn}&dob={dob}&course={code}`
  - Description: Returns attendance detail for the provided course code.
  - Query params: `usn`, `dob`, `course` (course code from `dashboard` response)
  - Response model: `models.AttendanceData`
  - Response fields: `stats`, `records[]`.
  - Errors: 404 if attendance detail URL not available, 401/500 on fetch failure.

- `GET /api/cie?usn={usn}&dob={dob}&course={code}`
  - Description: Returns parsed CIE (continuous internal evaluation) details for a course.
  - Query params: `usn`, `dob`, `course` (course code from `dashboard` response)
  - Response model: `models.CIEData`
  - Response fields: `components[]`, `finalCIE`, `attendancePercentage`.
  - Errors: 404 if CIE detail URL not available, 401/500 on fetch failure.

- `GET /api/exam-history?usn={usn}&dob={dob}`
  - Description: Returns parsed exam history (semester-wise) and summary statistics.
  - Response model: `models.ExamHistory`
  - Response highlights:
    - `lastUpdated` — ISO date string `YYYY-MM-DD` (converted from portal `DD/MM/YYYY`).
    - `student` — fields: `name`, `usn`, `program`, `semester`, `section`, `fatherName`, `yearOfAdmission`.
    - `degreeStats` — `creditsRequired`, `creditsEarned`, `creditsToEarn`, `cgpa`.
    - `semesters[]` — each contains `semesterName`, `creditsRegistered`, `creditsEarned`, `sgpa`, `cgpa`, and `courses[]`.

Examples

Curl exam-history:

```
curl 'http://localhost:3000/api/exam-history?usn=1XX00XX000&dob=2000-01-01'
```

Error handling
- Returns JSON `{ "error": "message" }` with appropriate HTTP status codes.
- Missing query parameters return 400.
- Auth/session failures typically return 401 or 500 depending on where the failure occurs.

Notes / Implementation details
- The service uses an in-memory cookie cache per `usn:dob` to reuse session cookies for a short window.
- If the portal redirects (session expired) the service transparently re-authenticates and retries once.
- Date normalization: `Last Updated On: DD/MM/YYYY` is converted to `YYYY-MM-DD`.

Data models
- See `api/models/models.go` for the full JSON schema used by the endpoints.
- `models.Student`: `name`, `usn`, `program`, `semester`, `section`, `fatherName`, `yearOfAdmission`
- `models.Course`: `code`, `name`, `ids`, `attendance`, `cie`, `lessonNotesUrl`
- `models.AttendanceData`: `stats`, `records`
- `models.CIEData`: `components`, `finalCIE`, `attendancePercentage`
- `models.ExamHistory`: `lastUpdated`, `student`, `degreeStats`, `semesters`

Contact
- Repo: arnav-kr/contineo
