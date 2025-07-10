// internal/repository/repository.go
package repository

import (
	"context"
	"gas_wells/internal/entity"
)

//Сначала создаем интерфейсы для репозиториев.
//Это гарантирует, что сервисы не зависят от конкретной БД.

type WellRepository interface {
	Create(ctx context.Context, well *entity.Well) error
	GetByID(ctx context.Context, id int) (*entity.Well, error)
	Update(ctx context.Context, well *entity.Well) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]*entity.Well, error)
}

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}
