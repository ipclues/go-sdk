package ipclues

// RateLimitMetadata represents rate-limit headers.
type RateLimitMetadata struct {
	Policy     string `json:"policy,omitempty"`
	Limit      *int   `json:"limit,omitempty"`
	Remaining  *int   `json:"remaining,omitempty"`
	Reset      *int   `json:"reset,omitempty"`
	RetryAfter *int   `json:"retry_after,omitempty"`
}

// ResponseMetadata contains status and IPClues response-header metadata.
type ResponseMetadata struct {
	Status           int                `json:"status"`
	StatusText       string             `json:"status_text"`
	RequestID        string             `json:"request_id,omitempty"`
	CreditsCharged   *int64             `json:"credits_charged,omitempty"`
	CreditsRemaining *int64             `json:"credits_remaining,omitempty"`
	Databases        []string           `json:"databases,omitempty"`
	DatabaseVersions map[string]string  `json:"database_versions,omitempty"`
	RateLimit        *RateLimitMetadata `json:"rate_limit,omitempty"`
}

// Result wraps the API response body with HTTP headers.
type Result struct {
	Data     any              `json:"data"`
	Metadata ResponseMetadata `json:"metadata"`
}
