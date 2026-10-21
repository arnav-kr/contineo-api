package services

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/arnav-kr/contineo-api/models"
	"github.com/arnav-kr/contineo-api/parser"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func init() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Printf("[cache] REDIS_URL not set — caching disabled")
		return
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Printf("[cache] warning: failed to parse REDIS_URL %q: %v — caching disabled", redisURL, err)
		return
	}

	rdb = redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("[cache] warning: Redis unreachable (%v) — caching disabled", err)
		rdb = nil
	} else {
		log.Printf("[cache] connected to Redis")
	}
}

// NormalizePortal resolves a portal identifier to either "even" or "odd".
// Falls back to DEFAULT_PORTAL environment variable, or "even".
func NormalizePortal(portal string) string {
	p := strings.ToLower(strings.TrimSpace(portal))
	p = strings.Trim(p, "/")
	if p == "odd" || p == "parentsodd" {
		return "odd"
	}
	if p == "even" {
		return "even"
	}
	defaultPortal := strings.ToLower(strings.TrimSpace(os.Getenv("DEFAULT_PORTAL")))
	if defaultPortal == "odd" || defaultPortal == "parentsodd" {
		return "odd"
	}
	return "even"
}

// GetBaseURL returns the Contineo portal base URL for the given portal type.
func GetBaseURL(portal string) string {
	if NormalizePortal(portal) == "odd" {
		return "https://nmitparents.contineo.in/parentsodd"
	}
	return "https://nmitparents.contineo.in"
}

// cookieRecord is the JSON-serialisable representation of a single cookie.
type cookieRecord struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Domain string `json:"domain"`
	Path   string `json:"path"`
}

// loadJar restores cookies from Redis into a fresh jar.
// Returns a new empty jar on any error (cache miss, Redis down, etc.).
func loadJar(usn, dob, portal string) (http.CookieJar, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	if rdb == nil {
		return jar, nil
	}

	norm := NormalizePortal(portal)
	key := "cookies:" + norm + ":" + usn + ":" + dob
	ctx := context.Background()

	raw, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		// Miss or Redis error — return empty jar, not a fatal error.
		return jar, nil
	}

	var records []cookieRecord
	if err := json.Unmarshal(raw, &records); err != nil {
		return jar, nil
	}

	u, _ := url.Parse(GetBaseURL(norm))
	var cookies []*http.Cookie
	for _, r := range records {
		cookies = append(cookies, &http.Cookie{
			Name:   r.Name,
			Value:  r.Value,
			Domain: r.Domain,
			Path:   r.Path,
		})
	}
	jar.SetCookies(u, cookies)

	return jar, nil
}

// saveJar persists cookies from the jar into Redis with a 10-minute TTL.
func saveJar(usn, dob, portal string, jar http.CookieJar) {
	if rdb == nil {
		return
	}

	norm := NormalizePortal(portal)
	u, _ := url.Parse(GetBaseURL(norm))
	cookies := jar.Cookies(u)

	records := make([]cookieRecord, 0, len(cookies))
	for _, c := range cookies {
		records = append(records, cookieRecord{
			Name:   c.Name,
			Value:  c.Value,
			Domain: c.Domain,
			Path:   c.Path,
		})
	}

	raw, err := json.Marshal(records)
	if err != nil {
		log.Printf("[cache] warning: failed to marshal cookies: %v", err)
		return
	}

	ctx := context.Background()
	key := "cookies:" + norm + ":" + usn + ":" + dob
	if err := rdb.Set(ctx, key, raw, 10*time.Minute).Err(); err != nil {
		log.Printf("[cache] warning: failed to save cookies to Redis: %v", err)
	}
}

// cacheGet retrieves and JSON-unmarshals a cached value. Returns false on miss or error.
func cacheGet(key string, dest any) bool {
	if rdb == nil {
		return false
	}
	ctx := context.Background()
	raw, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		return false
	}
	return json.Unmarshal(raw, dest) == nil
}

// cacheSet JSON-marshals and stores a value in Redis with the given TTL.
func cacheSet(key string, value any, ttl time.Duration) {
	if rdb == nil {
		return
	}
	raw, err := json.Marshal(value)
	if err != nil {
		log.Printf("[cache] warning: failed to marshal value for key %q: %v", key, err)
		return
	}
	ctx := context.Background()
	if err := rdb.Set(ctx, key, raw, ttl).Err(); err != nil {
		log.Printf("[cache] warning: failed to set key %q in Redis: %v", key, err)
	}
}

