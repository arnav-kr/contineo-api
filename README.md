# Contineo API

Unofficial Go REST API and scraper for Contineo parent and student portals (attendance, CIE internal marks, registered courses, and exam history).

Hosted at **[https://contineo-api.vercel.app](https://contineo-api.vercel.app)**

- **Swagger Docs:** [https://contineo-api.vercel.app/docs](https://contineo-api.vercel.app/docs)
- **OpenAPI 3.0 Spec:** [https://contineo-api.vercel.app/openapi.json](https://contineo-api.vercel.app/openapi.json)
- **Docs & Schemas:** [DOCS.md](DOCS.md) | [agent.md](agent.md) | [llms.txt](llms.txt)

---

## Endpoints

Base URL: `https://contineo-api.vercel.app`

| Endpoint | Method | Parameters | Description |
| :--- | :--- | :--- | :--- |
| `/api/dashboard` | `GET` | `usn`, `dob`, `portal`? | Student details, registered courses, attendance overview |
| `/api/attendance` | `GET` | `usn`, `dob`, `course`, `portal`? | Subject attendance log and stats |
| `/api/cie` | `GET` | `usn`, `dob`, `course`, `portal`? | Internal assessment and quiz marks breakdown |
| `/api/exam-history` | `GET` | `usn`, `dob`, `portal`? | Semester-by-semester grades, SGPA, credits, CGPA |
| `/health` | `GET` | — | Health check |

All requests take `usn` and `dob` (`YYYY-MM-DD`). The engine automatically checks both `even` and `odd` (`parentsodd`) portals with fallback if `portal` is not specified.

---

## Test Account (Mock Mode)

To test responses without real credentials, use the built-in test student:
- **USN:** `1NT24CS006`
- **DOB:** Any value (e.g. `2000-01-01`)
- **Demo Courses:** `21CS41`, `21CS42`, `21CS43`

```bash
# Dashboard
curl "https://contineo-api.vercel.app/api/dashboard?usn=1NT24CS006&dob=2000-01-01"

# Attendance
curl "https://contineo-api.vercel.app/api/attendance?usn=1NT24CS006&dob=2000-01-01&course=21CS41"

# CIE Marks
curl "https://contineo-api.vercel.app/api/cie?usn=1NT24CS006&dob=2000-01-01&course=21CS41"

# Exam History
curl "https://contineo-api.vercel.app/api/exam-history?usn=1NT24CS006&dob=2000-01-01"
```

---

## Using with AI

### ChatGPT (Custom GPT)
1. In **GPT Editor** -> **Configure** -> **Actions** -> **Create new action**.
2. **Import from URL**:
   ```
   https://contineo-api.vercel.app/openapi.json
   ```
3. Set instructions to take the student's `usn` and `dob` (`YYYY-MM-DD`), query `/api/dashboard` first to extract course codes, and then call `/api/attendance` or `/api/cie` for details.

### Google Gemini (Skills & Gems)
- **Skills**: Create a `/contineo` skill using the schema and workflow from [agent.md](agent.md).
- **Gems**: Paste the endpoint reference and mock USN (`1NT24CS006`) into the Gem instructions.

### Claude (Skills & Projects)
- **Claude Skills**: In Claude.ai, go to **Settings** -> **Capabilities** -> **Skills** (or **Customize** -> **Skills**) and upload the `.agents/skills/contineo-api-client` folder (or zip). This also syncs with Claude Code CLI (`/skills`).
- **Claude Projects**: Alternatively, attach [agent.md](agent.md) or [llms.txt](llms.txt) to **Project Knowledge**.

### Coding Agents (Cursor / Windsurf / Antigravity)
The repository includes an agent skill manifest in `.agents/skills/contineo-api-client/SKILL.md`. In Cursor/Windsurf, reference `@.agents/skills/contineo-api-client/SKILL.md`.

---

## Local Setup

```bash
git clone https://github.com/arnav-kr/contineo-api.git
cd contineo-api

go test -v .
go run main.go
```

Runs on `http://localhost:3000`.

---

## License

MIT
