package postgres

import (
	"fmt"

	"github.com/SaidovZohid/competition-project/storage/repo"
	"github.com/jmoiron/sqlx"
)

type urlRepo struct {
	db *sqlx.DB
}

func NewUrl(db *sqlx.DB) repo.UrlStorageI {
	return &urlRepo{
		db: db,
	}
}

func (ur *urlRepo) Create(url *repo.Url) (*repo.Url, error) {
	query := `
		insert into urls(
			user_id,
			original_url,
			hashed_url,
			max_clicks
		) values ($1, $2, $3, $4)
		returning id, created_at
	`

	row := ur.db.QueryRow(
		query,
		url.UserId,
		url.OriginalUrl,
		url.HashedUrl,
		url.MaxClicks,
	)

	err := row.Scan(
		&url.Id,
		&url.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (ur *urlRepo) Get(id int64) (*repo.Url, error) {
	var result repo.Url

	query := `
		SELECT
			id,
			user_id,
			original_url,
			hashed_url,
			max_clicks,
			created_at
		FROM urls
		WHERE id=$1
	`

	row := ur.db.QueryRow(query, id)
	err := row.Scan(
		&result.Id,
		&result.UserId,
		&result.OriginalUrl,
		&result.HashedUrl,
		&result.MaxClicks,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *urlRepo) GetAll(params *repo.GetAllUrlsParams) (*repo.GetAllUrlsResult, error) {
	result := repo.GetAllUrlsResult{
		Urls: make([]*repo.Url, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" limit %d offset %d ", params.Limit, offset)

	filter := ""
	if params.Search != "" {
		str := "%" + params.Search + "%"
		filter += fmt.Sprintf(`
		where original_url ilike '%s' or hashed_url ilike '%s'	`,
			str, str,
		)
	}

	query := `
		SELECT
			id,
			user_id,
			original_url,
			hashed_url,
			max_clicks,
			created_at
		FROM urls
		` + filter + `
		ORDER BY created_at desc
		` + limit

	rows, err := ur.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u repo.Url

		err := rows.Scan(
			&u.Id,
			&u.UserId,
			&u.OriginalUrl,
			&u.HashedUrl,
			&u.MaxClicks,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result.Urls = append(result.Urls, &u)
	}

	queryCount := `SELECT count(1) FROM urls ` + filter
	err = ur.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *urlRepo) Update(url *repo.Url) (*repo.Url, error) {
	var result repo.Url

	query := `update users set
				user_id=$1,
				original_url=$2,
				hashed_url=$3,
				max_clicks=$4
			where id=$5
			returning id
			`

	row := ur.db.QueryRow(
		query,
		url.UserId,
		url.OriginalUrl,
		url.HashedUrl,
		url.MaxClicks,
		url.Id,
	)

	err := row.Scan(
		&result.Id,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *urlRepo) Delete(url *repo.Url) (*repo.Url, error) {
	var result repo.Url

	query := `delete from urls
			where id=$1
			returning id`

	row := ur.db.QueryRow(
		query,
		url.Id,
	)

	err := row.Scan(
		&result.Id,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}