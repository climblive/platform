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
	mockedContenderID := domain.ContenderID(rand.Int())
	mockedOwnership := domain.OwnershipData{
		OrganizerID: 1,
		ContenderID: &mockedContenderID,
	}

	mockedRepo := new(repositoryMock)

	mockedRepo.
		On("GetContenderByCode", mock.Anything, nil, "ABCD1234").
		Return(domain.Contender{
			ID: mockedContenderID,
		}, nil)

	mockedRepo.
		On("GetContenderByCode", mock.Anything, nil, mock.AnythingOfType("string")).
		Return(domain.Contender{}, domain.ErrNotFound)

	authorizer := authorizer.NewAuthorizer(mockedRepo)

	t.Run("MissingAuthorization", func(t *testing.T) {
		dummyHandler := func(w http.ResponseWriter, r *http.Request) {
			role, err := authorizer.HasOwnership(r.Context(), mockedOwnership)

			assert.Equal(t, domain.NilRole, role)
			assert.ErrorIs(t, err, domain.ErrNotAuthorized)
		}

		r := httptest.NewRequest("GET", "http://localhost", nil)
		w := httptest.NewRecorder()

		handler := authorizer.Middleware(http.HandlerFunc(dummyHandler))
		handler.ServeHTTP(w, r)
	})

	t.Run("BadAuthorization", func(t *testing.T) {
		dummyHandler := func(w http.ResponseWriter, r *http.Request) {
			role, err := authorizer.HasOwnership(r.Context(), mockedOwnership)

			assert.Equal(t, domain.NilRole, role)
			assert.ErrorIs(t, err, domain.ErrNotAuthorized)
		}

		r := httptest.NewRequest("GET", "http://localhost", nil)
		w := httptest.NewRecorder()

		r.Header.Set("Authorization", "Regcode DEADBEEF")

		handler := authorizer.Middleware(http.HandlerFunc(dummyHandler))
		handler.ServeHTTP(w, r)
	})

	t.Run("BadSyntax", func(t *testing.T) {
		dummyHandler := func(w http.ResponseWriter, r *http.Request) {
			role, err := authorizer.HasOwnership(r.Context(), mockedOwnership)

			assert.Equal(t, domain.NilRole, role)
			assert.ErrorIs(t, err, domain.ErrNotAuthorized)
		}

		r := httptest.NewRequest("GET", "http://localhost", nil)
		w := httptest.NewRecorder()

		r.Header.Set("Authorization", "totally_wrong")

		handler := authorizer.Middleware(http.HandlerFunc(dummyHandler))
		handler.ServeHTTP(w, r)
	})

	t.Run("AuthorizedWithOwnership", func(t *testing.T) {
		dummyHandler := func(w http.ResponseWriter, r *http.Request) {
			role, err := authorizer.HasOwnership(r.Context(), mockedOwnership)

			assert.Equal(t, domain.ContenderRole, role)
			assert.NoError(t, err)
		}

		r := httptest.NewRequest("GET", "http://localhost", nil)
		w := httptest.NewRecorder()

		r.Header.Set("Authorization", "Regcode ABCD1234")

		handler := authorizer.Middleware(http.HandlerFunc(dummyHandler))
		handler.ServeHTTP(w, r)
	})

	t.Run("AuthorizedWithoutOwnership", func(t *testing.T) {
		otherContenderID := mockedContenderID + 1

		mockedOtherOwnership := domain.OwnershipData{
			OrganizerID: 1,
			ContenderID: &otherContenderID,
		}

		dummyHandler := func(w http.ResponseWriter, r *http.Request) {
			role, err := authorizer.HasOwnership(r.Context(), mockedOtherOwnership)

			assert.Equal(t, domain.NilRole, role)
			assert.ErrorIs(t, err, domain.ErrNoOwnership)
		}

		r := httptest.NewRequest("GET", "http://localhost", nil)
		w := httptest.NewRecorder()

		r.Header.Set("Authorization", "Regcode ABCD1234")

		handler := authorizer.Middleware(http.HandlerFunc(dummyHandler))
		handler.ServeHTTP(w, r)
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
