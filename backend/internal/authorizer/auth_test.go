package authorizer_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/climblive/platform/backend/internal/authorizer"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthorizer(t *testing.T) {
	fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedUserID := testutils.RandomResourceID[domain.UserID]()

	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
		ContenderID: &fakedContenderID,
	}

	t.Run("MissingAuthorization", func(t *testing.T) {
		mockedRepo := new(repositoryMock)
		mockedJWTDecoder := new(jwtDecoderMock)

		authorizer := authorizer.NewAuthorizer(mockedRepo, mockedJWTDecoder)

		dummyHandler := func(w http.ResponseWriter, r *http.Request) {
			role, err := authorizer.HasOwnership(r.Context(), fakedOwnership)

			assert.Equal(t, domain.NilRole, role)
			assert.ErrorIs(t, err, domain.ErrNotAuthenticated)

			_, err = authorizer.GetAuthentication(r.Context())
			assert.ErrorIs(t, err, domain.ErrNotAuthenticated)
		}

		r := httptest.NewRequest("GET", "http://localhost", nil)
		w := httptest.NewRecorder()

		handler := authorizer.Middleware(http.HandlerFunc(dummyHandler))
		handler.ServeHTTP(w, r)

		mockedRepo.AssertExpectations(t)
		mockedJWTDecoder.AssertExpectations(t)
	})

	t.Run("Contender", func(t *testing.T) {
		t.Run("BadAuthorization", func(t *testing.T) {
			mockedRepo := new(repositoryMock)
			mockedJWTDecoder := new(jwtDecoderMock)

			mockedRepo.
				On("GetContenderByCode", mock.Anything, nil, mock.AnythingOfType("string")).
				Return(domain.Contender{}, domain.ErrNotFound)

			authorizer := authorizer.NewAuthorizer(mockedRepo, mockedJWTDecoder)

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
			mockedJWTDecoder.AssertExpectations(t)
		})

		t.Run("BadSyntax", func(t *testing.T) {
			mockedRepo := new(repositoryMock)
			mockedJWTDecoder := new(jwtDecoderMock)

			authorizer := authorizer.NewAuthorizer(mockedRepo, mockedJWTDecoder)

			dummyHandler := func(w http.ResponseWriter, r *http.Request) {
				role, err := authorizer.HasOwnership(r.Context(), fakedOwnership)

				assert.Equal(t, domain.NilRole, role)
				assert.ErrorIs(t, err, domain.ErrNotAuthenticated)
			}

			r := httptest.NewRequest("GET", "http://localhost", nil)
			w := httptest.NewRecorder()

			r.Header.Set("Authorization", "totally_wrong")

			handler := authorizer.Middleware(http.HandlerFunc(dummyHandler))
			handler.ServeHTTP(w, r)

			mockedRepo.AssertExpectations(t)
			mockedJWTDecoder.AssertExpectations(t)
		})

		t.Run("AuthorizedWithOwnership", func(t *testing.T) {
			mockedRepo := new(repositoryMock)
			mockedJWTDecoder := new(jwtDecoderMock)

			mockedRepo.
				On("GetContenderByCode", mock.Anything, nil, "ABCD1234").
				Return(domain.Contender{
					ID: fakedContenderID,
				}, nil)

			authorizer := authorizer.NewAuthorizer(mockedRepo, mockedJWTDecoder)

			dummyHandler := func(w http.ResponseWriter, r *http.Request) {
				role, err := authorizer.HasOwnership(r.Context(), fakedOwnership)

				assert.Equal(t, domain.ContenderRole, role)
				assert.NoError(t, err)

				authentication, err := authorizer.GetAuthentication(r.Context())
				assert.NoError(t, err)
				assert.Equal(t, "ABCD1234", authentication.Regcode)
				assert.Empty(t, authentication.Username)
			}

			r := httptest.NewRequest("GET", "http://localhost", nil)
			w := httptest.NewRecorder()

			r.Header.Set("Authorization", "Regcode ABCD1234")

			handler := authorizer.Middleware(http.HandlerFunc(dummyHandler))
			handler.ServeHTTP(w, r)

			mockedRepo.AssertExpectations(t)
			mockedJWTDecoder.AssertExpectations(t)
		})

		t.Run("AuthorizedWithoutOwnership", func(t *testing.T) {
			mockedRepo := new(repositoryMock)
			mockedJWTDecoder := new(jwtDecoderMock)

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

			authorizer := authorizer.NewAuthorizer(mockedRepo, mockedJWTDecoder)

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
			mockedJWTDecoder.AssertExpectations(t)
		})

		t.Run("CodesConvertedToUpperCase", func(t *testing.T) {
			mockedRepo := new(repositoryMock)
			mockedJWTDecoder := new(jwtDecoderMock)

			mockedRepo.
				On("GetContenderByCode", mock.Anything, nil, "WXYZ1234").
				Return(domain.Contender{}, nil)

			authorizer := authorizer.NewAuthorizer(mockedRepo, mockedJWTDecoder)

			dummyHandler := func(w http.ResponseWriter, r *http.Request) {
				_, _ = authorizer.HasOwnership(r.Context(), fakedOwnership)
			}

			r := httptest.NewRequest("GET", "http://localhost", nil)
			w := httptest.NewRecorder()

			r.Header.Set("Authorization", "Regcode wxyz1234")

			handler := authorizer.Middleware(http.HandlerFunc(dummyHandler))
			handler.ServeHTTP(w, r)

			mockedRepo.AssertExpectations(t)
			mockedJWTDecoder.AssertExpectations(t)
		})
	})

	t.Run("Organizer", func(t *testing.T) {
		t.Run("BadToken", func(t *testing.T) {
			mockedRepo := new(repositoryMock)
			mockedJWTDecoder := new(jwtDecoderMock)

			mockedJWTDecoder.
				On("Decode", "some_jwt").
				Return(authorizer.Claims{}, authorizer.ErrExpiredCredentials)

			auth := authorizer.NewAuthorizer(mockedRepo, mockedJWTDecoder)

			dummyHandler := func(w http.ResponseWriter, r *http.Request) {
				role, err := auth.HasOwnership(r.Context(), fakedOwnership)

				assert.Equal(t, domain.NilRole, role)
				assert.ErrorIs(t, err, domain.ErrNotAuthenticated)
				assert.ErrorIs(t, err, authorizer.ErrExpiredCredentials)

				_, err = auth.GetAuthentication(r.Context())
				assert.ErrorIs(t, err, domain.ErrNotAuthenticated)
				assert.ErrorIs(t, err, authorizer.ErrExpiredCredentials)
			}

			r := httptest.NewRequest("GET", "http://localhost", nil)
			w := httptest.NewRecorder()

			r.Header.Set("Authorization", "Bearer some_jwt")

			handler := auth.Middleware(http.HandlerFunc(dummyHandler))
			handler.ServeHTTP(w, r)

			mockedRepo.AssertExpectations(t)
			mockedJWTDecoder.AssertExpectations(t)
		})

		t.Run("AuthorizedWithOwnership", func(t *testing.T) {
			mockedRepo := new(repositoryMock)
			mockedJWTDecoder := new(jwtDecoderMock)

			mockedRepo.
				On("GetUserByUsername", mock.Anything, nil, "john").
				Return(domain.User{
					ID:         fakedUserID,
					Username:   "john",
					Admin:      false,
					Organizers: []domain.Organizer{{ID: fakedOrganizerID}},
				}, nil)

			mockedJWTDecoder.
				On("Decode", "some_jwt").
				Return(authorizer.Claims{
					Username: "john",
				}, nil)

			authorizer := authorizer.NewAuthorizer(mockedRepo, mockedJWTDecoder)

			dummyHandler := func(w http.ResponseWriter, r *http.Request) {
				role, err := authorizer.HasOwnership(r.Context(), fakedOwnership)

				assert.Equal(t, domain.OrganizerRole, role)
				assert.NoError(t, err)

				authentication, err := authorizer.GetAuthentication(r.Context())

				assert.NoError(t, err)
				assert.Equal(t, "john", authentication.Username)
				assert.Empty(t, authentication.Regcode)
			}

			r := httptest.NewRequest("GET", "http://localhost", nil)
			w := httptest.NewRecorder()

			r.Header.Set("Authorization", "Bearer some_jwt")

			handler := authorizer.Middleware(http.HandlerFunc(dummyHandler))
			handler.ServeHTTP(w, r)

			mockedRepo.AssertExpectations(t)
			mockedJWTDecoder.AssertExpectations(t)
		})

		t.Run("AuthorizedAsAdmin", func(t *testing.T) {
			mockedRepo := new(repositoryMock)
			mockedJWTDecoder := new(jwtDecoderMock)

			mockedRepo.
				On("GetUserByUsername", mock.Anything, nil, "john").
				Return(domain.User{
					ID:         fakedUserID,
					Username:   "john",
					Admin:      true,
					Organizers: []domain.Organizer{},
				}, nil)

			mockedJWTDecoder.
				On("Decode", "some_jwt").
				Return(authorizer.Claims{
					Username: "john",
				}, nil)

			authorizer := authorizer.NewAuthorizer(mockedRepo, mockedJWTDecoder)

			dummyHandler := func(w http.ResponseWriter, r *http.Request) {
				role, err := authorizer.HasOwnership(r.Context(), fakedOwnership)

				assert.Equal(t, domain.AdminRole, role)
				assert.NoError(t, err)
			}

			r := httptest.NewRequest("GET", "http://localhost", nil)
			w := httptest.NewRecorder()

			r.Header.Set("Authorization", "Bearer some_jwt")

			handler := authorizer.Middleware(http.HandlerFunc(dummyHandler))
			handler.ServeHTTP(w, r)

			mockedRepo.AssertExpectations(t)
			mockedJWTDecoder.AssertExpectations(t)
		})

		t.Run("AuthorizedWithoutOwnership", func(t *testing.T) {
			mockedRepo := new(repositoryMock)
			mockedJWTDecoder := new(jwtDecoderMock)

			mockedRepo.
				On("GetUserByUsername", mock.Anything, nil, "john").
				Return(domain.User{
					ID:         fakedUserID,
					Username:   "john",
					Admin:      false,
					Organizers: []domain.Organizer{{ID: fakedOrganizerID + 1}},
				}, nil)

			mockedJWTDecoder.
				On("Decode", "some_jwt").
				Return(authorizer.Claims{
					Username: "john",
				}, nil)

			authorizer := authorizer.NewAuthorizer(mockedRepo, mockedJWTDecoder)

			dummyHandler := func(w http.ResponseWriter, r *http.Request) {
				role, err := authorizer.HasOwnership(r.Context(), fakedOwnership)

				assert.Equal(t, domain.NilRole, role)
				assert.ErrorIs(t, err, domain.ErrNoOwnership)
			}

			r := httptest.NewRequest("GET", "http://localhost", nil)
			w := httptest.NewRecorder()

			r.Header.Set("Authorization", "Bearer some_jwt")

			handler := authorizer.Middleware(http.HandlerFunc(dummyHandler))
			handler.ServeHTTP(w, r)

			mockedRepo.AssertExpectations(t)
			mockedJWTDecoder.AssertExpectations(t)
		})

		t.Run("CreateUser", func(t *testing.T) {
			mockedRepo := new(repositoryMock)
			mockedJWTDecoder := new(jwtDecoderMock)
			mockedTx := new(transactionMock)

			fakedNewOrganizerID := fakedOrganizerID + 1
			fakedNewUserID := fakedUserID + 1

			mockedRepo.
				On("GetUserByUsername", mock.Anything, nil, "john").
				Return(domain.User{}, domain.ErrNotFound)

			mockedRepo.
				On("Begin").
				Return(mockedTx, nil)

			mockedRepo.
				On("StoreOrganizer", mock.Anything, mockedTx, domain.Organizer{
					Name: "john's organizer",
				}).
				Return(domain.Organizer{
					ID:   fakedNewOrganizerID,
					Name: "john's organizer",
				}, nil)

			mockedRepo.
				On("StoreUser", mock.Anything, mockedTx, domain.User{
					Username: "john",
				}).
				Return(domain.User{
					ID:       fakedNewUserID,
					Username: "john",
				}, nil)

			mockedRepo.
				On("AddUserToOrganizer", mock.Anything, mockedTx, fakedNewUserID, fakedNewOrganizerID).
				Return(nil)

			mockedTx.On("Commit").Return(nil)

			mockedJWTDecoder.
				On("Decode", "some_jwt").
				Return(authorizer.Claims{
					Username: "john",
				}, nil)

			authorizer := authorizer.NewAuthorizer(mockedRepo, mockedJWTDecoder)

			dummyHandler := func(w http.ResponseWriter, r *http.Request) {
				role, err := authorizer.HasOwnership(r.Context(), fakedOwnership)

				assert.Equal(t, domain.NilRole, role)
				assert.ErrorIs(t, err, domain.ErrNoOwnership)
			}

			r := httptest.NewRequest("GET", "http://localhost", nil)
			w := httptest.NewRecorder()

			r.Header.Set("Authorization", "Bearer some_jwt")

			handler := authorizer.Middleware(http.HandlerFunc(dummyHandler))
			handler.ServeHTTP(w, r)

			mockedRepo.AssertExpectations(t)
			mockedJWTDecoder.AssertExpectations(t)
			mockedTx.AssertExpectations(t)
		})
	})
}

