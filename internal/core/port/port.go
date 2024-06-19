package port

import (
	"context"
	"min/internal/core/domain"
)

// ShortenerRepository is an interface that defines the methods for the repository storing the shortened URLs.
type ShortenerRepository interface {
	// GetOriginal returns the original URL for the given short URL.
	GetOriginal(ctx context.Context, short string) (string, error)
	// Add stores the original URL and returns the shortened URL.
	Add(ctx context.Context, short, original string) error
	// Remove deletes the shortened URL.
	Remove(ctx context.Context, short string) error
}

// ShortenerCache is an interface that defines the methods for the cache storing the shortened URLs.
type ShortenerCache interface {
	// GetOriginal returns the original URL for the given short URL.
	GetOriginal(ctx context.Context, short string) (string, error)
	// Add stores the original URL for the given short URL.
	Add(ctx context.Context, short, original string) error
	// Remove deletes the shortened URL.
	Remove(ctx context.Context, short string) error
}

// ShortenerService is an interface that defines the methods for the shortener
// service. It is responsible for shortening and resolving URLs.
type ShortenerService interface {
	// Resolve returns the original URL for the given short URL.
	Resolve(ctx context.Context, short string) (string, error)
	// Shorten returns the shortened URL for the given original URL.
	Shorten(ctx context.Context, url string, author *domain.User) (string, error)
	// Remove deletes the shortened URL.
	Remove(ctx context.Context, short string) error
}

// UserRepository defines the interface for the user repository. It is used to store and retrieve user data.
type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	ChangeLinksRemaining(ctx context.Context, username string, linksRemaining int64) error
}

// AuthService defines the interface for the auth service.
type AuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
	Register(ctx context.Context, newUser *domain.User) error
	ValidateToken(ctx context.Context, tokenString string) (*domain.User, error)
	ChangeLinksRemaining(ctx context.Context, username string, linksRemaining int64) error
}

// AuthClient defines the interface for the auth client. It is used to communicate with the auth server.
type AuthClient interface {
	Register(ctx context.Context, username, password string, role domain.Role) error
	Login(ctx context.Context, username, password string) (string, error)
	ValidateToken(ctx context.Context, token string) (*domain.User, error)
	ChangeLinksRemaining(ctx context.Context, username string, linksRemaining int64) error
}

type StatisticsRepository interface {
	AddEvent(ctx context.Context, event domain.Event) error
}

type StatisticsService interface {
	AddEvent(ctx context.Context, event domain.Event) error
}

type EventProducer interface {
	Produce(event *domain.Event) error
}
