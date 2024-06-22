package domain

import "time"

// Event represents an income request that occurred in the redirecting system.
type Event struct {
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	Timestamp   time.Time `json:"timestamp"`
	UserAgent   string    `json:"user_agent"`
	IP          string    `json:"ip"`
}

// NewEvent creates a new event with the given short URL, original URL, user agent, and IP.
func NewEvent(shortURL, originalURL, userAgent, ip string) *Event {
	return &Event{
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		Timestamp:   time.Now(),
		UserAgent:   userAgent,
		IP:          ip,
	}
}
