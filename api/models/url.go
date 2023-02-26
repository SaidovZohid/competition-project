package models

import "time"

type Url struct {
	Id          int64      `json:"id"`
	UserId      int64      `json:"user_id"`
	OriginalUrl string     `json:"original_url"`
	HashedUrl   string     `json:"hashed_url"`
	MaxClicks   *int64     `json:"max_clicks"`
	ExpiresAt   *time.Time `json:"expires_at"`
	CreatedAt   string     `json:"created_at"`
}
type CreateShortUrlRequest struct {
	OriginalUrl string `json:"original_url" binding:"required"`
	MaxClicks   int64  `json:"max_clicks"`
	Duration    string `json:"duration"`
}

type CreateUrlRequest struct {
	UserId      int64      `json:"user_id"`
	OriginalUrl string     `json:"original_url"`
	HashedUrl   string     `json:"hashed_url"`
	MaxClicks   *int64     `json:"max_clicks"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

type UpdateUrlRequest struct {
	HashedUrl string     `json:"hashed_url"`
	MaxClicks *int64     `json:"max_clicks"`
	ExpiresAt *time.Time `json:"expires_at"`
}
