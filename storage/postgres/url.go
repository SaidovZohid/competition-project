package postgres

import (
	"database/sql"
	"fmt"

	"github.com/SaidovZohid/competition-project/pkg/utils"
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
			max_clicks,
			expires_at
		) values ($1, $2, $3, $4, $5)
		returning id, created_at
	`

	row := ur.db.QueryRow(
		query,
		url.UserId,
		url.OriginalUrl,
		url.HashedUrl,
		utils.NullInt64(url.MaxClicks),
		url.ExpiresAt,
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

func (ur *urlRepo) Get(url string) (*repo.Url, error) {
	var result repo.Url

	query := fmt.Sprintf(`
		SELECT
			id,
			user_id,
			original_url,
			hashed_url,
			max_clicks,
			expires_at,
			created_at
		FROM urls
		WHERE hashed_url='%s'
	`, url)

	var maxClicks sql.NullInt64
	err := ur.db.QueryRow(query).Scan(
		&result.Id,
		&result.UserId,
		&result.OriginalUrl,
		&result.HashedUrl,
		&maxClicks,
		&result.ExpiresAt,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	result.MaxClicks = maxClicks.Int64

	return &result, nil
}

func (ur *urlRepo) GetAll(params *repo.GetAllUrlsParams) (*repo.GetAllUrlsResult, error) {
	result := repo.GetAllUrlsResult{
		Urls: make([]*repo.Url, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" limit %d offset %d ", params.Limit, offset)

	filter := " where true "
	if params.Search != "" {
		str := "%" + params.Search + "%"
		filter += fmt.Sprintf(`
		AND original_url ilike '%s' or hashed_url ilike '%s'	`,
			str, str,
		)
	}
	if params.UserID != 0 {
		filter += fmt.Sprintf(" AND user_id = %d ", params.UserID)
	}

	query := `
		SELECT
			id,
			user_id,
			original_url,
			hashed_url,
			max_clicks,
			expires_at,
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
		var (
			u         repo.Url
			maxClicks sql.NullInt64
		)
		err := rows.Scan(
			&u.Id,
			&u.UserId,
			&u.OriginalUrl,
			&u.HashedUrl,
			&maxClicks,
			&u.ExpiresAt,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		u.MaxClicks = maxClicks.Int64

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

	query := `
		update urls set
			hashed_url=$1,
			max_clicks=$2,
			expires_at=$3
		where id=$4 and user_id=$5
		returning 
			id, 
			user_id, 
			original_url,
			hashed_url,
			max_clicks,
			expires_at,
			created_at
	`
	var maxClicks sql.NullInt64
	err := ur.db.QueryRow(
		query,
		url.HashedUrl,
		utils.NullInt64(url.MaxClicks),
		url.ExpiresAt,
		url.Id,
		url.UserId,
	).Scan(
		&result.Id,
		&result.UserId,
		&result.OriginalUrl,
		&result.HashedUrl,
		&maxClicks,
		&result.ExpiresAt,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	result.MaxClicks = maxClicks.Int64

	return &result, nil
}

func (ur *urlRepo) Delete(id, userID int64) error {
	query := ` delete from urls where id=$1 and user_id=$2 `

	res, err := ur.db.Exec(
		query,
		id,
		userID,
	)
	if err != nil {
		return err
	}
	if count, _ := res.RowsAffected(); count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (ur *urlRepo) DecrementClick(url string) error {
	query := fmt.Sprintf(" UPDATE urls SET max_clicks = max_clicks - 1 WHERE hashed_url LIKE '%s'", "%"+url+"%")

	res, err := ur.db.Exec(query)
	if err != nil {
		return err
	}
	if count, _ := res.RowsAffected(); count == 0 {
		return sql.ErrNoRows
	}

	return nil
}
