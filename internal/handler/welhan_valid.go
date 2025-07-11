package handler

//Пример использования в обработчике (internal/handler/well_handler.go)

/*func (h *WellHandler) CreateWell(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string  `json:"name"`
		Pressure    float64 `json:"pressure"`
		Temperature float64 `json:"temperature"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.renderError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация
	v := validation.New()
	v.Check(validation.NotBlank(input.Name), "name", "Name is required")
	v.Check(validation.Between(input.Pressure, 0.1, 1000.0), "pressure", "Pressure must be between 0.1 and 1000")
	v.Check(validation.Between(input.Temperature, -50.0, 150.0), "temperature", "Temperature must be between -50 and 150")

	if !v.Valid() {
		h.renderValidationErrors(w, v.Errors)
		return
	}

	// Обработка запроса
	// ...
}

func (h *WellHandler) renderValidationErrors(w http.ResponseWriter, errors map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"errors": errors,
	})
} */
