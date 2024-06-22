package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	authv1 "min/api/gen/go/auth"
	"min/api/gen/go/auth/mocks"
	"min/internal/adapter/client/auth"
	"min/internal/core/domain"
)

func TestClient_Login(t *testing.T) {
	mockAuthClient := new(mocks.AuthClient)
	client := &auth.Client{
		Conn:   &grpc.ClientConn{},
		Client: mockAuthClient,
	}

	t.Run("successful login", func(t *testing.T) {
		mockAuthClient.On("Login", mock.Anything, &authv1.LoginRequest{
			Username: "testuser",
			Password: "testpassword",
		}).Return(&authv1.LoginResponse{Token: "testtoken"}, nil).Once()

		token, err := client.Login(context.Background(), "testuser", "testpassword")
		require.NoError(t, err)
		assert.Equal(t, "testtoken", token)
	})

	t.Run("failed login", func(t *testing.T) {
		mockAuthClient.On("Login", mock.Anything, &authv1.LoginRequest{
			Username: "testuser",
			Password: "testpassword",
		}).Return(nil, errors.New("login error"))

		token, err := client.Login(context.Background(), "testuser", "testpassword")
		require.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "failed to login")
	})
}

func TestClient_Register(t *testing.T) {
	mockAuthClient := new(mocks.AuthClient)
	client := &auth.Client{
		Conn:   &grpc.ClientConn{},
		Client: mockAuthClient,
	}

	t.Run("successful registration", func(t *testing.T) {
		mockAuthClient.On(
			"Register",
			mock.Anything,
			mock.Anything,
		).Return(&authv1.RegisterResponse{}, nil).Once()

		err := client.Register(context.Background(), "newuser", "newpassword", domain.USER)
		require.NoError(t, err)
	})

	t.Run("failed registration", func(t *testing.T) {
		mockAuthClient.On(
			"Register",
			mock.Anything,
			mock.Anything,
		).Return(nil, errors.New("register error")).Once()

		err := client.Register(context.Background(), "newuser", "newpassword", domain.USER)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to register")
	})
}

func TestClient_ValidateToken(t *testing.T) {
	mockAuthClient := new(mocks.AuthClient)
	client := &auth.Client{
		Conn:   &grpc.ClientConn{},
		Client: mockAuthClient,
	}

	t.Run("successful token validation", func(t *testing.T) {
		mockAuthClient.On("ValidateToken", mock.Anything, &authv1.ValidateTokenRequest{
			Token: "validtoken",
		}).Return(&authv1.ValidateTokenResponse{
			Username:       "user",
			Password:       "password",
			Role:           "USER",
			Plan:           "FREE",
			LinksRemaining: 10,
		}, nil)

		user, err := client.ValidateToken(context.Background(), "validtoken")
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "user", user.Username)
	})

	t.Run("failed token validation", func(t *testing.T) {
		mockAuthClient.On("ValidateToken", mock.Anything, &authv1.ValidateTokenRequest{
			Token: "invalidtoken",
		}).Return(nil, errors.New("validation error"))

		user, err := client.ValidateToken(context.Background(), "invalidtoken")
		require.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to validate token")
	})
}

func TestClient_ChangeLinksRemaining(t *testing.T) {
	mockAuthClient := new(mocks.AuthClient)
	client := &auth.Client{
		Conn:   &grpc.ClientConn{},
		Client: mockAuthClient,
	}

	t.Run("successful change links remaining", func(t *testing.T) {
		mockAuthClient.On(
			"ChangeLinksRemaining",
			mock.Anything,
			mock.Anything,
		).Return(&authv1.ChangeLinksRemainingResponse{}, nil).Once()

		err := client.ChangeLinksRemaining(context.Background(), "user", 5)
		require.NoError(t, err)
	})

	t.Run("failed change links remaining", func(t *testing.T) {
		mockAuthClient.On("ChangeLinksRemaining", mock.Anything, &authv1.ChangeLinksRemainingRequest{
			Username:       "user",
			LinksRemaining: 5,
		}).Return(nil, errors.New("change error"))

		err := client.ChangeLinksRemaining(context.Background(), "user", 5)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to change links remaining")
	})
}
