package postgres

import (
	"fmt"

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

	row := ur.db.QueryRow(
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
	)

	err := row.Scan(
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
			password,
			created_at
		FROM users
		WHERE id=$1
	`

	row := ur.db.QueryRow(query, id)
	err := row.Scan(
		&result.Id,
		&result.FirstName,
		&result.LastName,
		&result.Email,
		&result.Password,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *userRepo) GetAll(params *repo.GetAllUsersParams) (*repo.GetAllUsersResult, error) {
	result := repo.GetAllUsersResult{
		Users: make([]*repo.User, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" limit %d offset %d ", params.Limit, offset)

	filter := ""
	if params.Search != "" {
		str := "%" + params.Search + "%"
		filter += fmt.Sprintf(`
		where first_name ilike '%s' or last_name ilike '%s' or email ilike '%s'	`,
			str, str, str,
		)
	}

	query := `
		SELECT
			id,
			first_name,
			last_name,
			email,
			password,
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

	for rows.Next() {
		var u repo.User

		err := rows.Scan(
			&u.Id,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.Password,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

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

	row := ur.db.QueryRow(query, email)
	err := row.Scan(
		&result.Id,
		&result.FirstName,
		&result.LastName,
		&result.Email,
		&result.Password,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *userRepo) UpdatePassword(req *repo.UpdatePassword) error {
	query := ` UPDATE users SET password=$1 WHERE id=$2 `

	_, err := ur.db.Exec(query, req.Password, req.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (ur *userRepo) UpdateUser(user *repo.User) (*repo.User, error) {
	var result repo.User

	query := `update users set
				first_name=$1,
				last_name=$2,
				email=$3,
				password=$4,
			where id=$5
			returning id
			`

	row := ur.db.QueryRow(
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.Id,
	)

	err := row.Scan(
		&result.Id,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *userRepo) DeleteUser(user *repo.User) (*repo.User, error) {
	var result repo.User

	query := `delete from users
			where id=$1
			returning id`

	row := ur.db.QueryRow(
		query,
		user.Id,
	)

	err := row.Scan(
		&result.Id,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
