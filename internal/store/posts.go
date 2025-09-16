package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"userid"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostWithData struct {
	Post
	CommentCount int `json:"commentcount"`
}

type Poststore struct {
	db *sql.DB
}

func (s *Poststore) GetUserFeed(ctx context.Context, userID int64, fq PaginationQuery) ([]PostWithData, error) {
	conditions := []string{}
	args := []any{userID, userID, fq.Search, fq.Search}

	if len(fq.Tags) > 0 {
		for _, tag := range fq.Tags {
			conditions = append(conditions, "JSON_CONTAINS(p.tags, JSON_QUOTE(?))")
			args = append(args, tag)
		}
	}

	query := `
	SELECT DISTINCT
		p.id,
		p.user_id,
		p.title,
		p.content,
		p.created_at,
		p.version,
		p.tags,
		u.username
	FROM posts p
	LEFT JOIN users u ON p.user_id = u.id
	JOIN followers f ON f.user_id = p.user_id
	WHERE 
		(f.follower_id = ? OR p.user_id = ?)
		AND (LOWER(p.title) LIKE CONCAT('%', LOWER(?), '%')
			 OR LOWER(p.content) LIKE CONCAT('%', LOWER(?), '%'))`

	if len(conditions) > 0 {
		query += " AND (" + strings.Join(conditions, " OR ") + ")"
	}

	query += `
	ORDER BY p.created_at DESC
	LIMIT ? OFFSET ?;`

	args = append(args, fq.Limit, fq.Offset)

	ctx, cancel := context.WithTimeout(ctx, Querytimeout)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []PostWithData
	for rows.Next() {
		var post PostWithData
		var tagsSQL sql.NullString

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.Version,
			&tagsSQL,
			&post.User.Username,
		)
		if err != nil {
			return nil, err
		}

		if tagsSQL.Valid && tagsSQL.String != "" {
			if err := json.Unmarshal([]byte(tagsSQL.String), &post.Tags); err != nil {
				return nil, err
			}
		} else {
			post.Tags = []string{}
		}

		post.Comments = []Comment{}
		feeds = append(feeds, post)
	}

	return feeds, nil
}

func (s *Poststore) Create(ctx context.Context, post *Post) error {
	// convert slice -> JSON để lưu vào MySQL JSON column
	tagsJSON, err := json.Marshal(post.Tags)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO posts (content, title, user_id, tags, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())`

	ctx, cancel := context.WithTimeout(ctx, Querytimeout)
	defer cancel()

	result, err := s.db.ExecContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		tagsJSON,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	post.ID = id

	// lấy created_at, updated_at từ DB để đồng bộ
	err = s.db.QueryRowContext(
		ctx,
		`SELECT created_at, updated_at FROM posts WHERE id = ?`,
		post.ID,
	).Scan(&post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}

	// comments sẽ để trống, load ở hàm khác
	post.Comments = []Comment{}

	return nil
}

func (s *Poststore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `SELECT id, title, content, user_id, tags, created_at, updated_at, version
              FROM posts
              WHERE id = ?`
	ctx, cancel := context.WithTimeout(ctx, Querytimeout)
	defer cancel()

	var p Post
	var tagsSQL sql.NullString

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.Title,
		&p.Content,
		&p.UserID,
		&tagsSQL,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.Version,
	)
	if err != nil {
		log.Printf("DB scan error: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	// parse tags JSON thành []string
	if tagsSQL.Valid && tagsSQL.String != "" {
		if err := json.Unmarshal([]byte(tagsSQL.String), &p.Tags); err != nil {
			return nil, err
		}
	} else {
		p.Tags = []string{}
	}

	// comments để trống (load ở chỗ khác)
	p.Comments = []Comment{}
	return &p, nil
}

func (s *Poststore) Delete(ctx context.Context, postID int64) error {
	query := `DELETE FROM posts WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, Querytimeout)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil

}

func (s *Poststore) Update(ctx context.Context, post *Post) error {
	tagsJSON, err := json.Marshal(post.Tags)
	if err != nil {
		return err
	}

	// Bước 1: Update
	updateQuery := `UPDATE posts
	SET version = version + 1,
		title = ?, content = ?, tags = ?, updated_at = NOW()
	WHERE id = ? AND version = ?`

	res, err := s.db.ExecContext(ctx, updateQuery, post.Title, post.Content, tagsJSON, post.ID, post.Version)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no rows updated")
	}

	// Bước 2: Lấy version mới
	query := `SELECT version FROM posts WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, Querytimeout)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, post.ID)
	if err := row.Scan(&post.Version); err != nil {
		return err
	}

	return nil
}
