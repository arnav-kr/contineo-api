# Contineo API

A clean REST API for Contineo parent and student portals. Fetch attendance, CIE test marks, subject lists, and semester exam history as structured JSON without dealing with messy HTML tables or session logins.

- **Hosted URL:** [https://contineo-api.vercel.app](https://contineo-api.vercel.app)
- **Interactive Swagger Docs:** [https://contineo-api.vercel.app/docs](https://contineo-api.vercel.app/docs)
- **OpenAPI 3.0 Spec:** [https://contineo-api.vercel.app/openapi.json](https://contineo-api.vercel.app/openapi.json)
- **LLM Context:** [llms.txt](llms.txt) | [agent.md](agent.md)

---

## Quick Test (No credentials needed!)

Want to test the API or see what the responses look like? You don't need real college credentials. We have built-in mock student data:

- **USN:** `1NT24CS006`
- **DOB:** Any date (e.g. `2000-01-01`)
- **Demo Subjects:** `21CS41`, `21CS42`, `21CS43`

Try it straight in your browser:
- [View Dashboard](https://contineo-api.vercel.app/api/dashboard?usn=1NT24CS006&dob=2000-01-01)
- [View Attendance for 21CS41](https://contineo-api.vercel.app/api/attendance?usn=1NT24CS006&dob=2000-01-01&course=21CS41)
- [View CIE Marks for 21CS41](https://contineo-api.vercel.app/api/cie?usn=1NT24CS006&dob=2000-01-01&course=21CS41)
- [View Exam History & CGPA](https://contineo-api.vercel.app/api/exam-history?usn=1NT24CS006&dob=2000-01-01)

---

## Connect It to Your Favorite AI

Hook Contineo straight into your daily AI so you can just ask *"How's my attendance this week?"* or *"What did I score in my internal tests?"*.

### ChatGPT (Create a Custom GPT)
You can build your own personal college bot in 2 minutes:
1. Open ChatGPT, tap **Explore GPTs** in the sidebar, and click **+ Create**.
2. Head to the **Configure** tab and name it something like *College Buddy*.
3. Scroll down to **Actions** -> click **Create new action**.
4. Tap **Import from URL** and paste:
   ```
   https://contineo-api.vercel.app/openapi.json
   ```
5. In **Instructions**, add:
   ```text
   You help students check their college status via the Contineo API.
   Ask the user for their USN and DOB (format YYYY-MM-DD) if you don't know it yet.
   Call /api/dashboard first to see their enrolled courses and attendance percentages.
   When they ask about a specific subject's class logs or CIE scores, use /api/attendance or /api/cie with that course code.
   For past semesters or CGPA, check /api/exam-history.
   For demos without credentials, test USN "1NT24CS006" works with any DOB.
   ```
6. Save it for yourself, and you're good to go!

---

### Google Gemini (Skills & Gems)

#### Using Gemini Skills
If your Gemini account has **Skills** enabled:
1. Open [gemini.google.com](https://gemini.google.com) and go to **Skills**.
2. Create a new skill called `contineo` with trigger `/contineo`.
3. In instructions, tell it:
   ```text
   Contineo API Base: https://contineo-api.vercel.app
   - Dashboard: /api/dashboard?usn={usn}&dob={dob}
   - Attendance: /api/attendance?usn={usn}&dob={dob}&course={course}
   - CIE Marks: /api/cie?usn={usn}&dob={dob}&course={course}
   - Results: /api/exam-history?usn={usn}&dob={dob}
   Use YYYY-MM-DD for DOB. Test USN: 1NT24CS006 with any DOB.
   ```
4. Now in any chat, just type `/contineo` followed by your USN and DOB.

#### Using a Gem
Prefer a separate chat? Go to **Gem manager** -> **+ New Gem**, name it *Contineo Assistant*, and paste the same instructions.

---

### Claude (Projects)
1. Go to [claude.ai](https://claude.ai) and create a new **Project**.
2. In **Project Knowledge**, upload [`agent.md`](agent.md) or [`llms.txt`](llms.txt).
3. Under custom instructions, paste:
   ```text
   Use the Contineo API at https://contineo-api.vercel.app to look up college records.
   Endpoints include /api/dashboard, /api/attendance, /api/cie, and /api/exam-history.
   Always expect DOB as YYYY-MM-DD. Test USN: 1NT24CS006.
   ```

---

### Developer Coding Agents (Cursor, Windsurf, Antigravity)
- **Antigravity / Gemini CLI**: The repository already includes a pre-packaged skill in `.agents/skills/contineo-api-client`. It will be detected automatically.
- **Cursor / Windsurf**: Mention `@.agents/skills/contineo-api-client/SKILL.md` in chat or add it to `.cursorrules` / `.windsurfrules`.
- **Claude Code**: Append the skill to your `CLAUDE.md`.

---

## Endpoints Cheatsheet

All endpoints use `GET` and take student credentials as query parameters:

| Endpoint | Required Query Params | What You Get |
| :--- | :--- | :--- |
| `/api/dashboard` | `usn`, `dob`, `portal`? | Student details, course list, attendance averages & <75% warnings |
| `/api/attendance` | `usn`, `dob`, `course`, `portal`? | Full attendance log (dates, present/absent counts, percentage) |
| `/api/cie` | `usn`, `dob`, `course`, `portal`? | IA test marks, assignment scores, and final calculated CIE |
| `/api/exam-history` | `usn`, `dob`, `portal`? | Semester-by-semester grades, credit totals, SGPA, and CGPA |
| `/health` | *none* | Service status check |

> **Note on `portal`:** The API automatically checks both `even` and `odd` semesters, so you usually don't need to specify it unless you want to force one.

For full schema details and error responses, check [DOCS.md](DOCS.md) or [Swagger UI](https://contineo-api.vercel.app/docs).

---

## Running Locally

```bash
git clone https://github.com/arnav-kr/contineo-api.git
cd contineo-api

# Run tests
go test -v .

# Start API server on localhost:3000
go run main.go
```
