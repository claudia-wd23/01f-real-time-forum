package database

import (
	"database/sql"
	"errors"
	"time"

	"real-time-forum/backend/models"
)

var ErrPostNotFound = errors.New("post not found")

func (d *Database) CreatePost(userID int, title, content string) (*models.Post, error) {
    now := time.Now().UTC()
    ts := now.Format(time.RFC3339)

    res, err := d.DB.Exec(`
        INSERT INTO posts (user_id, title, content, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?)
    `, userID, title, content, ts, ts)
    if err != nil {
        return nil, err
    }

    id, err := res.LastInsertId()
    if err != nil {
        return nil, err
    }

    return &models.Post{
        ID:        int(id),
        UserID:    userID,
        Title:     title,
        Content:   content,
        CreatedAt: now,
        UpdatedAt: now,
    }, nil
}

func (d *Database) GetPost(id int) (*models.Post, error) {
    row := d.DB.QueryRow(`
        SELECT id, user_id, title, content, created_at, updated_at
        FROM posts
        WHERE id = ?
    `, id)

    var p models.Post
    var createdAtStr, updatedAtStr string

    if err := row.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &createdAtStr, &updatedAtStr); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrPostNotFound
        }
        return nil, err
    }

    createdAt, _ := time.Parse(time.RFC3339, createdAtStr)
    updatedAt, _ := time.Parse(time.RFC3339, updatedAtStr)

    p.CreatedAt = createdAt
    p.UpdatedAt = updatedAt

    return &p, nil
}

func (d *Database) ListPosts() ([]models.Post, error) {
    rows, err := d.DB.Query(`
        SELECT id, user_id, title, content, created_at, updated_at
        FROM posts
        ORDER BY created_at DESC
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []models.Post

    for rows.Next() {
        var p models.Post
        var createdAtStr, updatedAtStr string

        if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &createdAtStr, &updatedAtStr); err != nil {
            return nil, err
        }

        p.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
        p.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)

        posts = append(posts, p)
    }

    return posts, nil
}

func (d *Database) UpdatePost(id int, title, content string) error {
    now := time.Now().UTC().Format(time.RFC3339)

    res, err := d.DB.Exec(`
        UPDATE posts
        SET title = ?, content = ?, updated_at = ?
        WHERE id = ?
    `, title, content, now, id)
    if err != nil {
        return err
    }

    affected, _ := res.RowsAffected()
    if affected == 0 {
        return ErrPostNotFound
    }

    return nil
}

func (d *Database) DeletePost(id int) error {
    res, err := d.DB.Exec(`DELETE FROM posts WHERE id = ?`, id)
    if err != nil {
        return err
    }

    affected, _ := res.RowsAffected()
    if affected == 0 {
        return ErrPostNotFound
    }

    return nil
}
