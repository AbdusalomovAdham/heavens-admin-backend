package user

import (
	"context"
	"fmt"
	"log"
	"main/internal/entity"
	"main/internal/usecase/user"
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

func (r *Repository) GetByLogin(ctx context.Context, login string) (entity.User, error) {

	query := `SELECT * FROM users WHERE login = ?`

	rows, err := r.QueryContext(ctx, query, login)
	if err != nil {
		return entity.User{}, err
	}

	defer rows.Close()
	user := entity.User{}
	if err := r.ScanRows(ctx, rows, &user); err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *Repository) GetById(ctx context.Context, id int64) (entity.User, error) {

	query := `SELECT * FROM users WHERE id = ?`
	rows, err := r.QueryContext(ctx, query, id)
	if err != nil {
		return entity.User{}, err
	}

	defer rows.Close()
	user := entity.User{}
	if err := r.ScanRows(ctx, rows, &user); err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *Repository) Create(ctx context.Context, user user.Create) error {
	query := `
		INSERT INTO users (
		    avatar, first_name, last_name, middle_name,
		    region_id, district_id, gender_id, birth_date,
		    login, password, email, phone_number, mobile_number,
		    role, management_id, position_id, work_status_id,
		    salary, passport_number, jshshir, passport_scan,
		    car_prefix, car_number, diploma_file, cv_file,
		    status, created_by
		) VALUES (
		    ?, ?, ?, ?,
		    ?, ?, ?, ?,
		    ?, ?, ?, ?, ?,
		    ?, ?, ?, ?,
		    ?, ?, ?, ?,
		    ?, ?, ?, ?,
		    ?, ?
		)
		RETURNING id
		`
	log.Println("region id", user.RegionId)
	var id int64
	err := r.QueryRowContext(ctx, query,
		user.Avatar, user.FirstName, user.LastName, user.MiddleName,
		user.RegionId, user.DistrictId, user.GenderId, user.BirthDate,
		user.Login, user.Password, user.Email, user.PhoneNumber, user.MobileNumber,
		user.Role, user.ManagementId, user.PositionId, user.WorkStatusId,
		user.Salary, user.PassportNumber, user.JSHSHIR, user.PassportScan,
		user.CarPrefix, user.CarNumber, user.DimlomaFile, user.CVFile,
		user.Status, user.CreatedBy,
	).Scan(&id)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id, deletedBy int64) error {
	query := `UPDATE users SET deleted_at = NOW(), deleted_by = ? WHERE id = ?`

	if _, err := r.ExecContext(ctx, query, deletedBy, id); err != nil {
		return fmt.Errorf("delete user; %w", err)
	}

	return nil
}

func (r *Repository) GetList(ctx context.Context, filter entity.Filter) ([]user.UserPreview, int, error) {
	var list []user.UserPreview
	var limitQuery, offsetQuery string

	whereQuery := "WHERE u.deleted_at IS NULL"

	if filter.Limit != nil {
		limitQuery = fmt.Sprintf("LIMIT %d", *filter.Limit)
	}

	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf("OFFSET %d", *filter.Offset)
	}

	orderQuery := "ORDER BY u.id DESC"

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
        u.id,
        u.last_name,
        u.first_name,
        u.avatar,
        u.management_id,
        u.position_id,
        u.status
    FROM users u
    %s
    %s
    %s
    %s`,
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

	countQuery := `SELECT COUNT(u.id) FROM users u`

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

func (r *Repository) Update(ctx context.Context, id int64, data user.Update) (int64, error) {
	var detail entity.User
	fmt.Print("detail", data.LastName)
	query := fmt.Sprintf(`SELECT * FROM users WHERE id = %d`, id)
	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if err := r.ScanRows(ctx, rows, &detail); err != nil {
		return 0, err
	}

	setParts := []string{}
	args := []any{}
	i := 1

	if data.Avatar != nil {
		setParts = append(setParts, fmt.Sprintf("avatar = %s", "?"))
		args = append(args, *data.Avatar)
		i++
	}

	if data.FirstName != nil {
		setParts = append(setParts, fmt.Sprintf("first_name = %s", "?"))
		args = append(args, *data.FirstName)
		i++
	}

	if data.LastName != nil {
		setParts = append(setParts, fmt.Sprintf("last_name = %s", "?"))
		args = append(args, *data.LastName)
		i++
	}

	if data.MiddleName != nil {
		setParts = append(setParts, fmt.Sprintf("middle_name = %s", "?"))
		args = append(args, *data.MiddleName)
		i++
	}

	if data.RegionId != nil {
		setParts = append(setParts, fmt.Sprintf("region_id = %s", "?"))
		args = append(args, *data.RegionId)
		i++
	}

	if data.DistrictId != nil {
		setParts = append(setParts, fmt.Sprintf("district_id = %s", "?"))
		args = append(args, *data.DistrictId)
		i++
	}

	if data.GenderId != nil {
		setParts = append(setParts, fmt.Sprintf("gender_id = %s", "?"))
		args = append(args, *data.GenderId)
		i++
	}

	if data.BirthDate != nil {
		setParts = append(setParts, fmt.Sprintf("birth_date = %s", "?"))
		args = append(args, *data.BirthDate)
		i++
	}

	if data.Login != nil {
		setParts = append(setParts, fmt.Sprintf("login = %s", "?"))
		args = append(args, *data.Login)
		i++
	}

	if data.Password != nil {
		setParts = append(setParts, fmt.Sprintf("password = %s", "?"))
		args = append(args, *data.Password)
		i++
	}

	if data.Email != nil {
		setParts = append(setParts, fmt.Sprintf("email = %s", "?"))
		args = append(args, *data.Email)
		i++
	}

	if data.PhoneNumber != nil {
		setParts = append(setParts, fmt.Sprintf("phone_number = %s", "?"))
		args = append(args, *data.PhoneNumber)
	}

	if data.MobileNumber != nil {
		setParts = append(setParts, fmt.Sprintf("mobile_number = %s", "?"))
		args = append(args, *data.MobileNumber)
		i++
	}

	if data.Role != nil {
		setParts = append(setParts, fmt.Sprintf("role = %s", "?"))
		args = append(args, *data.Role)
		i++
	}

	if data.ManagementId != nil {
		setParts = append(setParts, fmt.Sprintf("management_id = %s", "?"))
		args = append(args, *data.ManagementId)
		i++
	}

	if data.PositionId != nil {
		setParts = append(setParts, fmt.Sprintf("position_id = %s", "?"))
		args = append(args, *data.PositionId)
		i++
	}

	if data.WorkStatusId != nil {
		setParts = append(setParts, fmt.Sprintf("work_status_id = %s", "?"))
		args = append(args, *data.WorkStatusId)
		i++
	}

	if data.Salary != nil {
		setParts = append(setParts, fmt.Sprintf("salary = %s", "?"))
		args = append(args, *data.Salary)
		i++
	}

	if data.PassportNumber != nil {
		setParts = append(setParts, fmt.Sprintf("passport_number = %s", "?"))
		args = append(args, *data.PassportNumber)
		i++
	}

	if data.JSHSHIR != nil {
		setParts = append(setParts, fmt.Sprintf("jshshir = %s", "?"))
		args = append(args, *data.JSHSHIR)
	}

	if data.PassportScan != nil {
		setParts = append(setParts, fmt.Sprintf("passport_scan = %s", "?"))
		args = append(args, *data.PassportScan)
		i++
	}

	if data.CarPrefix != nil {
		setParts = append(setParts, fmt.Sprintf("car_prefix = %s", "?"))
		args = append(args, *data.CarPrefix)
		i++
	}

	if data.CarNumber != nil {
		setParts = append(setParts, fmt.Sprintf("car_number = %s", "?"))
		args = append(args, *data.CarNumber)
		i++
	}

	if data.DimlomaFile != nil {
		setParts = append(setParts, fmt.Sprintf("diploma_file = %s", "?"))
		args = append(args, *data.DimlomaFile)
		i++
	}

	if data.CVFile != nil {
		setParts = append(setParts, fmt.Sprintf("cv_file = %s", "?"))
		args = append(args, *data.CVFile)
		i++
	}

	if data.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = %s", "?"))
		args = append(args, *data.Status)
		i++
	}

	if data.UpdatedBy != nil {
		setParts = append(setParts, fmt.Sprintf("updated_by = %s", "?"))
		args = append(args, *data.UpdatedBy)
		i++
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = %s", "?"))
	args = append(args, time.Now())
	i++

	if len(setParts) == 0 {
		return 0, fmt.Errorf("nothing to update")
	}

	updateQuery := fmt.Sprintf("UPDATE users SET %s WHERE id = %d",
		strings.Join(setParts, ", "),
		id,
	)
	args = append(args, id)

	_, err = r.ExecContext(ctx, updateQuery, args...)
	if err != nil {
		return 0, err
	}

	return detail.Id, nil
}