type repositoryMock struct {
	mock.Mock
}

func (m *repositoryMock) Begin() (domain.Transaction, error) {
	args := m.Called()
	return args.Get(0).(domain.Transaction), args.Error(1)
}

func (m *repositoryMock) GetContenderByCode(ctx context.Context, tx domain.Transaction, registrationCode string) (domain.Contender, error) {
	args := m.Called(ctx, tx, registrationCode)
	return args.Get(0).(domain.Contender), args.Error(1)
}

func (m *repositoryMock) GetUserByUsername(ctx context.Context, tx domain.Transaction, username string) (domain.User, error) {
	args := m.Called(ctx, tx, username)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *repositoryMock) StoreUser(ctx context.Context, tx domain.Transaction, user domain.User) (domain.User, error) {
	args := m.Called(ctx, tx, user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *repositoryMock) StoreOrganizer(ctx context.Context, tx domain.Transaction, organizer domain.Organizer) (domain.Organizer, error) {
	args := m.Called(ctx, tx, organizer)
	return args.Get(0).(domain.Organizer), args.Error(1)
}

func (m *repositoryMock) AddUserToOrganizer(ctx context.Context, tx domain.Transaction, userID domain.UserID, organizerID domain.OrganizerID) error {
	args := m.Called(ctx, tx, userID, organizerID)
	return args.Error(0)
}

type jwtDecoderMock struct {
	mock.Mock
}

func (m *jwtDecoderMock) Decode(jwt string) (authorizer.Claims, error) {
	args := m.Called(jwt)
	return args.Get(0).(authorizer.Claims), args.Error(1)
}

type transactionMock struct {
	mock.Mock
}

func (m *transactionMock) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *transactionMock) Rollback() {
	m.Called()
}
