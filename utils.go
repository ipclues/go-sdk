package ipclues

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

var Version = "0.1.0"

func UserAgent() string {
	return fmt.Sprintf(
		"IPCluesGo/%s (go/%s; os/%s; arch/%s)",
		Version,
		strings.TrimPrefix(runtime.Version(), "go"),
		runtime.GOOS,
		runtime.GOARCH,
	)
}

func headerParams(s string) map[string]string {
	out := map[string]string{}
	for _, part := range strings.Split(s, ";") {
		if i := strings.IndexByte(part, '='); i >= 0 {
			k := strings.Trim(strings.TrimSpace(part[:i]), `"`)
			v := strings.Trim(strings.TrimSpace(part[i+1:]), `"`)
			out[k] = v
		}
	}
	return out
}

func parseMetadata(status int, statusText string, h http.Header) ResponseMetadata {
	m := ResponseMetadata{Status: status, StatusText: statusText, RequestID: h.Get("I3-Request-Id"),
		CreditsCharged: parseInt64(h.Get("I3-Credits-Charged")), CreditsRemaining: parseInt64(h.Get("I3-Credits-Remaining")),
		Databases: splitCSV(h.Get("I3-Databases")), DatabaseVersions: parseKVCSV(h.Get("I3-Database-Versions"))}
	m.RateLimit = parseRateLimit(h.Get("Ratelimit"), h.Get("Ratelimit-Policy"), h.Get("Retry-After"))
	return m
}

func parseRateLimit(value, policy, retry string) *RateLimitMetadata {
	if value == "" && policy == "" && retry == "" {
		return nil
	}
	r := &RateLimitMetadata{Policy: policy}
	p := headerParams(value)
	pp := headerParams(policy)
	r.Remaining = parseInt(p["r"])
	r.Reset = parseInt(p["t"])
	r.Limit = parseInt(pp["q"])
	if r.Reset == nil {
		r.Reset = parseInt(pp["w"])
	}
	r.RetryAfter = parseInt(retry)
	return r
}

func parseKVCSV(s string) map[string]string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	out := map[string]string{}
	for item := range strings.SplitSeq(s, ",") {
		k, v, ok := strings.Cut(strings.TrimSpace(item), "=")
		if ok && k != "" && v != "" {
			out[strings.TrimSpace(k)] = strings.TrimSpace(v)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func parseInt(s string) *int {
	if s == "" {
		return nil
	}
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return nil
	}
	return &n
}

func parseInt64(s string) *int64 {
	if s == "" {
		return nil
	}
	n, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return nil
	}
	return &n
}

func splitCSV(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	var out []string
	for item := range strings.SplitSeq(s, ",") {
		if item = strings.TrimSpace(item); item != "" {
			out = append(out, item)
		}
	}
	return out
}
