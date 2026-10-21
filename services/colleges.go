package services

import (
	"net/url"
	"os"
	"strings"
)

// CollegeConfig defines domain settings for an institution's Contineo portal.
type CollegeConfig struct {
	Code       string // e.g. "nmit"
	USNPrefix  string // e.g. "1NT"
	BaseHost   string // e.g. "nmitparents.contineo.in"
	OddSubpath string // e.g. "parentsodd" (default)
}

// DefaultBaseHost is the fallback Contineo host when unconfigured or unknown.
var DefaultBaseHost = "nmitparents.contineo.in"

func init() {
	if host := os.Getenv("CONTINEO_BASE_HOST"); host != "" {
		DefaultBaseHost = host
	}
}

// Registry maps USN prefixes and college codes to their Contineo portal domains.
var collegeRegistry = map[string]CollegeConfig{
	"1nt": {
		Code:       "nmit",
		USNPrefix:  "1NT",
		BaseHost:   "nmitparents.contineo.in",
		OddSubpath: "parentsodd",
	},
	"nmit": {
		Code:       "nmit",
		USNPrefix:  "1NT",
		BaseHost:   "nmitparents.contineo.in",
		OddSubpath: "parentsodd",
	},
	"sjec": {
		Code:       "sjec",
		USNPrefix:  "4SO",
		BaseHost:   "sjecparents.contineo.in",
		OddSubpath: "parentsodd",
	},
	"jnmc": {
		Code:       "jnmc",
		USNPrefix:  "JNMC",
		BaseHost:   "jnmcparents.contineo.in",
		OddSubpath: "parentsodd",
	},
	"wit": {
		Code:       "wit",
		USNPrefix:  "WIT",
		BaseHost:   "wit.contineo.in",
		OddSubpath: "parentsodd",
	},
}

// ResolveCollegeHost determines the appropriate Contineo base domain from USN,
// custom college param, or environment variable CONTINEO_BASE_HOST.
func ResolveCollegeHost(usn string, collegeOverride ...string) string {
	// 1. Explicit override parameter (e.g., college=sjec or college=https://sjecparents.contineo.in)
	if len(collegeOverride) > 0 && strings.TrimSpace(collegeOverride[0]) != "" {
		raw := strings.TrimSpace(collegeOverride[0])
		if strings.Contains(raw, "://") {
			if u, err := url.Parse(raw); err == nil && u.Host != "" {
				return u.Host
			}
		}
		clean := strings.ToLower(raw)
		if cfg, exists := collegeRegistry[clean]; exists {
			return cfg.BaseHost
		}
		if strings.Contains(clean, ".") {
			return clean
		}
		return clean + ".contineo.in"
	}

	// 2. USN prefix auto-detection (e.g. 1NT24CS006 -> 1NT -> nmitparents.contineo.in)
	cleanUSN := strings.ToUpper(strings.TrimSpace(usn))
	if len(cleanUSN) >= 3 {
		prefix := strings.ToLower(cleanUSN[:3])
		if cfg, exists := collegeRegistry[prefix]; exists {
			return cfg.BaseHost
		}
	}

	// 3. Environment variable or default host
	if host := os.Getenv("CONTINEO_BASE_HOST"); host != "" {
		return strings.TrimSpace(host)
	}

	return DefaultBaseHost
}

// RegisterCollege allows dynamic registration of college mappings at runtime.
func RegisterCollege(cfg CollegeConfig) {
	if cfg.Code != "" {
		collegeRegistry[strings.ToLower(cfg.Code)] = cfg
	}
	if cfg.USNPrefix != "" {
		collegeRegistry[strings.ToLower(cfg.USNPrefix)] = cfg
	}
}
