package repo

import "time"

type Url struct {
	Id          int64
	UserId      int64
	OriginalUrl string
	HashedUrl   string
	MaxClicks   int64
	ExpiresAt   *time.Time
	CreatedAt   time.Time
}

type GetAllUrlsResult struct {
	Urls  []*Url
	Count int32
}

type GetAllUrlsParams struct {
	Limit  int32
	Page   int32
	Search string
}

type UrlStorageI interface {
	Create(u *Url) (*Url, error)
	Get(ulr string) (*Url, error)
	GetAll(params *GetAllUrlsParams) (*GetAllUrlsResult, error)
	Update(u *Url) (*Url, error)
	Delete(u *Url) (*Url, error)
}
