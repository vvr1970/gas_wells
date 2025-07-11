// internal/service/well_service.go
package service

import (
	"context"
	"errors"
	"fmt"
	"gas_wells/internal/entity"
	"gas_wells/internal/pkg/logger"
	"gas_wells/internal/repository"
	"math"
)

type WellService struct {
	repo   repository.WellRepository
	logger logger.Logger
}

func NewWellService(repo repository.WellRepository, log logger.Logger) *WellService {
	return &WellService{
		repo:   repo,
		logger: log.With("layer", "service"),
	}
}

// CreateWell создает новую скважину с валидацией и расчетами
func (s *WellService) CreateWell(ctx context.Context, well *entity.Well) (*entity.Well, error) {
	// Валидация входных данных
	if well.Name == "" {
		return nil, errors.New("well name cannot be empty")
	}
	if well.Diameter <= 0 {
		return nil, errors.New("pressure must be positive")
	}
	if well.Temp <= -273.15 {
		return nil, errors.New("temperature cannot be below absolute zero")
	}

	// Выполняем расчеты
	result, err := s.calculateWellParameters(well.Pbuf, well.Temp)
	if err != nil {
		return nil, fmt.Errorf("calculation failed: %w", err)
	}
	well.Pmax = result

	// Сохраняем в БД
	if err := s.repo.Create(ctx, well); err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	return well, nil
}

// GetWell возвращает скважину по ID
func (s *WellService) GetWell(ctx context.Context, id int) (*entity.Well, error) {
	if id <= 0 {
		return nil, errors.New("invalid well ID")
	}

	well, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}
	if well == nil {
		return nil, errors.New("well not found")
	}

	return well, nil
}

// UpdateWell обновляет данные скважины
func (s *WellService) UpdateWell(ctx context.Context, well *entity.Well) (*entity.Well, error) {
	// Проверяем существование скважины
	existing, err := s.repo.GetByID(ctx, well.ID)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}
	if existing == nil {
		return nil, errors.New("well not found")
	}

	// Валидация
	if well.Pbuf <= 0 {
		return nil, errors.New("pressure must be positive")
	}

	// Пересчитываем параметры при изменении давления/температуры
	if well.Pbuf != existing.Pbuf || well.Temp != existing.Temp {
		result, err := s.calculateWellParameters(well.Pbuf, well.Temp)
		if err != nil {
			return nil, fmt.Errorf("calculation failed: %w", err)
		}
		well.Pmax = result
	}

	// Обновляем в БД
	if err := s.repo.Update(ctx, well); err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	return well, nil
}

// DeleteWell удаляет скважину
func (s *WellService) DeleteWell(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid well ID")
	}

	// Проверяем существование
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("repository error: %w", err)
	}

	return s.repo.Delete(ctx, id)
}

// ListWells возвращает список скважин с пагинацией
func (s *WellService) ListWells(ctx context.Context, limit int, offset int) ([]*entity.Well, error) {
	wells, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	return wells, nil
}

// calculateWellParameters содержит бизнес-логику расчетов
func (s *WellService) calculateWellParameters(pressure, temperature float64) (float64, error) {
	// Пример сложной бизнес-логики
	if pressure > 1000 {
		return 0, errors.New("pressure exceeds maximum allowed value")
	}

	// Формула расчета (пример)
	const efficiencyFactor = 0.85
	result := pressure * temperature * efficiencyFactor

	// Округляем до 2 знаков после запятой
	return math.Round(result*100) / 100, nil
}
