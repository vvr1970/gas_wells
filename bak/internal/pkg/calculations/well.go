package calculations

import "errors"

// internal/pkg/calculations/well.go
func ProcessWellData(well entity.Well) (float64, error) {
	if well.Pressure <= 0 {
		return 0, errors.New("pressure must be positive")
	}
	return well.Pressure * well.Temperature * 0.85, nil
}
