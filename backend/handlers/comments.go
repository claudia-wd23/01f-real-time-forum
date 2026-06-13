package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"real-time-forum/backend/database"
	"real-time-forum/backend/models"
)

type CommentsHandler struct {
    DB *database.Database
}

type createCommentRequest struct {
    PostID  int    `json:"post_id"`
    Content string `json:"content"`
}

func (h *CommentsHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
    user := r.Context().Value("user").(*models.User)

    var req createCommentRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    comment, err := h.DB.CreateComment(user.ID, req.PostID, req.Content)
    if err != nil {
        http.Error(w, "could not create comment", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(comment)
}

func (h *CommentsHandler) ListComments(w http.ResponseWriter, r *http.Request) {
    postIDStr := r.URL.Query().Get("post_id")
    postID, _ := strconv.Atoi(postIDStr)

    comments, err := h.DB.ListComments(postID)
    if err != nil {
        http.Error(w, "could not list comments", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(comments)
}

type updateCommentRequest struct {
    Content string `json:"content"`
}

func (h *CommentsHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
    user := r.Context().Value("user").(*models.User)

    idStr := r.URL.Query().Get("id")
    id, _ := strconv.Atoi(idStr)

    comment, err := h.DB.GetComment(id)
    if err != nil {
        http.Error(w, "comment not found", http.StatusNotFound)
        return
    }

    if comment.UserID != user.ID {
        http.Error(w, "forbidden", http.StatusForbidden)
        return
    }

    var req updateCommentRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    if err := h.DB.UpdateComment(id, req.Content); err != nil {
        http.Error(w, "could not update comment", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *CommentsHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
    user := r.Context().Value("user").(*models.User)

    idStr := r.URL.Query().Get("id")
    id, _ := strconv.Atoi(idStr)

    comment, err := h.DB.GetComment(id)
    if err != nil {
        http.Error(w, "comment not found", http.StatusNotFound)
        return
    }

    if comment.UserID != user.ID {
        http.Error(w, "forbidden", http.StatusForbidden)
        return
    }

    if err := h.DB.DeleteComment(id); err != nil {
        http.Error(w, "could not delete comment", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

/*package database

import (
	"database/sql"
	"errors"
	"time"

	"backend/models"
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
}*/