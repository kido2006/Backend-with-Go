package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	Querytimeout = 5 * time.Second
	ErrConflict  = errors.New("record already exists")
)

type PostStores interface {
	Create(context.Context, *Post) error
	GetByID(context.Context, int64) (*Post, error)
	Delete(context.Context, int64) error
	Update(context.Context, *Post) error
	GetUserFeed(context.Context, int64, PaginationQuery) ([]PostWithData, error)
}

type UserStores interface {
	GetByID(context.Context, int64) (*User, error)
	GetByEmail(context.Context, string) (*User, error)
	Create(context.Context, *sql.Tx, *User) error
	CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error
	Activate(context.Context, string) error
	Delete(context.Context, int64) error
}

type Storage struct {
	Posts    PostStores
	Users    UserStores
	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostID(context.Context, int64) ([]Comment, error)
	}
	Followers interface {
		Follow(context.Context, int64, int64) error
		UnFollow(context.Context, int64, int64) error
	}
	Roles interface {
		GetByName(context.Context, string) (*Role, error)
	}
}

func NewSQL(db *sql.DB) Storage {
	return Storage{
		Posts:     &Poststore{db},
		Users:     &Userstore{db},
		Comments:  &Commentstore{db},
		Followers: &FollowerStore{db},
		Roles:     &RoleStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
