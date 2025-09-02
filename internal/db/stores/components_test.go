package stores

import (
	"PartTrack/internal/db/models"
	"PartTrack/internal/db/views"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type MockComponentsStore struct{}

func (m *MockComponentsStore) GetOne(id uint64) (*models.Component, error) {
	return &models.Component{
		Id:   100,
		Name: "some random name",
	}, nil
}

func TestGetOneHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/protected/components/100", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/protected/components/:id")
	c.SetParamNames("id")
	c.SetParamValues("100")

	h := views.NewComponentsHandler()
	if assert.NoError(t, h.GetOne(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Test User")
	}
}
