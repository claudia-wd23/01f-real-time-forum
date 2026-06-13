package database

import (
	"database/sql"
	"errors"
	"time"

	"real-time-forum/backend/models"
)

var ErrCommentNotFound = errors.New("comment not found")

func (d *Database) CreateComment(userID, postID int, content string) (*models.Comment, error) {
    now := time.Now().UTC()
    ts := now.Format(time.RFC3339)

    res, err := d.DB.Exec(`
        INSERT INTO comments (post_id, user_id, content, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?)
    `, postID, userID, content, ts, ts)
    if err != nil {
        return nil, err
    }

    id, err := res.LastInsertId()
    if err != nil {
        return nil, err
    }

    return &models.Comment{
        ID:        int(id),
        PostID:    postID,
        UserID:    userID,
        Content:   content,
        CreatedAt: now,
        UpdatedAt: now,
    }, nil
}

func (d *Database) GetComment(id int) (*models.Comment, error) {
    row := d.DB.QueryRow(`
        SELECT id, post_id, user_id, content, created_at, updated_at
        FROM comments
        WHERE id = ?
    `, id)

    var c models.Comment
    var createdAtStr, updatedAtStr string

    if err := row.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &createdAtStr, &updatedAtStr); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrCommentNotFound
        }
        return nil, err
    }

    c.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
    c.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)

    return &c, nil
}

func (d *Database) ListComments(postID int) ([]models.Comment, error) {
    rows, err := d.DB.Query(`
        SELECT id, post_id, user_id, content, created_at, updated_at
        FROM comments
        WHERE post_id = ?
        ORDER BY created_at ASC
    `, postID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var comments []models.Comment

    for rows.Next() {
        var c models.Comment
        var createdAtStr, updatedAtStr string

        if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &createdAtStr, &updatedAtStr); err != nil {
            return nil, err
        }

        c.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
        c.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)

        comments = append(comments, c)
    }

    return comments, nil
}

func (d *Database) UpdateComment(id int, content string) error {
    now := time.Now().UTC().Format(time.RFC3339)

    res, err := d.DB.Exec(`
        UPDATE comments
        SET content = ?, updated_at = ?
        WHERE id = ?
    `, content, now, id)
    if err != nil {
        return err
    }

    affected, _ := res.RowsAffected()
    if affected == 0 {
        return ErrCommentNotFound
    }

    return nil
}

func (d *Database) DeleteComment(id int) error {
    res, err := d.DB.Exec(`DELETE FROM comments WHERE id = ?`, id)
    if err != nil {
        return err
    }

    affected, _ := res.RowsAffected()
    if affected == 0 {
        return ErrCommentNotFound
    }

    return nil
}
