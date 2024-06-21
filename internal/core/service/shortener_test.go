package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"min/internal/core/domain"
	"min/internal/core/service"
	"min/internal/mocks"
)

func TestShortener_Resolve(t *testing.T) {
	repoMock := new(mocks.ShortenerRepository)
	cacheMock := new(mocks.ShortenerCache)
	authClientMock := new(mocks.AuthClient)
	shortener := service.NewShortener(repoMock, cacheMock, 8, authClientMock)

	t.Run("resolve from cache", func(t *testing.T) {
		cacheMock.On(
			"GetOriginal",
			mock.Anything,
			"shortUrl",
		).Return("http://original.url", nil).Once()

		original, err := shortener.Resolve(context.Background(), "shortUrl")
		require.NoError(t, err)
		assert.Equal(t, "http://original.url", original)
		cacheMock.AssertCalled(t, "GetOriginal", mock.Anything, "shortUrl")
	})

	t.Run("resolve from repository", func(t *testing.T) {
		cacheMock.On("GetOriginal", mock.Anything, "shortUrl").Return("", nil).Once()
		repoMock.On(
			"GetOriginal",
			mock.Anything,
			"shortUrl",
		).Return("http://original.url", nil).Once()

		original, err := shortener.Resolve(context.Background(), "shortUrl")
		require.NoError(t, err)
		assert.Equal(t, "http://original.url", original)
		cacheMock.AssertCalled(t, "GetOriginal", mock.Anything, "shortUrl")
		repoMock.AssertCalled(t, "GetOriginal", mock.Anything, "shortUrl")
	})

	t.Run("short URL not found", func(t *testing.T) {
		cacheMock.On("GetOriginal", mock.Anything, "shortUrl").Return("", nil).Once()
		repoMock.On("GetOriginal", mock.Anything, "shortUrl").Return("", nil).Once()

		original, err := shortener.Resolve(context.Background(), "shortUrl")
		require.Error(t, err)
		assert.Empty(t, original)
		assert.Equal(t, "short URL not found", err.Error())
	})
}

func TestShortener_Shorten(t *testing.T) {
	repoMock := new(mocks.ShortenerRepository)
	cacheMock := new(mocks.ShortenerCache)
	authClientMock := new(mocks.AuthClient)
	shortener := service.NewShortener(repoMock, cacheMock, 8, authClientMock)

	t.Run("successful shorten", func(t *testing.T) {
		user := &domain.User{Username: "user", LinksRemaining: 5}
		repoMock.On("Add", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		cacheMock.On(
			"Add",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil).Once()
		authClientMock.On(
			"ChangeLinksRemaining",
			mock.Anything,
			user.Username,
			int64(4),
		).Return(nil).Once()

		short, err := shortener.Shorten(context.Background(), "http://original.url", user)
		require.NoError(t, err)
		assert.NotEmpty(t, short)
		repoMock.AssertCalled(t, "Add", mock.Anything, mock.Anything, "http://original.url")
		cacheMock.AssertCalled(t, "Add", mock.Anything, mock.Anything, "http://original.url")
		authClientMock.AssertCalled(t, "ChangeLinksRemaining", mock.Anything, user.Username, int64(4))
	})

	t.Run("no links remaining", func(t *testing.T) {
		user := &domain.User{Username: "user", LinksRemaining: 0}

		short, err := shortener.Shorten(context.Background(), "http://original.url", user)
		require.Error(t, err)
		assert.Empty(t, short)
		assert.Equal(t, "no links remaining, please upgrade your account or remove some existing links", err.Error())
	})

	t.Run("failed to add to repository", func(t *testing.T) {
		user := &domain.User{Username: "user", LinksRemaining: 5}
		repoMock.On(
			"Add",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(errors.New("repo error"))

		short, err := shortener.Shorten(context.Background(), "http://original.url", user)
		require.Error(t, err)
		assert.Empty(t, short)
		assert.Contains(t, err.Error(), "failed to add short URL to repository")
	})
}

func TestShortener_Remove(t *testing.T) {
	repoMock := new(mocks.ShortenerRepository)
	cacheMock := new(mocks.ShortenerCache)
	authClientMock := new(mocks.AuthClient)
	shortener := service.NewShortener(repoMock, cacheMock, 8, authClientMock)

	t.Run("successful remove", func(t *testing.T) {
		repoMock.On("Remove", mock.Anything, "shortUrl").Return(nil).Once()
		cacheMock.On("Remove", mock.Anything, "shortUrl").Return(nil).Once()

		err := shortener.Remove(context.Background(), "shortUrl")
		require.NoError(t, err)
		repoMock.AssertCalled(t, "Remove", mock.Anything, "shortUrl")
		cacheMock.AssertCalled(t, "Remove", mock.Anything, "shortUrl")
	})

	t.Run("failed to remove from repository", func(t *testing.T) {
		repoMock.On("Remove", mock.Anything, "shortUrl").Return(errors.New("repo error")).Once()

		err := shortener.Remove(context.Background(), "shortUrl")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to remove short URL from repository")
	})

	t.Run("failed to remove from cache", func(t *testing.T) {
		repoMock.On("Remove", mock.Anything, "shortUrl").Return(nil).Once()
		cacheMock.On("Remove", mock.Anything, "shortUrl").Return(errors.New("cache error")).Once()

		err := shortener.Remove(context.Background(), "shortUrl")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to remove short URL from cache")
	})
}
