// internal/repository/well_repo.go
package repository

import (
	"context"
	"errors"
	"gas_wells/internal/entity"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WellRepo struct {
	db *pgxpool.Pool // Используем пул соединений
}

func NewWellRepo(db *pgxpool.Pool) *WellRepo {
	return &WellRepo{db: db}
}

// Create - добавление новой скважины
func (r *WellRepo) Create(ctx context.Context, well *entity.Well) error {
	query := `
		INSERT INTO wells (name, location, gammag, temp, tempust, depth,
					pbuf, ptb, ppl, pz, q, roughness, diametr, a, b, mu,
					wgf, rog, hw, qmin, pmax, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8,$9, $10, $11, $12, $13, $14,
				$15, $16, $17, $18, $19, $20, $21, $22) 
		RETURNING id, created_at
	`
	err := r.db.QueryRow(
		ctx,
		query,
		well.Name,
		well.Location,
		well.GammaG,
		well.Temp,
		well.TempUst,
		well.Depth,
		well.Pbuf,
		well.Ptb,
		well.Ppl,
		well.Pz,
		well.Q,
		well.Roughness,
		well.Diameter,
		well.A,
		well.B,
		well.Mu,
		well.WGF,
		well.Rog,
		well.Hw,
		well.Qmin,
		well.Pmax,
		well.Status,
	).Scan(&well.ID, &well.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// GetByID - получение скважины по ID
func (r *WellRepo) GetByID(ctx context.Context, id int) (*entity.Well, error) {
	query := `SELECT * FROM wells WHERE id = $1`
	well := &entity.Well{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&well.Name,
		&well.Location,
		&well.GammaG,
		&well.Temp,
		&well.TempUst,
		&well.Depth,
		&well.Pbuf,
		&well.Ptb,
		&well.Ppl,
		&well.Pz,
		&well.Q,
		&well.Roughness,
		&well.Diameter,
		&well.A,
		&well.B,
		&well.Mu,
		&well.WGF,
		&well.Rog,
		&well.Hw,
		&well.Qmin,
		&well.Pmax,
		&well.Status,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return well, nil
}

// Update - обновление данных скважины
func (r *WellRepo) Update(ctx context.Context, well *entity.Well) error {
	query := `
	UPDATE wells 
	SET name=$1, location=$2, gammag=$3, temp=$4, tempust=$5, 
		depth=$6, pbuf=$7, ptb=$8, ppl=$9, pz=$10, q=$11, roughness=$12,
		diametr=$13, a=$14, b=$15, mu=$16, wgf=$17,	rog=$18, hw=$19,
		qmin=$20, pmax=$21, status=$22, updated_at = NOW() 
		WHERE id = $23	
	`
	result, err := r.db.Exec(
		ctx,
		query,
		well.Name,
		well.Location,
		well.GammaG,
		well.Temp,
		well.TempUst,
		well.Depth,
		well.Pbuf,
		well.Ptb,
		well.Ppl,
		well.Pz,
		well.Q,
		well.Roughness,
		well.Diameter,
		well.A,
		well.B,
		well.Mu,
		well.WGF,
		well.Rog,
		well.Hw,
		well.Qmin,
		well.Pmax,
		well.Status,
		well.ID,
	)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no rows affected (well not found)")
	}
	return nil
}

// Delete - удаление скважины по ID
func (r *WellRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM wells WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		slog.Error("failed to delete well", "id", id, "error", err)
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("well not found")
	}
	return nil
}

// List - получение списка всех скважин (с пагинацией)
func (r *WellRepo) List(ctx context.Context, limit, offset int) ([]*entity.Well, error) {
	query := `
		SELECT id, name, location, pbuf, status, result, created_at, updated_at
		FROM wells 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wells []*entity.Well
	for rows.Next() {
		well := &entity.Well{}
		err := rows.Scan(
			&well.ID,
			&well.Name,
			&well.Location,
			&well.Pbuf,
			&well.Status,
			&well.Result,
			&well.CreatedAt,
			&well.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		wells = append(wells, well)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return wells, nil
}
