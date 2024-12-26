package authorizer_test

import (
	"context"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/climblive/platform/backend/internal/authorizer"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthorizer(t *testing.T) {
	fakedContenderID := domain.ContenderID(rand.Int())
	fakedOrganizerID := domain.OrganizerID(rand.Int())

	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
		ContenderID: &fakedContenderID,
	}

	t.Run("MissingAuthorization", func(t *testing.T) {
		mockedRepo := new(repositoryMock)

		authorizer := authorizer.NewAuthorizer(mockedRepo)

		dummyHandler := func(w http.ResponseWriter, r *http.Request) {
			role, err := authorizer.HasOwnership(r.Context(), fakedOwnership)

			assert.Equal(t, domain.NilRole, role)
			assert.ErrorIs(t, err, domain.ErrNotAuthorized)
		}

		r := httptest.NewRequest("GET", "http://localhost", nil)
		w := httptest.NewRecorder()

		handler := authorizer.Middleware(http.HandlerFunc(dummyHandler))
		handler.ServeHTTP(w, r)

		mockedRepo.AssertExpectations(t)
	})

	t.Run("BadAuthorization", func(t *testing.T) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetContenderByCode", mock.Anything, nil, mock.AnythingOfType("string")).
			Return(domain.Contender{}, domain.ErrNotFound)

		authorizer := authorizer.NewAuthorizer(mockedRepo)

		dummyHandler := func(w http.ResponseWriter, r *http.Request) {
			role, err := authorizer.HasOwnership(r.Context(), fakedOwnership)

			assert.Equal(t, domain.NilRole, role)
			assert.ErrorIs(t, err, domain.ErrNotAuthorized)
		}

		r := httptest.NewRequest("GET", "http://localhost", nil)
		w := httptest.NewRecorder()

		r.Header.Set("Authorization", "Regcode DEADBEEF")

		handler := authorizer.Middleware(http.HandlerFunc(dummyHandler))
		handler.ServeHTTP(w, r)

		mockedRepo.AssertExpectations(t)
	})

	t.Run("BadSyntax", func(t *testing.T) {
		mockedRepo := new(repositoryMock)

		authorizer := authorizer.NewAuthorizer(mockedRepo)

		dummyHandler := func(w http.ResponseWriter, r *http.Request) {
			role, err := authorizer.HasOwnership(r.Context(), fakedOwnership)

			assert.Equal(t, domain.NilRole, role)
			assert.ErrorIs(t, err, domain.ErrNotAuthorized)
		}

		r := httptest.NewRequest("GET", "http://localhost", nil)
		w := httptest.NewRecorder()

		r.Header.Set("Authorization", "totally_wrong")

		handler := authorizer.Middleware(http.HandlerFunc(dummyHandler))
		handler.ServeHTTP(w, r)

		mockedRepo.AssertExpectations(t)
	})

	t.Run("AuthorizedWithOwnership", func(t *testing.T) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetContenderByCode", mock.Anything, nil, "ABCD1234").
			Return(domain.Contender{
				ID: fakedContenderID,
			}, nil)

		authorizer := authorizer.NewAuthorizer(mockedRepo)

		dummyHandler := func(w http.ResponseWriter, r *http.Request) {
			role, err := authorizer.HasOwnership(r.Context(), fakedOwnership)

			assert.Equal(t, domain.ContenderRole, role)
			assert.NoError(t, err)
		}

		r := httptest.NewRequest("GET", "http://localhost", nil)
		w := httptest.NewRecorder()

		r.Header.Set("Authorization", "Regcode ABCD1234")

		handler := authorizer.Middleware(http.HandlerFunc(dummyHandler))
		handler.ServeHTTP(w, r)

		mockedRepo.AssertExpectations(t)
	})

	t.Run("AuthorizedWithoutOwnership", func(t *testing.T) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetContenderByCode", mock.Anything, nil, "ABCD1234").
			Return(domain.Contender{
				ID: fakedContenderID,
			}, nil)

		otherContenderID := fakedContenderID + 1

		fakedOtherOwnership := domain.OwnershipData{
			OrganizerID: fakedOrganizerID,
			ContenderID: &otherContenderID,
		}

		authorizer := authorizer.NewAuthorizer(mockedRepo)

		dummyHandler := func(w http.ResponseWriter, r *http.Request) {
			role, err := authorizer.HasOwnership(r.Context(), fakedOtherOwnership)

			assert.Equal(t, domain.NilRole, role)
			assert.ErrorIs(t, err, domain.ErrNoOwnership)
		}

		r := httptest.NewRequest("GET", "http://localhost", nil)
		w := httptest.NewRecorder()

		r.Header.Set("Authorization", "Regcode ABCD1234")

		handler := authorizer.Middleware(http.HandlerFunc(dummyHandler))
		handler.ServeHTTP(w, r)

		mockedRepo.AssertExpectations(t)
	})

	t.Run("CodesConvertedToUpperCase", func(t *testing.T) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetContenderByCode", mock.Anything, nil, "WXYZ1234").
			Return(domain.Contender{}, nil)

		authorizer := authorizer.NewAuthorizer(mockedRepo)

		dummyHandler := func(w http.ResponseWriter, r *http.Request) {
			_, _ = authorizer.HasOwnership(r.Context(), fakedOwnership)
		}

		r := httptest.NewRequest("GET", "http://localhost", nil)
		w := httptest.NewRecorder()

		r.Header.Set("Authorization", "Regcode wxyz1234")

		handler := authorizer.Middleware(http.HandlerFunc(dummyHandler))
		handler.ServeHTTP(w, r)

		mockedRepo.AssertExpectations(t)
	})
}

type repositoryMock struct {
	mock.Mock
}

func (m *repositoryMock) Begin() domain.Transaction {
	args := m.Called()
	return args.Get(0).(domain.Transaction)
}

func (m *repositoryMock) GetContenderByCode(ctx context.Context, tx domain.Transaction, registrationCode string) (domain.Contender, error) {
	args := m.Called(ctx, tx, registrationCode)
	return args.Get(0).(domain.Contender), args.Error(1)
}
