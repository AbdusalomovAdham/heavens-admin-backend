package role

import (
	"context"
	"encoding/json"
	"fmt"
	"main/internal/entity"
	"main/internal/usecase/role"
	"strings"

	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) Create(ctx context.Context, role role.Create, userId int64) (int64, error) {
	var id int64

	query := `INSERT INTO roles (name, created_by, status) VALUES (?, ?, ?) RETURNING id`
	if err := r.QueryRowContext(ctx, query, role.Name, userId, role.Status).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) Delete(ctx context.Context, id int64, userId int64) error {

	query := `UPDATE roles SET deleted_at = NOW(), deleted_by = ? WHERE id = ?  RETURNING id`

	_, err := r.ExecContext(ctx, query, userId, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetById(ctx context.Context, id int64) (role.RoleById, error) {
	var detail role.RoleById

	query := `SELECT id, status, created_at, name FROM roles WHERE id = ?`

	rows, err := r.QueryContext(ctx, query, id)
	if err != nil {
		return role.RoleById{}, err
	}
	defer rows.Close()

	if err := r.ScanRows(ctx, rows, &detail); err != nil {
		return role.RoleById{}, err
	}
	return detail, nil
}

func (r *Repository) GetList(ctx context.Context, filter entity.Filter, lang string) ([]role.Get, int, error) {
	var list []role.Get
	var limitQuery, offsetQuery string

	whereQuery := "WHERE r.deleted_at IS NULL"

	if filter.Limit != nil {
		limitQuery = fmt.Sprintf("LIMIT %d", *filter.Limit)
	}

	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf("OFFSET %d", *filter.Offset)
	}

	orderQuery := "ORDER BY r.id DESC"
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
	        r.id,
	        r.name->>'%s' as name,
	        r.status,
	        r.created_at
	    FROM roles r
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

	countQuery := `SELECT COUNT(r.id) FROM roles r WHERE r.deleted_at IS NULL`

	countRows, err := r.QueryContext(ctx, countQuery)
	if err != nil {
		return nil, 0, err
	}
	defer countRows.Close()

	count := 0

	if err = r.ScanRows(ctx, countRows, &count); err != nil {
		return nil, 0, fmt.Errorf("select role count: %w", err)
	}
	return list, count, nil
}

func (r *Repository) Update(ctx context.Context, id int64, data role.Update, userId int64) error {
	var role entity.Role
	var nameJSON []byte
	query := `SELECT id, name, status, updated_at, updated_by FROM roles WHERE id = ?`
	if err := r.QueryRowContext(ctx, query, id).Scan(&role.Id, &nameJSON, &role.Status, &role.UpdatedAt, &role.UpdatedBy); err != nil {
		return err
	}

	if err := json.Unmarshal(nameJSON, &role.Name); err != nil {
		return err
	}

	if data.Name.Uz != nil {
		role.Name.Uz = data.Name.Uz
	}

	if data.Name.Ru != nil {
		role.Name.Ru = data.Name.Ru
	}

	if data.Name.En != nil {
		role.Name.En = data.Name.En
	}

	if data.Status != nil {
		role.Status = *data.Status
	}

	query = `UPDATE roles SET name = ?, status = ?, updated_by = ?, updated_at = NOW() WHERE id = ?`
	if _, err := r.ExecContext(ctx, query, role.Name, role.Status, userId, id); err != nil {
		return err
	}
	return nil
}