func GetClientExported(usn, dob string, portal ...string) (*http.Client, error) {
	p := ""
	if len(portal) > 0 {
		p = portal[0]
	}
	jar, err := loadJar(usn, dob, p)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Jar:       jar,
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	return client, nil
}

func login(client *http.Client, usn, dob, portal string) (string, error) {
	norm := NormalizePortal(portal)
	baseURL := GetBaseURL(norm)

	data := url.Values{}
	data.Set("username", usn)
	data.Set("passwd", dob)
	data.Set("option", "com_user")
	data.Set("task", "login")

	req, err := http.NewRequest("POST", baseURL+"/index.php", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Referer", baseURL+"/index.php")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusForbidden {
		return "", errors.New("portal access forbidden (geo-blocked or blocked by firewall)")
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	html := string(bodyBytes)

	if !strings.Contains(html, "com_studentdashboard") && !strings.Contains(html, "cn-stu-data") {
		return "", errors.New("invalid credentials or portal is down")
	}

	saveJar(usn, dob, norm, client.Jar)
	return html, nil
}

// GetDashboard retrieves student overview and courses.
// If portal is omitted or empty, it automatically attempts primary portal and falls back to alternate portal.
func GetDashboard(usn, dob string, portal ...string) (*models.Overview, error) {
	requestedPortal := ""
	if len(portal) > 0 {
		requestedPortal = strings.TrimSpace(portal[0])
	}

	// Explicit portal request: target that portal directly
	if requestedPortal == "odd" || requestedPortal == "parentsodd" || requestedPortal == "even" {
		return getDashboardForPortal(usn, dob, requestedPortal)
	}

	// Auto-detection mode: try primary portal first. If credentials fail, try alternate portal.
	primary := NormalizePortal(requestedPortal)
	alternate := "odd"
	if primary == "odd" {
		alternate = "even"
	}

	overview, err := getDashboardForPortal(usn, dob, primary)
	if err == nil {
		return overview, nil
	}

	// Try the alternate portal if the primary returned invalid credentials / down
	if strings.Contains(err.Error(), "invalid credentials") || strings.Contains(err.Error(), "portal is down") {
		altOverview, altErr := getDashboardForPortal(usn, dob, alternate)
		if altErr == nil {
			return altOverview, nil
		}
	}

	return nil, err
}

func getDashboardForPortal(usn, dob, portal string) (*models.Overview, error) {
	norm := NormalizePortal(portal)
	cacheKey := "dashboard:" + norm + ":" + usn + ":" + dob

	var cached models.Overview
	if cacheGet(cacheKey, &cached) {
		return &cached, nil
	}

	client, err := GetClientExported(usn, dob, norm)
	if err != nil {
		return nil, err
	}

	baseURL := GetBaseURL(norm)
	params := url.Values{}
	params.Set("option", "com_studentdashboard")
	params.Set("format", "html")
	params.Set("t", fmt.Sprintf("%d", time.Now().UnixMilli()))

	req, err := http.NewRequest("GET", baseURL+"/index.php?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Referer", baseURL+"/index.php")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusForbidden {
		return nil, errors.New("portal access forbidden (geo-blocked or blocked by firewall)")
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	html := string(bodyBytes)

	if !strings.Contains(html, "cn-stu-data") {
		html, err = login(client, usn, dob, norm)
		if err != nil {
			return nil, err
		}
	}

	overview, err := parser.ParseOverview(html, baseURL)
	if err != nil {
		return nil, err
	}
	overview.Portal = norm

	cacheSet(cacheKey, overview, 5*time.Minute)

	return overview, nil
}

func GetAttendanceDetail(usn, dob string, detailUrl string, portal ...string) (*models.AttendanceData, error) {
	p := ""
	if len(portal) > 0 {
		p = portal[0]
	}
	if strings.Contains(detailUrl, "parentsodd") {
		p = "odd"
	}
	norm := NormalizePortal(p)
	baseURL := GetBaseURL(norm)

	cacheKey := "attendance:" + norm + ":" + usn + ":" + dob + ":" + detailUrl

	var cached models.AttendanceData
	if cacheGet(cacheKey, &cached) {
		return &cached, nil
	}

	client, err := GetClientExported(usn, dob, norm)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(detailUrl, "http") {
		detailUrl = baseURL + "/" + strings.TrimLeft(detailUrl, "/")
	}

	if !strings.Contains(detailUrl, "format=html") {
		detailUrl = fmt.Sprintf("%s&format=html&t=%d", detailUrl, time.Now().UnixMilli())
	}

	req, err := http.NewRequest("GET", detailUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Referer", baseURL+"/index.php?option=com_studentdashboard")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusForbidden {
		return nil, errors.New("portal access forbidden (geo-blocked or blocked by firewall)")
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	html := string(bodyBytes)

	if !strings.Contains(res.Request.URL.Host, "nmitparents.contineo.in") || !strings.Contains(html, "uk-table") {
		_, err = login(client, usn, dob, norm)
		if err != nil {
			return nil, err
		}

		req, err = http.NewRequest("GET", detailUrl, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")
		req.Header.Set("Referer", baseURL+"/index.php?option=com_studentdashboard")

		res, err = client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		bodyBytes, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		html = string(bodyBytes)
	}

	data, err := parser.ParseAttendanceDetail(html)
	if err != nil {
		return nil, err
	}

	cacheSet(cacheKey, data, 5*time.Minute)

	return data, nil
}

func GetCIEDetail(usn, dob string, detailUrl string, portal ...string) (*models.CIEData, error) {
	p := ""
	if len(portal) > 0 {
		p = portal[0]
	}
	if strings.Contains(detailUrl, "parentsodd") {
		p = "odd"
	}
	norm := NormalizePortal(p)
	baseURL := GetBaseURL(norm)

	cacheKey := "cie:" + norm + ":" + usn + ":" + dob + ":" + detailUrl

	var cached models.CIEData
	if cacheGet(cacheKey, &cached) {
		return &cached, nil
	}

	client, err := GetClientExported(usn, dob, norm)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(detailUrl, "http") {
		detailUrl = baseURL + "/" + strings.TrimLeft(detailUrl, "/")
	}

	if !strings.Contains(detailUrl, "format=html") {
		detailUrl = fmt.Sprintf("%s&format=html&t=%d", detailUrl, time.Now().UnixMilli())
	}

	req, err := http.NewRequest("GET", detailUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Referer", baseURL+"/index.php?option=com_studentdashboard")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusForbidden {
		return nil, errors.New("portal access forbidden (geo-blocked or blocked by firewall)")
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	html := string(bodyBytes)

	if !strings.Contains(res.Request.URL.Host, "nmitparents.contineo.in") || (!strings.Contains(html, "barPadding") && !strings.Contains(html, "CIE")) {
		_, err = login(client, usn, dob, norm)
		if err != nil {
			return nil, err
		}

		req, err = http.NewRequest("GET", detailUrl, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")
		req.Header.Set("Referer", baseURL+"/index.php?option=com_studentdashboard")

		res, err = client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		bodyBytes, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		html = string(bodyBytes)
	}

	data, err := parser.ParseCIEDetail(html)
	if err != nil {
		return nil, err
	}

	cacheSet(cacheKey, data, 5*time.Minute)

	return data, nil
}

func GetExamHistory(usn, dob string, portal ...string) (*models.ExamHistory, error) {
	overview, err := GetDashboard(usn, dob, portal...)
	if err != nil {
		return nil, err
	}

	norm := NormalizePortal(overview.Portal)
	baseURL := GetBaseURL(norm)
	cacheKey := "exam-history:" + norm + ":" + usn + ":" + dob

	var cached models.ExamHistory
	if cacheGet(cacheKey, &cached) {
		return &cached, nil
	}

	client, err := GetClientExported(usn, dob, norm)
	if err != nil {
		return nil, err
	}

	targetUrl := fmt.Sprintf("%s/index.php?option=com_history&task=getResult&usn=%s", baseURL, usn)
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Referer", baseURL+"/index.php?option=com_studentdashboard")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusForbidden {
		return nil, errors.New("portal access forbidden (geo-blocked or blocked by firewall)")
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	html := string(bodyBytes)

	if !strings.Contains(html, "Student's History") && !strings.Contains(html, "Student Name") {
		_, err = login(client, usn, dob, norm)
		if err != nil {
			return nil, err
		}

		req, _ = http.NewRequest("GET", targetUrl, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0")
		req.Header.Set("Referer", baseURL+"/index.php?option=com_studentdashboard")

		res, err = client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		bodyBytes, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		html = string(bodyBytes)
	}

	data, err := parser.ParseExamHistory(html)
	if err != nil {
		return nil, err
	}

	cacheSet(cacheKey, data, 5*time.Minute)

	return data, nil
}
