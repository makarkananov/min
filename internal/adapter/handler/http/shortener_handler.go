package http

import (
	log "github.com/sirupsen/logrus"
	"min/internal/core/domain"
	"min/internal/core/port"
	"net/http"
	"net/url"
)

// ShortenerHandler provides methods for handling redirect requests and shorten requests.
type ShortenerHandler struct {
	shortenerService port.ShortenerService
	eventProducer    port.EventProducer
}

// NewShortenerHandler creates a new instance of ShortenerHandler.
func NewShortenerHandler(
	shortenerService port.ShortenerService,
	eventProducer port.EventProducer,
) *ShortenerHandler {
	return &ShortenerHandler{
		shortenerService: shortenerService,
		eventProducer:    eventProducer,
	}
}

// Redirect handles redirect requests by trying to resolve the short URL and redirecting to the original URL.
func (sh *ShortenerHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	log.Infof("Got request to redirect: %s", r.URL.Path)
	short := r.URL.Path[1:]
	if short == "" {
		log.Errorf("Short URL is required")
		http.Error(w, "Short URL is required", http.StatusBadRequest)
		return
	}

	original, err := sh.shortenerService.Resolve(r.Context(), short)
	if err != nil {
		log.Errorf("Failed to resolve URL: %v", err)
		http.Error(w, "Failed to resolve URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = sh.eventProducer.Produce(domain.NewEvent(short, original, r.UserAgent(), r.RemoteAddr))
	if err != nil {
		log.Errorf("Failed to produce event: %v", err)
		http.Error(w, "Failed to produce event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("Successfully redirecting to: %s", original)
	http.Redirect(w, r, original, http.StatusPermanentRedirect)
}

// Shorten handles shorten requests by shortening the original URL and returning the shortened URL.
func (sh *ShortenerHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	original := r.URL.Query().Get("url")
	log.Infof("Got request to shorten: %s", original)
	if original == "" {
		log.Errorf("Original URL is required")
		http.Error(w, "Original URL is required", http.StatusBadRequest)
		return
	}

	log.Infof("Request to shorten made by user: %v", r.Context().Value(currentUserKey))
	userData := r.Context().Value(currentUserKey)
	user, _ := userData.(*domain.User)
	if user == nil {
		log.Errorf("User is required to perform this action")
		http.Error(w, "User is required to perform this action", http.StatusBadRequest)
		return
	}

	short, err := sh.shortenerService.Shorten(r.Context(), original, user)
	if err != nil {
		log.Errorf("Failed to shorten URL: %v", err)
		http.Error(w, "Failed to shorten URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	u := &url.URL{
		Scheme: "http",
		Host:   r.Host,
		Path:   "/" + short,
	}
	fullURL := u.String()

	_, err = w.Write([]byte(fullURL))
	if err != nil {
		log.Errorf("Failed to write response: %v", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

// Remove handles remove requests by removing the short URL from the cache.
func (sh *ShortenerHandler) Remove(w http.ResponseWriter, r *http.Request) {
	short := r.URL.Query().Get("url")
	log.Infof("Got request to remove: %s", short)
	if short == "" {
		log.Errorf("Short URL is required")
		http.Error(w, "Short URL is required", http.StatusBadRequest)
		return
	}

	err := sh.shortenerService.Remove(r.Context(), short)
	if err != nil {
		log.Errorf("Failed to remove URL: %v", err)
		http.Error(w, "Failed to remove URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("Successfully removed: %s", short)
	w.WriteHeader(http.StatusOK)
}
