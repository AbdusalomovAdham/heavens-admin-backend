package attendancetemplate

import (
	"context"
	"encoding/json"
	"fmt"
	"main/internal/entity"
	attendancetemplate "main/internal/usecase/attendance_template"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r Repository) Create(ctx context.Context, data attendancetemplate.Create, startAtTime, finishAtTime *time.Time) (int64, error) {
	var id int64
	query := `INSERT INTO attendance_templates (name, status, created_at, created_by, start_at, finish_at, color, type_number) VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`
	if err := r.QueryRowContext(ctx, query, data.Name, data.Status, time.Now(), data.CreatedBy, startAtTime, finishAtTime, data.Color, data.TypeNumber).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) GetList(ctx context.Context, filter entity.Filter, lang string) ([]attendancetemplate.Get, int, error) {
	var list []attendancetemplate.Get
	var limitQuery, offsetQuery string

	whereQuery := "WHERE at.deleted_at IS NULL"

	if filter.Limit != nil {
		limitQuery = fmt.Sprintf("LIMIT %d", *filter.Limit)
	}

	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf("OFFSET %d", *filter.Offset)
	}

	orderQuery := "ORDER BY at.id DESC"
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
	        at.id,
	        at.name->>'%s' as name,
	        at.status,
	        at.created_at,
			at.finish_at,
			at.start_at,
			at.color,
			at.type_number
	    FROM attendance_templates at
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

	countQuery := `SELECT COUNT(at.id) FROM attendance_templates at WHERE at.deleted_at IS NULL`

	countRows, err := r.QueryContext(ctx, countQuery)
	if err != nil {
		return nil, 0, err
	}
	defer countRows.Close()

	count := 0

	if err = r.ScanRows(ctx, countRows, &count); err != nil {
		return nil, 0, fmt.Errorf("select attendance count: %w", err)
	}
	return list, count, nil
}

func (r *Repository) GetById(ctx context.Context, id int64) (attendancetemplate.AttendanceTemplateById, error) {
	var detail attendancetemplate.AttendanceTemplateById

	query := `SELECT id, status, created_at, finish_at, start_at, color, name, type_number FROM attendance_templates WHERE id = ? AND deleted_at IS NULL`

	rows, err := r.QueryContext(ctx, query, id)
	if err != nil {
		return attendancetemplate.AttendanceTemplateById{}, err
	}
	defer rows.Close()

	if err := r.ScanRows(ctx, rows, &detail); err != nil {
		return attendancetemplate.AttendanceTemplateById{}, err
	}
	return detail, nil
}

func (r *Repository) Update(ctx context.Context, id int64, data attendancetemplate.Update, userId int64, startAtTime, finishAtTime *time.Time) error {
	var attendancetemplate entity.AttendanceTemplate
	var nameJSON []byte
	query := `SELECT id, name, status, updated_at, updated_by, start_at, finish_at FROM attendance_templates WHERE id = ?`
	if err := r.QueryRowContext(ctx, query, id).Scan(&attendancetemplate.Id, &nameJSON, &attendancetemplate.Status, &attendancetemplate.UpdatedAt, &attendancetemplate.UpdatedBy, &attendancetemplate.StartAt, &attendancetemplate.FinishAt); err != nil {
		return err
	}

	if err := json.Unmarshal(nameJSON, &attendancetemplate.Name); err != nil {
		return err
	}

	if data.Name.Uz != nil {
		attendancetemplate.Name.Uz = data.Name.Uz
	}

	if data.Name.Ru != nil {
		attendancetemplate.Name.Ru = data.Name.Ru
	}

	if data.Name.En != nil {
		attendancetemplate.Name.En = data.Name.En
	}

	if data.Status != nil {
		attendancetemplate.Status = *data.Status
	}

	if startAtTime != nil {
		attendancetemplate.StartAt = *startAtTime
	}

	if finishAtTime != nil {
		attendancetemplate.FinishAt = *finishAtTime
	}

	query = `UPDATE attendance_templates SET name = ?, status = ?, start_at = ?, finish_at = ?, updated_by = ?, updated_at = NOW() WHERE id = ?`
	if _, err := r.ExecContext(ctx, query, attendancetemplate.Name, attendancetemplate.Status, attendancetemplate.StartAt, attendancetemplate.FinishAt, userId, id); err != nil {
		return err
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64, userId int64) error {
	query := `UPDATE attendance_templates SET deleted_by = ?, deleted_at = NOW() WHERE id = ?`
	if _, err := r.ExecContext(ctx, query, userId, id); err != nil {
		return err
	}
	return nil
}
