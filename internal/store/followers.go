package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

type Follower struct {
	UserID     int64 `json:"user_id"`
	FollowerID int64 `json:"follower_id"`
	CreatedAt  int64 `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, followerID, UserID int64) error {
	query := `
		INSERT INTO followers (user_id, follower_id, created_at)
		VALUE (?, ?, NOW())`

	ctx, cancel := context.WithTimeout(ctx, Querytimeout)
	defer cancel()
	_, err := s.db.ExecContext(ctx, query, UserID, followerID)

	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return fmt.Errorf("you are already following this user")
		}
	}

	return err
}

func (s *FollowerStore) UnFollow(ctx context.Context, followerID, UserID int64) error {
	query := `
		DELETE FROM followers 
		WHERE user_id = ? AND follower_id = ?`

	ctx, cancel := context.WithTimeout(ctx, Querytimeout)
	defer cancel()
	_, err := s.db.ExecContext(ctx, query, UserID, followerID)
	return err
}
