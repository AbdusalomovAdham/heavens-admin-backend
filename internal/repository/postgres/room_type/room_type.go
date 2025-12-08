package room

import (
	"context"
	"fmt"
	"log"
	"main/internal/entity"
	roomtype "main/internal/usecase/room_type"
	"strings"

	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r Repository) Create(ctx context.Context, data roomtype.Create, userId int64) (int64, error) {
	var id int64
	query := `INSERT INTO room_types (name, status, created_by) VALUES (?, ?, ?) RETURNING id`

	err := r.QueryRowContext(ctx, query, data.Name, data.Status, userId).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r Repository) Update(ctx context.Context, data roomtype.Create, userId, id int64) (int64, error) {
	var detail entity.RoomType
	query := `SELECT id, name, status FROM room_types WHERE id = ?`

	rows, err := r.QueryContext(ctx, query, id)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if err := r.ScanRows(ctx, rows, &detail); err != nil {
		return 0, err
	}

	if data.Name != "" {
		detail.Name = data.Name
	}

	if userId != 0 {
		detail.CreatedBy = userId
	}

	if data.Status != nil {
		detail.Status = *data.Status
	}

	query = `UPDATE room_types SET name = ?, status = ?, created_by = ? WHERE id = ?`
	_, err = r.ExecContext(ctx, query, detail.Name, detail.Status, detail.CreatedBy, id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r Repository) GetList(ctx context.Context, filter *entity.Filter, userId int64) ([]entity.RoomType, uint32, error) {
	var list []entity.RoomType
	var limitQuery, offsetQuery string
	whereQuery := "WHERE r.deleted_at IS NULL AND r.deleted_at IS NULL"

	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(`LIMIT %d`, filter.Limit)
	}

	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(`OFFSET %d`, *filter.Offset)
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
			r.name,
			r.status,
			r.created_at
		FROM room_types r
		%s
		%s
		%s
		%s
		`, whereQuery, orderQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, err
	}
	countQuery := `SELECT COUNT(r.id) FROM room_types r WHERE r.deleted_at IS NULL`

	countRows, err := r.QueryContext(ctx, countQuery)
	if err != nil {
		return nil, 0, err
	}
	defer countRows.Close()

	count := 0

	if err = r.ScanRows(ctx, countRows, &count); err != nil {
		return nil, 0, fmt.Errorf("select room type count: %w", err)
	}

	formatcount := uint32(count)
	return list, formatcount, nil
}

func (r *Repository) GetById(ctx context.Context, id int64) (entity.RoomType, error) {
	var detail entity.RoomType

	query := fmt.Sprintf(`
		SELECT
			r.id,
			r.name,
			r.status,
			r.created_at
		FROM room_types r
		WHERE r.id = %d AND r.deleted_at IS NULL
		`, id)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return detail, err
	}
	defer rows.Close()

	err = r.ScanRows(ctx, rows, &detail)
	if err != nil {
		return detail, err
	}

	return detail, nil
}

func (r Repository) Delete(ctx context.Context, id, userID int64) error {
	log.Println("room id", id, "delete id", userID)
	query := `UPDATE room_types SET deleted_at = NOW(), deleted_by = ? WHERE id = ?`

	if _, err := r.ExecContext(ctx, query, userID, id); err != nil {
		return fmt.Errorf("delete room type; %w", err)
	}

	return nil
}
