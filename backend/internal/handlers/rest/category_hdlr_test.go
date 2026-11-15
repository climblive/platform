package rest_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/handlers/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoryHandler(t *testing.T) {
	category := domain.Category{
		ID:            42,
		CompetitionID: 123,
		Name:          "Open",
		TimeBegin:     time.Date(2026, time.November, 15, 10, 0, 0, 0, time.UTC),
		TimeEnd:       time.Date(2026, time.November, 15, 16, 0, 0, 0, time.UTC),
	}
	const categoryJSON = `{"id":42,"competitionId":123,"name":"Open","timeBegin":"2026-11-15T10:00:00Z","timeEnd":"2026-11-15T16:00:00Z"}`
	const templateJSON = `{"name":"Open","timeBegin":"2026-11-15T10:00:00Z","timeEnd":"2026-11-15T16:00:00Z"}`

	for _, tc := range []struct {
		name   string
		method string
		path   string
		body   string
		status int
		result string
		setup  func(*categoryUseCaseMock)
	}{
		{
			name: "Get", method: http.MethodGet, path: "/categories/42",
			status: http.StatusOK, result: categoryJSON,
			setup: func(uc *categoryUseCaseMock) {
				uc.On("GetCategory", mock.Anything, domain.CategoryID(42)).Return(category, nil)
			},
		},
		{
			name: "List", method: http.MethodGet, path: "/competitions/123/categories",
			status: http.StatusOK, result: "[" + categoryJSON + "]",
			setup: func(uc *categoryUseCaseMock) {
				uc.On("GetCategoriesByCompetition", mock.Anything, domain.CompetitionID(123)).Return([]domain.Category{category}, nil)
			},
		},
		{
			name: "Create", method: http.MethodPost, path: "/competitions/123/categories", body: templateJSON,
			status: http.StatusCreated, result: categoryJSON,
			setup: func(uc *categoryUseCaseMock) {
				uc.On("CreateCategory", mock.Anything, domain.CompetitionID(123), domain.CategoryTemplate{
					Name: category.Name, TimeBegin: category.TimeBegin, TimeEnd: category.TimeEnd,
				}).Return(category, nil)
			},
		},
		{
			name: "Patch", method: http.MethodPatch, path: "/categories/42", body: `{"name":"Open"}`,
			status: http.StatusOK, result: categoryJSON,
			setup: func(uc *categoryUseCaseMock) {
				uc.On("PatchCategory", mock.Anything, domain.CategoryID(42), domain.CategoryPatch{
					Name: domain.NewPatch("Open"),
				}).Return(category, nil)
			},
		},
		{
			name: "Delete", method: http.MethodDelete, path: "/categories/42",
			status: http.StatusNoContent,
			setup: func(uc *categoryUseCaseMock) {
				uc.On("DeleteCategory", mock.Anything, domain.CategoryID(42)).Return(nil)
			},
		},
		{
			name: "InvalidCategoryID", method: http.MethodGet, path: "/categories/invalid",
			status: http.StatusBadRequest,
		},
		{
			name: "InvalidCompetitionID", method: http.MethodGet, path: "/competitions/invalid/categories",
			status: http.StatusBadRequest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			uc := new(categoryUseCaseMock)
			if tc.setup != nil {
				tc.setup(uc)
			}
			mux := rest.NewMux()
			rest.InstallCategoryHandler(mux, uc)

			request := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			mux.ServeHTTP(response, request)

			assert.Equal(t, tc.status, response.Code)
			if tc.result != "" {
				assert.JSONEq(t, tc.result, response.Body.String())
			} else {
				assert.Empty(t, response.Body.String())
			}
			uc.AssertExpectations(t)
		})
	}
}

type categoryUseCaseMock struct {
	mock.Mock
}

func (m *categoryUseCaseMock) GetCategory(ctx context.Context, categoryID domain.CategoryID) (domain.Category, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).(domain.Category), args.Error(1)
}

func (m *categoryUseCaseMock) GetCategoriesByCompetition(ctx context.Context, competitionID domain.CompetitionID) ([]domain.Category, error) {
	args := m.Called(ctx, competitionID)
	return args.Get(0).([]domain.Category), args.Error(1)
}

func (m *categoryUseCaseMock) CreateCategory(ctx context.Context, competitionID domain.CompetitionID, template domain.CategoryTemplate) (domain.Category, error) {
	args := m.Called(ctx, competitionID, template)
	return args.Get(0).(domain.Category), args.Error(1)
}

func (m *categoryUseCaseMock) PatchCategory(ctx context.Context, categoryID domain.CategoryID, patch domain.CategoryPatch) (domain.Category, error) {
	args := m.Called(ctx, categoryID, patch)
	return args.Get(0).(domain.Category), args.Error(1)
}

func (m *categoryUseCaseMock) DeleteCategory(ctx context.Context, categoryID domain.CategoryID) error {
	return m.Called(ctx, categoryID).Error(0)
}
