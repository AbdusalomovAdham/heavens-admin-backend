package room

import (
	"context"
	"fmt"
	"log"
	"main/internal/entity"
	"main/internal/usecase/room"
	"strings"

	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r Repository) Create(ctx context.Context, roomData room.Create, userId int64) (int64, error) {
	var id int64
	query := `INSERT INTO rooms (room_number, room_type, employee_id, created_by, corpus) VALUES (?, ?, ?, ?, ?) RETURNING id`

	err := r.QueryRowContext(ctx, query, roomData.RoomNumber, roomData.RoomType, roomData.EmployeeId, userId, roomData.Corpus).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r Repository) Delete(ctx context.Context, id, userID int64) error {
	log.Println("room id", id, "delete id", userID)
	query := `UPDATE rooms SET deleted_at = NOW(), deleted_by = ? WHERE id = ?`

	if _, err := r.ExecContext(ctx, query, userID, id); err != nil {
		return fmt.Errorf("delete room; %w", err)
	}

	return nil
}

func (r Repository) GetList(ctx context.Context, filter *entity.Filter) ([]room.RoomPreview, uint32, error) {
	var list []room.RoomPreview

	var limitQuery, offsetQuery string
	whereQuery := "WHERE r.deleted_at IS NULL AND u.deleted_at IS NULL"

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

	query := fmt.Sprintf(
		`SELECT
		        r.id,
		        r.employee_id,
		        r.room_number,
		        r.room_type,
		        r.created_at,
				r.corpus,
		        r.status,
				u.first_name,
				u.last_name,
				u.middle_name,
				q.path
		FROM rooms r
		LEFT JOIN users u ON r.employee_id = u.id
		LEFT JOIN qr_codes q ON q.room_id = r.id
		%s %s %s %s`,
		whereQuery,
		orderQuery,
		limitQuery,
		offsetQuery,
	)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	if err := r.ScanRows(ctx, rows, &list); err != nil {
		return nil, 0, err
	}

	countQuery := `SELECT COUNT(r.id) FROM rooms r WHERE r.deleted_at IS NULL`

	countRows, err := r.QueryContext(ctx, countQuery)
	if err != nil {
		return nil, 0, err
	}
	defer countRows.Close()

	count := 0

	if err = r.ScanRows(ctx, countRows, &count); err != nil {
		return nil, 0, fmt.Errorf("select user count: %w", err)
	}

	formatount := uint32(count)
	return list, formatount, nil
}

func (r Repository) GetById(ctx context.Context, id int64) (room.RoomPreview, error) {
	var result room.RoomPreview

	query := `SELECT
			    r.id,
			    r.employee_id,
			    r.room_number,
			    r.room_type,
			    r.created_at,
			    r.status,
			    r.corpus,
			    u.first_name,
			    u.last_name,
			    u.middle_name,
				q.path
			FROM rooms r
			LEFT JOIN users u ON r.employee_id = u.id
			LEFT JOIN qr_codes q ON q.room_id = r.id
			WHERE u.deleted_at IS NULL AND r.id = ? AND r.deleted_at IS NULL
`
	rows, err := r.QueryContext(ctx, query, id)
	if err != nil {
		return room.RoomPreview{}, err
	}
	defer rows.Close()

	if err := r.ScanRows(ctx, rows, &result); err != nil {
		return room.RoomPreview{}, err
	}

	return result, nil
}
