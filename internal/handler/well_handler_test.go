// handler/well_handler_test.go
package handler_test

import (
	"bytes"
	"encoding/json"
	"gas_wells/internal/entity"
	"gas_wells/internal/handler"
	"gas_wells/internal/service/mocks"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateWell(t *testing.T) {
	mockService := new(mocks.WellService)
	handler := handler.NewWellHandler(mockService, slog.Default())

	well := &entity.Well{
		Name:        "Test Well",
		Pressure:    100,
		Temperature: 20,
	}

	mockService.On("CreateWell", mock.Anything, well).Return(well, nil)

	body, _ := json.Marshal(well)
	req := httptest.NewRequest("POST", "/api/v1/wells", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.CreateWell(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	mockService.AssertExpectations(t)
}
