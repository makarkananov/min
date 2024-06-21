package http

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"min/internal/core/domain"
	"min/internal/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShortenerHandler_Redirect(t *testing.T) {
	shortenerServiceMock := new(mocks.ShortenerService)
	eventProducerMock := new(mocks.EventProducer)
	handler := NewShortenerHandler(shortenerServiceMock, eventProducerMock)

	t.Run("successful redirect", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/shortUrl", nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()

		shortenerServiceMock.On(
			"Resolve",
			mock.Anything,
			"shortUrl",
		).Return("http://original.url", nil).Once()
		eventProducerMock.On("Produce", mock.Anything).Return(nil)

		handler.Redirect(rr, req)

		assert.Equal(t, http.StatusPermanentRedirect, rr.Code)
		assert.Equal(t, "http://original.url", rr.Header().Get("Location"))
		shortenerServiceMock.AssertCalled(t, "Resolve", mock.Anything, "shortUrl")
		eventProducerMock.AssertCalled(t, "Produce", mock.Anything)
	})

	t.Run("missing short URL", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()

		handler.Redirect(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("resolve error", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/shortUrl", nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()

		shortenerServiceMock.On(
			"Resolve",
			mock.Anything,
			"shortUrl",
		).Return("", errors.New("resolve error"))

		handler.Redirect(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		shortenerServiceMock.AssertCalled(t, "Resolve", mock.Anything, "shortUrl")
	})

	t.Run("produce event error", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/shortUrl", nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()

		shortenerServiceMock.On("Resolve", mock.Anything, "shortUrl").Return("http://original.url", nil)
		eventProducerMock.On("Produce", mock.Anything).Return(errors.New("produce error"))

		handler.Redirect(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		shortenerServiceMock.AssertCalled(t, "Resolve", mock.Anything, "shortUrl")
		eventProducerMock.AssertCalled(t, "Produce", mock.Anything)
	})
}

func TestShortenerHandler_Shorten(t *testing.T) {
	shortenerServiceMock := new(mocks.ShortenerService)
	eventProducerMock := new(mocks.EventProducer)
	handler := NewShortenerHandler(shortenerServiceMock, eventProducerMock)

	t.Run("successful shorten", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/shorten?url=http://original.url", nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()

		ctx := context.WithValue(req.Context(), currentUserKey, &domain.User{Username: "user1"})
		req = req.WithContext(ctx)

		shortenerServiceMock.On(
			"Shorten",
			mock.Anything,
			"http://original.url",
			&domain.User{Username: "user1"},
		).Return("shortUrl", nil).Once()

		handler.Shorten(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "http://"+req.Host+"/shortUrl", rr.Body.String())
		shortenerServiceMock.AssertCalled(
			t,
			"Shorten",
			mock.Anything,
			"http://original.url",
			&domain.User{Username: "user1"},
		)
	})

	t.Run("missing original URL", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/shorten", nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()

		handler.Shorten(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("missing user", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/shorten?url=http://original.url", nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()

		handler.Shorten(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("shorten error", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/shorten?url=http://original.url", nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()

		ctx := context.WithValue(req.Context(), currentUserKey, &domain.User{Username: "user1"})
		req = req.WithContext(ctx)

		shortenerServiceMock.On(
			"Shorten",
			mock.Anything,
			"http://original.url",
			&domain.User{Username: "user1"},
		).Return("", errors.New("shorten error"))

		handler.Shorten(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		shortenerServiceMock.AssertCalled(
			t,
			"Shorten",
			mock.Anything,
			"http://original.url",
			&domain.User{Username: "user1"},
		)
	})
}

func TestShortenerHandler_Remove(t *testing.T) {
	shortenerServiceMock := new(mocks.ShortenerService)
	eventProducerMock := new(mocks.EventProducer)
	handler := NewShortenerHandler(shortenerServiceMock, eventProducerMock)

	t.Run("successful remove", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/remove?url=shortUrl", nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()

		shortenerServiceMock.On("Remove", mock.Anything, "shortUrl").Return(nil).Once()

		handler.Remove(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		shortenerServiceMock.AssertCalled(t, "Remove", mock.Anything, "shortUrl")
	})

	t.Run("missing short URL", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/remove", nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()

		handler.Remove(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("remove error", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/remove?url=shortUrl", nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()

		shortenerServiceMock.On("Remove", mock.Anything, "shortUrl").Return(errors.New("remove error"))

		handler.Remove(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		shortenerServiceMock.AssertCalled(t, "Remove", mock.Anything, "shortUrl")
	})
}
