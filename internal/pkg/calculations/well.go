package calculations

import (
	"errors"
	"gas_wells/internal/entity"
)

// internal/pkg/calculations/well.go
func ProcessWellData(well entity.Well) (float64, error) {
	if well.Pbuf <= 0 {
		return 0, errors.New("pressure must be positive")
	}
	return well.Pbuf * well.Temp * 0.85, nil
}
