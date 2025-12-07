package position

import (
	"context"
	"encoding/json"
	"fmt"
	"main/internal/entity"

	"main/internal/usecase/position"
	"strings"

	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) Create(ctx context.Context, position position.Create, userId int64) (int64, error) {
	var id int64

	query := `INSERT INTO departments (name, created_by, status) VALUES (?, ?, ?) RETURNING id`
	if err := r.QueryRowContext(ctx, query, position.Name, userId, position.Status).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) Delete(ctx context.Context, id int64, userId int64) error {

	query := `UPDATE positions SET deleted_at = NOW(), deleted_by = ? WHERE id = ?  RETURNING id`

	_, err := r.ExecContext(ctx, query, userId, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetById(ctx context.Context, id int64) (position.PositionById, error) {
	var detail position.PositionById

	query := `SELECT id, status, created_at, name FROM positions WHERE id = ?`

	rows, err := r.QueryContext(ctx, query, id)
	if err != nil {
		return position.PositionById{}, err
	}
	defer rows.Close()

	if err := r.ScanRows(ctx, rows, &detail); err != nil {
		return position.PositionById{}, err
	}
	return detail, nil
}

func (r *Repository) GetList(ctx context.Context, filter entity.Filter, lang string) ([]position.Get, int, error) {
	var list []position.Get
	var limitQuery, offsetQuery string

	whereQuery := "WHERE p.deleted_at IS NULL"

	if filter.Limit != nil {
		limitQuery = fmt.Sprintf("LIMIT %d", *filter.Limit)
	}

	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf("OFFSET %d", *filter.Offset)
	}

	orderQuery := "ORDER BY p.id DESC"
	if filter.Order != nil && *filter.Order != "" {
		parts := strings.Fields(*filter.Order)
		if len(parts) == 2 {
			column := parts[0]
			direction := strings.ToUpper(parts[1])
			if direction != "ASC" && direction != "DESC" {
				direction = "ASC"
			}
			orderQuery = fmt.Sprintf("ORDER BY %s %s", column, direction)
		}
	}

	query := fmt.Sprintf(`
	    SELECT
	        p.id,
	        p.name->>'%s' as name,
	        p.status,
	        p.created_at
	    FROM positions p
	    %s %s %s %s
	`, lang, whereQuery, orderQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	if err := r.ScanRows(ctx, rows, &list); err != nil {
		return nil, 0, err
	}

	countQuery := `SELECT COUNT(p.id) FROM positions p`

	countRows, err := r.QueryContext(ctx, countQuery)
	if err != nil {
		return nil, 0, err
	}
	defer countRows.Close()

	count := 0

	if err = r.ScanRows(ctx, countRows, &count); err != nil {
		return nil, 0, fmt.Errorf("select position count: %w", err)
	}
	return list, count, nil
}

func (r *Repository) Update(ctx context.Context, id int64, data position.Update, userId int64) error {
	var position entity.Position
	var nameJSON []byte
	query := `SELECT id, name, status, updated_at, updated_by FROM positions WHERE id = ?`
	if err := r.QueryRowContext(ctx, query, id).Scan(&position.Id, &nameJSON, &position.Status, &position.UpdatedAt, &position.UpdatedBy); err != nil {
		return err
	}

	if err := json.Unmarshal(nameJSON, &position.Name); err != nil {
		return err
	}

	if data.Name.Uz != nil {
		position.Name.Uz = data.Name.Uz
	}

	if data.Name.Ru != nil {
		position.Name.Ru = data.Name.Ru
	}

	if data.Name.En != nil {
		position.Name.En = data.Name.En
	}

	if data.Status != nil {
		position.Status = *data.Status
	}

	query = `UPDATE positions SET name = ?, status = ?, updated_by = ?, updated_at = NOW() WHERE id = ?`
	if _, err := r.ExecContext(ctx, query, position.Name, position.Status, userId, id); err != nil {
		return err
	}
	return nil
}
