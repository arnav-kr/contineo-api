# Deploying Contineo API to Vercel

The Contineo Go API is fully configured for deployment on Vercel using the official **Go Framework Preset** and serverless execution.

---

## Deployment Options

### Method 1: Vercel Dashboard (GitHub / GitLab / Bitbucket)

Zero configuration required:

1. Import the repository in [Vercel](https://vercel.com/new).
2. Vercel automatically detects the **Go** Framework Preset from the root `go.mod` and `main.go`.
3. Click **Deploy**.

---

### Method 2: Vercel CLI

Deploy directly from your terminal:

```bash
# Preview deployment:
vercel

# Production deployment:
vercel --prod
```

---

## Configuration (`vercel.json`)

```json
{
  "$schema": "https://openapi.vercel.sh/vercel.json",
  "framework": "go",
  "build": {
    "env": {
      "GO_BUILD_FLAGS": "-ldflags '-s -w'"
    }
  }
}
```

- Strips debugging symbols with `-ldflags '-s -w'` for faster cold starts and minimal binary footprint.

---

## Optional Environment Variables

| Variable | Description | Default |
| :--- | :--- | :--- |
| `PORT` | Managed automatically by Vercel. | `3000` |
| `REDIS_URL` | Distributed cache URL (e.g., [Upstash Redis](https://upstash.com)). If omitted, Redis caching is gracefully skipped without startup delays. | `""` |
| `DB_HOST` | PostgreSQL host. If omitted, database storage is safely disabled. | `""` |
| `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE` | PostgreSQL connection parameters (used only if `DB_HOST` is set). | — |

---

## Endpoints Available Once Deployed

- `GET /` — API running status
- `GET /health` or `GET /api/health` — Health check endpoint
- `GET /docs` — Swagger UI API documentation
- `GET /openapi.json` — OpenAPI 3.0 specification
- `GET /api/dashboard?usn={usn}&dob={dob}` — Student overview & courses
- `GET /api/attendance?usn={usn}&dob={dob}&course={code}` — Attendance details
- `GET /api/cie?usn={usn}&dob={dob}&course={code}` — Continuous Internal Evaluation details
- `GET /api/exam-history?usn={usn}&dob={dob}` — Semester exam records & SGPA/CGPA
