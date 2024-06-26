package service

import (
	"context"
	"errors"
	"fmt"
	"min/internal/core/domain"
	"min/internal/core/port"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService is the service for authentication.
type AuthService struct {
	authRep      port.UserRepository
	tokenMaxTime time.Duration
}

// NewAuthService creates a new AuthService.
func NewAuthService(authRep port.UserRepository, tokenMaxTime time.Duration) *AuthService {
	return &AuthService{
		authRep:      authRep,
		tokenMaxTime: tokenMaxTime,
	}
}

// Login logs in the user and returns a token.
func (a *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := a.authRep.GetByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return "", errors.New("user not found")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(a.tokenMaxTime).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret")) // For demonstration purposes only. Should be hidden.
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// Register registers a new user.
func (a *AuthService) Register(ctx context.Context, newUser *domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	newUser.Password = string(hashedPassword)

	err = a.authRep.Save(ctx, newUser)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

// ValidateToken validates the token and returns the user associated with it.
func (a *AuthService) ValidateToken(ctx context.Context, tokenString string) (*domain.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil // For demonstration purposes only. Should be hidden.
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, errors.New("invalid username in claims")
	}

	user, err := a.authRep.GetByUsername(ctx, username)

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// ChangeLinksRemaining changes the remaining links for the specified user.
func (a *AuthService) ChangeLinksRemaining(ctx context.Context, username string, linksRemaining int64) error {
	err := a.authRep.ChangeLinksRemaining(ctx, username, linksRemaining)
	if err != nil {
		return fmt.Errorf("failed to change links remaining: %w", err)
	}

	return nil
}
