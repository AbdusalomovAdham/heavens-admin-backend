package departament

import (
	"context"
	"encoding/json"
	"fmt"
	"main/internal/entity"
	"main/internal/usecase/departament"
	"main/internal/usecase/user"
	"strings"

	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) Create(ctx context.Context, departament departament.Create, userId int64) (int64, error) {
	var id int64

	query := `INSERT INTO departaments (name, created_by, status) VALUES (?, ?, ?) RETURNING id`
	if err := r.QueryRowContext(ctx, query, departament.Name, userId, departament.Status).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) Delete(ctx context.Context, id int64, userId int64) error {

	query := `UPDATE departaments SET deleted_at = NOW(), deleted_by = ? WHERE id = ?  RETURNING id`

	_, err := r.ExecContext(ctx, query, userId, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetById(ctx context.Context, id int64) (departament.DepartamentById, error) {
	var detail departament.DepartamentById

	query := `SELECT id, status, created_at, name FROM departaments WHERE id = ?`

	rows, err := r.QueryContext(ctx, query, id)
	if err != nil {
		return departament.DepartamentById{}, err
	}
	defer rows.Close()

	if err := r.ScanRows(ctx, rows, &detail); err != nil {
		return departament.DepartamentById{}, err
	}
	return detail, nil
}

func (r *Repository) GetList(ctx context.Context, filter user.Filter, lang string) ([]departament.Get, int, error) {
	var list []departament.Get
	var limitQuery, offsetQuery string

	whereQuery := "WHERE d.deleted_at IS NULL"

	if filter.Limit != nil {
		limitQuery = fmt.Sprintf("LIMIT %d", *filter.Limit)
	}

	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf("OFFSET %d", *filter.Offset)
	}

	orderQuery := "ORDER BY d.id DESC"
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
	        d.id,
	        d.name->>'%s' as name,
	        d.status,
	        d.created_at
	    FROM departaments d
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

	countQuery := `SELECT COUNT(d.id) FROM departaments d`

	countRows, err := r.QueryContext(ctx, countQuery)
	if err != nil {
		return nil, 0, err
	}
	defer countRows.Close()

	count := 0

	if err = r.ScanRows(ctx, countRows, &count); err != nil {
		return nil, 0, fmt.Errorf("select user count: %w", err)
	}
	return list, count, nil
}

func (r *Repository) Update(ctx context.Context, id int64, data departament.Update, userId int64) error {
	var departament entity.Departament
	var nameJSON []byte
	query := `SELECT id, name, status, updated_at, updated_by FROM departaments WHERE id = ?`
	if err := r.QueryRowContext(ctx, query, id).Scan(&departament.Id, &nameJSON, &departament.Status, &departament.UpdatedAt, &departament.UpdatedBy); err != nil {
		return err
	}

	if err := json.Unmarshal(nameJSON, &departament.Name); err != nil {
		return err
	}

	if data.Name.Uz != nil {
		departament.Name.Uz = data.Name.Uz
	}

	if data.Name.Ru != nil {
		departament.Name.Ru = data.Name.Ru
	}

	if data.Name.En != nil {
		departament.Name.En = data.Name.En
	}

	if data.Status != nil {
		departament.Status = *data.Status
	}

	query = `UPDATE departaments SET name = ?, status = ?, updated_by = ?, updated_at = NOW() WHERE id = ?`
	if _, err := r.ExecContext(ctx, query, departament.Name, departament.Status, userId, id); err != nil {
		return err
	}
	return nil
}
