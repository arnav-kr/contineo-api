# Contineo API

Unofficial REST API for retrieving student academic data, course enrollments, attendance logs, CIE scores, and semester examination history from Contineo student portals.

Hosted deployment: **[https://contineo-api.vercel.app](https://contineo-api.vercel.app)**  
Interactive Swagger Docs: **[https://contineo-api.vercel.app/docs](https://contineo-api.vercel.app/docs)**  
OpenAPI 3.0 Spec: **[https://contineo-api.vercel.app/openapi.json](https://contineo-api.vercel.app/openapi.json)**  
LLM Reference: **[llms.txt](llms.txt)** | Agent Integration Guide: **[agent.md](agent.md)**

---

## Using with AI Assistants (ChatGPT, Claude, Gemini, etc.)

You can give everyday AI chat assistants the ability to check your Contineo attendance, marks, and academic status.

### 1. ChatGPT (Custom GPT)
Create your own student assistant GPT that fetches real-time data:
1. Go to **ChatGPT** -> Click on your profile or sidebar -> **Explore GPTs** -> **+ Create**.
2. Go to the **Configure** tab:
   - **Name:** Contineo Student Assistant
   - **Description:** Checks attendance, CIE marks, and exam history from Contineo.
3. Scroll down to **Actions** -> Click **Create new action**.
4. In **Schema**, click **Import from URL** and paste:
   ```
   https://contineo-api.vercel.app/openapi.json
   ```
5. In **Instructions**, copy and paste:
   ```text
   You are a student assistant connected to the Contineo API.
   When asked for attendance, marks, or grades, prompt the user for their USN and DOB (format YYYY-MM-DD) if not provided.
   First call GET /api/dashboard with the usn and dob to fetch student courses and overview.
   If the user asks for detailed attendance or CIE marks for a specific subject, find the matching course code from the dashboard response and call /api/attendance or /api/cie with that course code.
   For overall credits, SGPA, and semester grades, call /api/exam-history.
   For testing purposes without real student credentials, you can use USN "1NT24CS006" with any DOB.
   ```
6. Click **Save** / **Publish** (select "Only me" for private use).

---

### 2. Claude (Projects / Custom Instructions)
1. Go to [claude.ai](https://claude.ai) -> Open or create a **Project** (or use Custom Instructions).
2. Under **Project Knowledge** (or chat attachment), upload the [`agent.md`](agent.md) or [`llms.txt`](llms.txt) file.
3. In **Custom Instructions for this project**, paste:
   ```text
   You have access to the Contineo API hosted at https://contineo-api.vercel.app.
   API endpoints:
   - GET https://contineo-api.vercel.app/api/dashboard?usn={usn}&dob={dob}
   - GET https://contineo-api.vercel.app/api/attendance?usn={usn}&dob={dob}&course={courseCode}
   - GET https://contineo-api.vercel.app/api/cie?usn={usn}&dob={dob}&course={courseCode}
   - GET https://contineo-api.vercel.app/api/exam-history?usn={usn}&dob={dob}
   When I ask about attendance or grades, generate the exact URL for me to open or use web search/artifacts if enabled.
   Mock demo USN for testing: 1NT24CS006 (any DOB).
   ```

---

### 3. Google Gemini (Gems)
1. Go to [gemini.google.com](https://gemini.google.com) -> In the sidebar click **Gem manager** -> **+ New Gem**.
2. **Name:** Contineo Assistant
3. In **Instructions**, paste:
   ```text
   You are an assistant for college students using Contineo.
   The REST API is hosted at https://contineo-api.vercel.app.
   Endpoints:
   - Dashboard & enrolled courses: https://contineo-api.vercel.app/api/dashboard?usn={usn}&dob={dob}
   - Attendance logs: https://contineo-api.vercel.app/api/attendance?usn={usn}&dob={dob}&course={course}
   - CIE marks: https://contineo-api.vercel.app/api/cie?usn={usn}&dob={dob}&course={course}
   - Exam history: https://contineo-api.vercel.app/api/exam-history?usn={usn}&dob={dob}
   Always verify the date of birth is in YYYY-MM-DD format. Demo account for testing: USN 1NT24CS006 with any DOB.
   ```
4. Save the Gem.

---

### 4. Developer / Agent Skills (Cursor, Antigravity, Windsurf)

If you are using AI coding agents:
- **Antigravity / Gemini CLI:** Skills inside `.agents/skills/contineo-api-client` are auto-discovered, or copy to `~/.gemini/config/skills/`.
- **Cursor:** Mention `@.agents/skills/contineo-api-client/SKILL.md` or add a rule in `.cursorrules`.
- **Windsurf:** Add a rule in `.windsurfrules` pointing to `.agents/skills/contineo-api-client/SKILL.md`.
- **Claude Code:** Add to `CLAUDE.md`.

---

## Mock Test Account

For testing API integrations without real student credentials:
- **USN:** `1NT24CS006`
- **DOB:** Any string (e.g. `2000-01-01`)
- **Enrolled Courses:** `21CS41`, `21CS42`, `21CS43`

All endpoints return simulated payloads for this USN.

---

## API Endpoints

Base URL: `https://contineo-api.vercel.app`

| Method | Endpoint | Query Parameters | Description |
| :--- | :--- | :--- | :--- |
| `GET` | `/health` | None | Service health status |
| `GET` | `/api/dashboard` | `usn`, `dob`, `portal`? | Student profile, course attendance summary, active portal |
| `GET` | `/api/attendance` | `usn`, `dob`, `course`, `portal`? | Session-by-session class attendance logs |
| `GET` | `/api/cie` | `usn`, `dob`, `course`, `portal`? | Internal assessment and quiz mark breakdown |
| `GET` | `/api/exam-history` | `usn`, `dob`, `portal`? | Historical semester grades, SGPA, credits, CGPA |

Detailed schema documentation is available in [DOCS.md](DOCS.md).

---

## Local Development

```bash
# Clone the repository
git clone https://github.com/arnav-kr/contineo-api.git
cd contineo-api

# Run unit tests
go test -v .

# Start local server (port 3000)
go run main.go
```
