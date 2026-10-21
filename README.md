# Contineo API

Unofficial REST API for retrieving student academic data, course enrollments, attendance logs, CIE scores, and semester examination history from Contineo student portals.

Hosted deployment: **[https://contineo-api.vercel.app](https://contineo-api.vercel.app)**  
Interactive Swagger Docs: **[https://contineo-api.vercel.app/docs](https://contineo-api.vercel.app/docs)**  
OpenAPI 3.0 Spec: **[https://contineo-api.vercel.app/openapi.json](https://contineo-api.vercel.app/openapi.json)**  
LLM Reference: **[llms.txt](llms.txt)** | Agent Integration Guide: **[agent.md](agent.md)**

---

## Installing the Skill in AI Assistants

You can install and use the `contineo-api-client` skill across various AI coding assistants and agent frameworks.

### 1. Antigravity / Gemini CLI
The skill is located in `.agents/skills/contineo-api-client`.
- **Project-level**: Antigravity automatically discovers skills inside `.agents/skills/` when opening this repository.
- **Global installation**: Copy the skill directory to your global Antigravity config:
  ```bash
  # Linux / macOS
  mkdir -p ~/.gemini/config/skills
  cp -r .agents/skills/contineo-api-client ~/.gemini/config/skills/

  # Windows (PowerShell)
  New-Item -ItemType Directory -Force -Path "$HOME\.gemini\config\skills"
  Copy-Item -Recurse -Force ".agents\skills\contineo-api-client" "$HOME\.gemini\config\skills\"
  ```

### 2. Cursor
To use with Cursor's Agent and Composer:
- Add a project rule in `.cursorrules` or `.cursor/rules/contineo.mdc`:
  ```markdown
  Read and follow the guidelines in .agents/skills/contineo-api-client/SKILL.md whenever querying student attendance, CIE marks, or exam results via https://contineo-api.vercel.app.
  ```
- Alternatively, reference `@.agents/skills/contineo-api-client/SKILL.md` directly in Cursor chat.

### 3. Claude Code / Anthropic
Add the skill to Claude Code memory or system prompt:
- Copy the skill file into Claude's prompt or memory file:
  ```bash
  cat .agents/skills/contineo-api-client/SKILL.md >> CLAUDE.md
  ```

### 4. Windsurf / Codeium
Add the skill to `.windsurfrules`:
```markdown
When fetching student portal data or answering questions about student attendance, CIE, or exam grades, follow the workflow and schema defined in `.agents/skills/contineo-api-client/SKILL.md` using the API at `https://contineo-api.vercel.app`.
```

### 5. OpenAI Custom GPTs / Assistants
1. Navigate to **GPT Editor** -> **Configure** -> **Actions**.
2. Click **Import from URL** and paste:
   ```
   https://contineo-api.vercel.app/openapi.json
   ```
3. In **Instructions**, copy and paste the contents of [`.agents/skills/contineo-api-client/SKILL.md`](.agents/skills/contineo-api-client/SKILL.md) or [`agent.md`](agent.md).

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
