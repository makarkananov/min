package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"min/internal/core/domain"
	"min/internal/core/port"
)

type Shortener struct {
	repository    port.ShortenerRepository
	cache         port.ShortenerCache
	authClient    port.AuthClient
	shortenLength int
}

func NewShortener(
	repository port.ShortenerRepository,
	cache port.ShortenerCache,
	shortenLength int,
	authClient port.AuthClient,
) *Shortener {
	return &Shortener{
		repository:    repository,
		cache:         cache,
		shortenLength: shortenLength,
		authClient:    authClient,
	}
}

func (s *Shortener) Resolve(ctx context.Context, short string) (string, error) {
	original, err := s.cache.GetOriginal(ctx, short)
	if err != nil {
		return "", fmt.Errorf("failed to get short url from cache: %w", err)
	}

	if original == "" {
		original, err = s.repository.GetOriginal(ctx, short)
		if err != nil {
			return "", fmt.Errorf("failed to get original url from repository: %w", err)
		}
	}

	if original == "" {
		return "", errors.New("short URL not found")
	}

	return original, nil
}

func (s *Shortener) Shorten(ctx context.Context, url string, author *domain.User) (string, error) {
	if author.LinksRemaining <= 0 {
		return "", errors.New("no links remaining, please upgrade your account or remove some existing links")
	}

	shorten, err := s.generateShortURL()
	if err != nil {
		return "", fmt.Errorf("failed to generate short URL: %w", err)
	}

	if err := s.repository.Add(ctx, shorten, url); err != nil {
		return "", fmt.Errorf("failed to add short URL to repository: %w", err)
	}

	if err := s.cache.Add(ctx, shorten, url); err != nil {
		return "", fmt.Errorf("failed to add short URL to cache: %w", err)
	}

	err = s.authClient.ChangeLinksRemaining(ctx, author.Username, author.LinksRemaining-1)
	if err != nil {
		log.Infof("Failed to change links remaining: %v", err)
		return "", fmt.Errorf("failed to change links remaining: %w", err)
	}

	return shorten, nil
}

func (s *Shortener) Remove(ctx context.Context, short string) error {
	if err := s.repository.Remove(ctx, short); err != nil {
		return fmt.Errorf("failed to remove short URL from repository: %w", err)
	}

	if err := s.cache.Remove(ctx, short); err != nil {
		return fmt.Errorf("failed to remove short URL from cache: %w", err)
	}

	return nil
}

// generateShortURL generates a random short URL.
func (s *Shortener) generateShortURL() (string, error) {
	// Generate shortenLength random bytes
	bytes := make([]byte, s.shortenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	// Encode the bytes to a base64 URL safe string
	return base64.URLEncoding.EncodeToString(bytes), nil
}
