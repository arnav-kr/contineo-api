package main

import (
	_ "embed"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/arnav-kr/contineo-api/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/storage/postgres/v3"
)

//go:embed openapi.json
var openapiSpec []byte

var (
	app     *fiber.App
	appOnce sync.Once
)

func getApp() *fiber.App {
	appOnce.Do(func() {
		app = buildApp()
	})
	return app
}

func buildApp() *fiber.App {
	instance := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	store := initPostgresStorage()
	if store != nil {
		defer store.Close()
	}

	instance.Use(cors.New())

	instance.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Contineo API is running")
	})

	instance.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	instance.Get("/openapi.json", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Send(openapiSpec)
	})

	instance.Get("/docs", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Contineo API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  <style>body{margin:0}.swagger-ui .topbar{display:none}</style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    SwaggerUIBundle({
      url: "/openapi.json",
      dom_id: "#swagger-ui",
      presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
      layout: "BaseLayout",
      deepLinking: true,
    });
  </script>
</body>
</html>`)
	})

	api := instance.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	api.Get("/dashboard", func(c *fiber.Ctx) error {
		usn := c.Query("usn")
		dob := c.Query("dob")

		if usn == "" || dob == "" {
			return fiber.NewError(fiber.StatusBadRequest, "usn and dob query parameters are required")
		}

		overview, err := services.GetDashboard(usn, dob)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		return c.JSON(overview)
	})

	api.Get("/attendance", func(c *fiber.Ctx) error {
		usn := c.Query("usn")
		dob := c.Query("dob")
		courseCode := c.Query("course")

		if usn == "" || dob == "" || courseCode == "" {
			return fiber.NewError(fiber.StatusBadRequest, "usn, dob, and course query parameters are required")
		}

		overview, err := services.GetDashboard(usn, dob)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		var detailUrl string
		for _, c := range overview.Courses {
			if c.Code == courseCode {
				if c.Attendance.DetailURL != nil {
					detailUrl = *c.Attendance.DetailURL
				}
				break
			}
		}

		if detailUrl == "" {
			return fiber.NewError(fiber.StatusNotFound, "Attendance details not available for this course")
		}

		detail, err := services.GetAttendanceDetail(usn, dob, detailUrl)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(detail)
	})

	api.Get("/cie", func(c *fiber.Ctx) error {
		usn := c.Query("usn")
		dob := c.Query("dob")
		courseCode := c.Query("course")

		if usn == "" || dob == "" || courseCode == "" {
			return fiber.NewError(fiber.StatusBadRequest, "usn, dob, and course query parameters are required")
		}

		overview, err := services.GetDashboard(usn, dob)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		var detailUrl string
		for _, c := range overview.Courses {
			if c.Code == courseCode {
				if c.CIE.DetailURL != nil {
					detailUrl = *c.CIE.DetailURL
				}
				break
			}
		}

		if detailUrl == "" {
			return fiber.NewError(fiber.StatusNotFound, "CIE details not available for this course")
		}

		detail, err := services.GetCIEDetail(usn, dob, detailUrl)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(detail)
	})

	api.Get("/exam-history", func(c *fiber.Ctx) error {
		usn := c.Query("usn")
		dob := c.Query("dob")

		if usn == "" || dob == "" {
			return fiber.NewError(fiber.StatusBadRequest, "usn and dob query parameters are required")
		}

		history, err := services.GetExamHistory(usn, dob)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(history)
	})

	return instance
}

// Handler is exported for Vercel Serverless Function compatibility.
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "" {
		r.RequestURI = r.URL.RequestURI()
	}
	adaptor.FiberApp(getApp())(w, r)
}

func main() {
	server := getApp()
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Starting Contineo API on port %s", port)
	log.Fatal(server.Listen(":" + port))
}

func initPostgresStorage() *postgres.Storage {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		log.Printf("[db] DB_HOST not set — Postgres storage disabled")
		return nil
	}

	port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

	store := postgres.New(postgres.Config{
		Host:       dbHost,
		Port:       port,
		Username:   getEnv("DB_USER", "postgres"),
		Password:   getEnv("DB_PASSWORD", "postgres"),
		Database:   getEnv("DB_NAME", "contineo_db"),
		Table:      "fiber_storage",
		SSLMode:    getEnv("DB_SSLMODE", "disable"),
		Reset:      false,
		GCInterval: 10 * time.Second,
	})

	return store
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
