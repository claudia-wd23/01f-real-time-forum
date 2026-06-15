package database

import (
    "database/sql"
    "errors"
    "strings"
    "time"

    "real-time-forum/backend/models"
)

var ErrPostNotFound = errors.New("post not found")

func (d *Database) CreatePost(userID int, title, content string, categories []string) (*models.Post, error) {
    now := time.Now().UTC()
    ts := now.Format(time.RFC3339)

    res, err := d.DB.Exec(`
        INSERT INTO posts (user_id, title, content, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?)
    `, userID, title, content, ts, ts)
    if err != nil {
        return nil, err
    }

    postID, err := res.LastInsertId()
    if err != nil {
        return nil, err
    }

    for _, catName := range categories {
        var catID int
        err := d.DB.QueryRow(`SELECT id FROM categories WHERE name = ?`, catName).Scan(&catID)
        if err != nil {
            continue
        }
        d.DB.Exec(`INSERT OR IGNORE INTO post_categories (post_id, category_id) VALUES (?, ?)`, postID, catID)
    }

    return &models.Post{
        ID:         int(postID),
        UserID:     userID,
        Title:      title,
        Content:    content,
        Categories: categories,
        CreatedAt:  now,
        UpdatedAt:  now,
    }, nil
}

func (d *Database) GetPost(id int) (*models.Post, error) {
    row := d.DB.QueryRow(`
        SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at, p.updated_at,
               GROUP_CONCAT(c.name) AS categories
        FROM posts p
        LEFT JOIN users u ON p.user_id = u.id
        LEFT JOIN post_categories pc ON p.id = pc.post_id
        LEFT JOIN categories c ON pc.category_id = c.id
        WHERE p.id = ?
        GROUP BY p.id
    `, id)

    return scanPost(row)
}

func (d *Database) ListPosts(categoryFilter string) ([]models.Post, error) {
    query := `
        SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at, p.updated_at,
               GROUP_CONCAT(c.name) AS categories
        FROM posts p
        LEFT JOIN users u ON p.user_id = u.id
        LEFT JOIN post_categories pc ON p.id = pc.post_id
        LEFT JOIN categories c ON pc.category_id = c.id
    `
    args := []interface{}{}
    if categoryFilter != "" {
        query += ` WHERE p.id IN (
            SELECT pc2.post_id FROM post_categories pc2
            JOIN categories c2 ON pc2.category_id = c2.id
            WHERE c2.name = ?
        )`
        args = append(args, categoryFilter)
    }
    query += ` GROUP BY p.id ORDER BY p.created_at DESC`

    rows, err := d.DB.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []models.Post
    for rows.Next() {
        p, err := scanPost(rows)
        if err != nil {
            return nil, err
        }
        posts = append(posts, *p)
    }
    return posts, nil
}

func (d *Database) ListCategories() ([]models.Category, error) {
    rows, err := d.DB.Query(`SELECT id, name FROM categories ORDER BY name`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var cats []models.Category
    for rows.Next() {
        var c models.Category
        if err := rows.Scan(&c.ID, &c.Name); err != nil {
            return nil, err
        }
        cats = append(cats, c)
    }
    return cats, nil
}

func (d *Database) UpdatePost(id int, title, content string) error {
    now := time.Now().UTC().Format(time.RFC3339)

    res, err := d.DB.Exec(`
        UPDATE posts SET title = ?, content = ?, updated_at = ? WHERE id = ?
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

// scanner works for both *sql.Row and *sql.Rows
type scanner interface {
    Scan(dest ...interface{}) error
}

func scanPost(s scanner) (*models.Post, error) {
    var p models.Post
    var createdAtStr, updatedAtStr string
    var catRaw sql.NullString

    if err := s.Scan(&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content,
        &createdAtStr, &updatedAtStr, &catRaw); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrPostNotFound
        }
        return nil, err
    }

    p.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
    p.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)

    if catRaw.Valid && catRaw.String != "" {
        p.Categories = strings.Split(catRaw.String, ",")
    } else {
        p.Categories = []string{}
    }

    return &p, nil
}
