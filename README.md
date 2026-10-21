<div align="center">

# Contineo API

A high-performance Go REST API and scraping engine for Contineo student and parent portals. Normalizes unstructured portal HTML into clean, typed JSON payloads with serverless deployment on Vercel (`bom1` region), Redis session caching, and mock test mode.

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Vercel](https://img.shields.io/badge/Vercel-Deployed-black?style=for-the-badge&logo=vercel&logoColor=white)](https://contineo-api.vercel.app)
[![OpenAPI](https://img.shields.io/badge/OpenAPI-3.0.3-6BA539?style=for-the-badge&logo=openapiinitiative&logoColor=white)](https://contineo-api.vercel.app/openapi.json)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge)](LICENSE)

[Live API](https://contineo-api.vercel.app) • [Swagger UI](https://contineo-api.vercel.app/docs) • [OpenAPI Spec](https://contineo-api.vercel.app/openapi.json) • [Full Documentation](DOCS.md)

</div>

---

## Features

- **Automated Multi-Portal Routing**: Auto-resolves credentials between standard (`even`) and `parentsodd` (`odd`) portals with automatic fallback.
- **Strict Typed Schemas**: Clean Go models for student profiles, course registries, class attendance logs, CIE components, and degree credit audits.
- **Serverless & Edge Ready**: Deployed on Vercel Serverless Functions pinned to Mumbai, India (`bom1`) to prevent portal geo-blocking.
- **Built-in Mock Account**: Dedicated test student (`1NT24CS006`) returning deterministic test fixtures across all endpoints for testing without credentials.
- **Agent & AI Native**: Includes pre-built agent skill manifests (`.agents/skills/`), `llms.txt`, and `agent.md` for tool-calling agents and custom LLM actions.
- **Redis Session Caching**: Preserves session cookies and caches parsed responses with TTLs when `REDIS_URL` is configured.

---

## Architecture Flow

```
                      Client Request
             (usn, dob, course?, portal?)
                          │
                          ▼
             ┌─────────────────────────┐
             │  Contineo API (Fiber)   │
             │   Vercel Serverless     │
             └────────────┬────────────┘
                          │
         ┌────────────────┴────────────────┐
         ▼                                 ▼
   [Mock Student]                   [Live Student]
   USN: 1NT24CS006                         │
         │                         ┌───────┴───────┐
         ▼                         ▼               ▼
   Mock Fixture            Primary Portal    Fallback Portal
   (Instant JSON)             ("even")          ("odd")
                                   │               │
                                   └───────┬───────┘
                                           │
                                           ▼
                                 [Goquery HTML Parser]
                                           │
                                           ▼
                                  Typed JSON Response
```

---

## Quick Reference & Endpoints

Base URL: `https://contineo-api.vercel.app`

| Endpoint | Method | Parameters | Description |
| :--- | :--- | :--- | :--- |
| `/health` | `GET` | — | Service uptime and timestamp |
| `/api/dashboard` | `GET` | `usn`, `dob`, `portal`? | Student details, registered courses, attendance overview, warnings |
| `/api/attendance` | `GET` | `usn`, `dob`, `course`, `portal`? | Full attendance log and present/absent counts |
| `/api/cie` | `GET` | `usn`, `dob`, `course`, `portal`? | CIE component marks (IA tests, assignments) and final score |
| `/api/exam-history` | `GET` | `usn`, `dob`, `portal`? | Semester-by-semester grades, credit completion, SGPA, CGPA |

### Authentication & Parameters

- `usn` (string, required): Student University Seat Number.
- `dob` (string, required): Date of birth in `YYYY-MM-DD` format.
- `course` (string): Course code obtained from `/api/dashboard` (required for `/api/attendance` and `/api/cie`).
- `portal` (string, optional): Target portal override (`even` or `odd`). Omit to let the engine auto-resolve.

---

## Mock Mode (Testing)

Test API responses without real student credentials using the reserved mock USN:

- **USN:** `1NT24CS006`
- **DOB:** Any string (e.g. `2000-01-01`)
- **Demo Courses:** `21CS41`, `21CS42`, `21CS43`

```bash
# Dashboard
curl "https://contineo-api.vercel.app/api/dashboard?usn=1NT24CS006&dob=2000-01-01"

# Attendance for 21CS41
curl "https://contineo-api.vercel.app/api/attendance?usn=1NT24CS006&dob=2000-01-01&course=21CS41"

# CIE Marks for 21CS41
curl "https://contineo-api.vercel.app/api/cie?usn=1NT24CS006&dob=2000-01-01&course=21CS41"

# Exam History
curl "https://contineo-api.vercel.app/api/exam-history?usn=1NT24CS006&dob=2000-01-01"
```

---

## AI Assistant & Skill Integration

### 1. ChatGPT (Custom GPT Action)
1. In **GPT Editor** -> **Configure** -> **Actions** -> click **Create new action**.
2. Click **Import from URL** and paste:
   ```
   https://contineo-api.vercel.app/openapi.json
   ```
3. Set instructions to prompt for `usn` and `dob` (`YYYY-MM-DD`), fetch `/api/dashboard` first, then drill down into `/api/attendance` or `/api/cie` using the retrieved course codes.

### 2. Google Gemini (Skills & Gems)
- **Gemini Skill**: Create a skill with shortcut `/contineo` that points to `https://contineo-api.vercel.app/openapi.json` and uses the prompt template from [`agent.md`](agent.md).
- **Gem**: Create a Gem with the endpoints and parameter contracts defined in [`llms.txt`](llms.txt).

### 3. Claude (Projects)
Upload [`agent.md`](agent.md) or [`llms.txt`](llms.txt) to **Project Knowledge** on [claude.ai](https://claude.ai).

### 4. Coding Agents (Antigravity, Cursor, Windsurf)
- **Antigravity / Gemini CLI**: Auto-discovers the skill at [`.agents/skills/contineo-api-client/SKILL.md`](.agents/skills/contineo-api-client/SKILL.md).
- **Cursor / Windsurf**: Add a reference in `.cursorrules` or `.windsurfrules` pointing to `.agents/skills/contineo-api-client/SKILL.md`.

---

## Local Development

### Prerequisites
- Go `1.22+`

### Setup
```bash
# Clone repository
git clone https://github.com/arnav-kr/contineo-api.git
cd contineo-api

# Run test suite
go test -v .

# Run local development server
go run main.go
```

The server binds to port `3000` by default (`http://localhost:3000`).

---

## License

This project is licensed under the MIT License.
