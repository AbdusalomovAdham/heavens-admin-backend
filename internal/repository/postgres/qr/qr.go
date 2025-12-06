package qr

import (
	"context"

	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) Create(ctx context.Context, room_id, user_id int64, path string) (int64, string, error) {
	var id int64
	var oldPath string
	query := `SELECT id, path FROM qr_codes WHERE room_id = ?`

	rows, err := r.QueryContext(ctx, query, room_id)
	if err != nil {
		return 0, "", err
	}

	defer rows.Close()
	if err := r.ScanRows(ctx, rows, &id, &oldPath); err != nil {
		return 0, "", err
	}

	if id != 0 {
		query = `UPDATE qr_codes SET path = ? WHERE id = ?`
		_, err = r.ExecContext(ctx, query, path, id)
		return id, oldPath, err
	}

	query = `INSERT INTO qr_codes (room_id, path) VALUES (?, ?) RETURNING id`
	err = r.QueryRowContext(ctx, query, room_id, path).Scan(&id)
	if err != nil {
		return 0, "", err
	}

	return id, "", nil
}
