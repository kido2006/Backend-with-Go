package store

import (
	"context"
	"database/sql"
	"time"
)

type Comment struct {
	ID        int64    `json:"id"`
	PostID    int64    `json:"postid"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags"`
	UserID    int64    `json:"userid"`
	CreatedAt int64    `json:"createdat"`
	User      User     `json:"user"`
}

type Commentstore struct {
	db *sql.DB
}

// Lấy comment theo postID
func (s *Commentstore) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	query := `
    SELECT 
        c.id, 
        c.post_id, 
        c.user_id, 
        c.content, 
        UNIX_TIMESTAMP(c.created_at), 
        u.username, 
        u.id
    FROM comments c
    JOIN users u ON u.id = c.user_id
    WHERE c.post_id = ?
    ORDER BY c.created_at DESC;
`

	ctx, cancel := context.WithTimeout(ctx, Querytimeout)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var c Comment
		c.User = User{}
		err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.UserID,
			&c.Content,
			&c.CreatedAt,
			&c.User.Username,
			&c.User.ID,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

// Tạo comment mới
func (s *Commentstore) Create(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO comments (post_id, user_id, content, created_at)
		VALUES (?, ?, ?, NOW())
	`

	ctx, cancel := context.WithTimeout(ctx, Querytimeout)
	defer cancel()

	res, err := s.db.ExecContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	)
	if err != nil {
		return err
	}

	// Lấy ID vừa insert
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	comment.ID = id

	// Gán created_at theo timestamp hiện tại
	comment.CreatedAt = time.Now().Unix()

	return nil
}
