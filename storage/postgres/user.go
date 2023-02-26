package postgres

import (
	"database/sql"
	"fmt"

	"github.com/SaidovZohid/competition-project/pkg/utils"
	"github.com/SaidovZohid/competition-project/storage/repo"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) repo.UserStorageI {
	return &userRepo{
		db: db,
	}
}

func (ur *userRepo) Create(user *repo.User) (*repo.User, error) {
	query := `
		insert into users(
			first_name,
			last_name,
			email,
			password
		) values ($1, $2, $3, $4)
		returning id, created_at
	`

	err := ur.db.QueryRow(
		query,
		utils.NullString(user.FirstName),
		utils.NullString(user.LastName),
		user.Email,
		user.Password,
	).Scan(
		&user.Id,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *userRepo) Get(id int64) (*repo.User, error) {
	var result repo.User

	query := `
		SELECT
			id,
			first_name,
			last_name,
			email,
			created_at
		FROM users
		WHERE id=$1
	`
	var (
		firstName, lastName sql.NullString
	)
	row := ur.db.QueryRow(query, id)
	err := row.Scan(
		&result.Id,
		&firstName,
		&lastName,
		&result.Email,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	result.FirstName = firstName.String
	result.LastName = lastName.String

	return &result, nil
}

func (ur *userRepo) GetAll(params *repo.GetAllUsersParams) (*repo.GetAllUsersResult, error) {
	result := repo.GetAllUsersResult{
		Users: make([]*repo.User, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)

	filter := ""
	if params.Search != "" {
		str := "%" + params.Search + "%"
		filter += fmt.Sprintf(`
		WHERE first_name ILIKE '%s' OR last_name ILIKE '%s' OR email ILIKE '%s'	`,
			str, str, str,
		)
	}

	query := `
		SELECT
			id,
			first_name,
			last_name,
			email,
			created_at
		FROM users
		` + filter + `
		ORDER BY created_at desc
		` + limit

	rows, err := ur.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var (
		firstName, lastName sql.NullString
	)
	for rows.Next() {
		var u repo.User

		err := rows.Scan(
			&u.Id,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		u.FirstName = firstName.String
		u.LastName = lastName.String
		result.Users = append(result.Users, &u)
	}

	queryCount := `SELECT count(1) FROM users ` + filter
	err = ur.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *userRepo) GetByEmail(email string) (*repo.User, error) {
	var result repo.User

	query := `
		select
			id,
			first_name,
			last_name,
			email,
			password,
			created_at
		from users
		where email=$1
	`

	var (
		firstName, lastName sql.NullString
	)
	row := ur.db.QueryRow(query, email)
	err := row.Scan(
		&result.Id,
		&firstName,
		&lastName,
		&result.Email,
		&result.Password,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	result.FirstName = firstName.String
	result.LastName = lastName.String

	return &result, nil
}

func (ur *userRepo) UpdateUser(user *repo.User) (*repo.User, error) {
	var result repo.User

	query := `
		UPDATE users SET
			first_name=$1,
			last_name=$2
		WHERE id=$3
		RETURNING id, email, created_at
	`

	err := ur.db.QueryRow(
		query,
		user.FirstName,
		user.LastName,
		user.Id,
	).Scan(
		&result.Id,
		&result.Email,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *userRepo) DeleteUser(id int64) error {
	query := ` DELETE FROM users WHERE id=$1 `

	res, err := ur.db.Exec(
		query,
		id,
	)
	if err != nil {
		return err
	}
	if count, _ := res.RowsAffected(); count == 0 {
		return sql.ErrNoRows
	}
	return nil
}
