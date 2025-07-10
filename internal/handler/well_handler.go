// internal/handler/well_handler.go
package handler

import (
	"fmt"
	"gas_wells/internal/entity"
	"gas_wells/internal/pkg/logger"
	"gas_wells/internal/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type WellHandler struct {
	service *service.WellService
	logger  logger.Logger
	//templates *template.Template
}

func NewWellHandler(service *service.WellService, log logger.Logger) *WellHandler {
	return &WellHandler{
		service: service,
		logger:  log.With("layer", "handler"),
	}
}

func (h *WellHandler) RegisterRoutes(r chi.Router) {
	r.Get("/wells", h.ListWells)
	r.Get("/wells/create", h.CreateWellForm)
	r.Post("/wells", h.CreateWell)
	r.Get("/wells/{id}", h.GetWell)
	r.Get("/wells/{id}/edit", h.EditWellForm)
	r.Put("/wells/{id}", h.UpdateWell)
	r.Delete("/wells/{id}", h.DeleteWell)
}

// ListWells - отображает список всех скважин
func (h *WellHandler) ListWells(w http.ResponseWriter, r *http.Request) {
	wells, err := h.service.ListWells(r.Context())
	if err != nil {
		h.logger.Error("failed to list wells", "error", err)
		h.renderError(w, "Failed to load wells", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title": "Список скважин",
		"Wells": wells,
	}

	h.renderTemplate(w, "wells/list.html", data)
}

// CreateWellForm - форма создания новой скважины
func (h *WellHandler) CreateWellForm(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Новая скважина",
		"Well":  &entity.Well{}, // Пустая скважина для формы
	}

	h.renderTemplate(w, "wells/edit.html", data)
}

// CreateWell - обработчик создания скважины
func (h *WellHandler) CreateWell(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderError(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	well := &entity.Well{
		Name:        r.FormValue("name"),
		Pressure:    parseFloat(r.FormValue("pressure")),
		Temperature: parseFloat(r.FormValue("temperature")),
	}

	createdWell, err := h.service.CreateWell(r.Context(), well)
	if err != nil {
		h.logger.Error("failed to create well", "error", err)
		h.renderError(w, "Failed to create well", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/wells/%d", createdWell.ID), http.StatusSeeOther)
}

// GetWell - отображает детали скважины
func (h *WellHandler) GetWell(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.renderError(w, "Invalid well ID", http.StatusBadRequest)
		return
	}

	well, err := h.service.GetWell(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to get well", "id", id, "error", err)
		h.renderError(w, "Well not found", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Title": fmt.Sprintf("Скважина %s", well.Name),
		"Well":  well,
	}

	h.renderTemplate(w, "wells/view.html", data)
}

// EditWellForm - форма редактирования скважины
func (h *WellHandler) EditWellForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.renderError(w, "Invalid well ID", http.StatusBadRequest)
		return
	}

	well, err := h.service.GetWell(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to get well for edit", "id", id, "error", err)
		h.renderError(w, "Well not found", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Title": fmt.Sprintf("Редактирование %s", well.Name),
		"Well":  well,
	}

	h.renderTemplate(w, "wells/edit.html", data)
}

// UpdateWell - обработчик обновления скважины
func (h *WellHandler) UpdateWell(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.renderError(w, "Invalid well ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.renderError(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	well := &entity.Well{
		ID:          id,
		Name:        r.FormValue("name"),
		Pressure:    parseFloat(r.FormValue("pressure")),
		Temperature: parseFloat(r.FormValue("temperature")),
	}

	_, err = h.service.UpdateWell(r.Context(), well)
	if err != nil {
		h.logger.Error("failed to update well", "id", id, "error", err)
		h.renderError(w, "Failed to update well", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/wells/%d", id), http.StatusSeeOther)
}

// DeleteWell - обработчик удаления скважины
func (h *WellHandler) DeleteWell(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.renderError(w, "Invalid well ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteWell(r.Context(), id); err != nil {
		h.logger.Error("failed to delete well", "id", id, "error", err)
		h.renderError(w, "Failed to delete well", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/wells", http.StatusSeeOther)
}

// Вспомогательные методы

func (h *WellHandler) renderTemplate(w http.ResponseWriter, tmpl string, data map[string]interface{}) {
	w.Header().Set("Content-Type", "text/html")
	err := h.templates.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		h.logger.Error("failed to render template", "template", tmpl, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *WellHandler) renderError(w http.ResponseWriter, message string, statusCode int) {
	data := map[string]interface{}{
		"Error":      message,
		"StatusCode": statusCode,
	}
	w.WriteHeader(statusCode)
	h.renderTemplate(w, "error.html", data)
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
