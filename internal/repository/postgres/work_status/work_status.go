package workstatuses

import (
	"context"
	"encoding/json"
	"fmt"
	"main/internal/entity"
	workstatus "main/internal/usecase/work_status"
	"strings"

	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) Create(ctx context.Context, workStatus workstatus.Create, userId int64) (int64, error) {
	var id int64

	query := `INSERT INTO work_statuses (name, created_by, status) VALUES (?, ?, ?) RETURNING id`
	if err := r.QueryRowContext(ctx, query, workStatus.Name, userId, workStatus.Status).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) Delete(ctx context.Context, id int64, userId int64) error {

	query := `UPDATE work_statuses SET deleted_at = NOW(), deleted_by = ? WHERE id = ?  RETURNING id`

	_, err := r.ExecContext(ctx, query, userId, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetById(ctx context.Context, id int64) (workstatus.WorkStatusById, error) {
	var detail workstatus.WorkStatusById

	query := `SELECT id, status, created_at, name FROM work_statuses WHERE id = ?`

	rows, err := r.QueryContext(ctx, query, id)
	if err != nil {
		return workstatus.WorkStatusById{}, err
	}
	defer rows.Close()

	if err := r.ScanRows(ctx, rows, &detail); err != nil {
		return workstatus.WorkStatusById{}, err
	}
	return detail, nil
}

func (r *Repository) GetList(ctx context.Context, filter entity.Filter, lang string) ([]workstatus.Get, int, error) {
	var list []workstatus.Get
	var limitQuery, offsetQuery string

	whereQuery := "WHERE ws.deleted_at IS NULL"

	if filter.Limit != nil {
		limitQuery = fmt.Sprintf("LIMIT %d", *filter.Limit)
	}

	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf("OFFSET %d", *filter.Offset)
	}

	orderQuery := "ORDER BY ws.id DESC"
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
	        ws.id,
	        ws.name->>'%s' as name,
	        ws.status,
	        ws.created_at
	    FROM work_statuses ws
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

	countQuery := `SELECT COUNT(ws.id) FROM work_statuses ws`

	countRows, err := r.QueryContext(ctx, countQuery)
	if err != nil {
		return nil, 0, err
	}
	defer countRows.Close()

	count := 0

	if err = r.ScanRows(ctx, countRows, &count); err != nil {
		return nil, 0, fmt.Errorf("select work status count: %w", err)
	}
	return list, count, nil
}

func (r *Repository) Update(ctx context.Context, id int64, data workstatus.Update, userId int64) error {
	var workStatus entity.WorkStatus
	var nameJSON []byte
	query := `SELECT id, name, status, updated_at, updated_by FROM work_statuses WHERE id = ?`
	if err := r.QueryRowContext(ctx, query, id).Scan(&workStatus.Id, &nameJSON, &workStatus.Status, &workStatus.UpdatedAt, &workStatus.UpdatedBy); err != nil {
		return err
	}

	if err := json.Unmarshal(nameJSON, &workStatus.Name); err != nil {
		return err
	}

	if data.Name.Uz != nil {
		workStatus.Name.Uz = data.Name.Uz
	}

	if data.Name.Ru != nil {
		workStatus.Name.Ru = data.Name.Ru
	}

	if data.Name.En != nil {
		workStatus.Name.En = data.Name.En
	}

	if data.Status != nil {
		workStatus.Status = *data.Status
	}

	query = `UPDATE work_statuses SET name = ?, status = ?, updated_by = ?, updated_at = NOW() WHERE id = ?`
	if _, err := r.ExecContext(ctx, query, workStatus.Name, workStatus.Status, userId, id); err != nil {
		return err
	}
	return nil
}
