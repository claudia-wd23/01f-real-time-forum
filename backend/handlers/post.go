package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"real-time-forum/backend/database"
	"real-time-forum/backend/models"
)

type PostsHandler struct {
    DB *database.Database
}

type createPostRequest struct {
    Title   string `json:"title"`
    Content string `json:"content"`
}

func (h *PostsHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
    user := r.Context().Value("user").(*models.User)

    var req createPostRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    post, err := h.DB.CreatePost(user.ID, req.Title, req.Content)
    if err != nil {
        http.Error(w, "could not create post", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(post)
}

func (h *PostsHandler) GetPost(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, _ := strconv.Atoi(idStr)

    post, err := h.DB.GetPost(id)
    if err != nil {
        http.Error(w, "post not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(post)
}

func (h *PostsHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
    posts, err := h.DB.ListPosts()
    if err != nil {
        http.Error(w, "could not list posts", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(posts)
}

type updatePostRequest struct {
    Title   string `json:"title"`
    Content string `json:"content"`
}

func (h *PostsHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
    user := r.Context().Value("user").(*models.User)

    idStr := r.URL.Query().Get("id")
    id, _ := strconv.Atoi(idStr)

    post, err := h.DB.GetPost(id)
    if err != nil {
        http.Error(w, "post not found", http.StatusNotFound)
        return
    }

    if post.UserID != user.ID {
        http.Error(w, "forbidden", http.StatusForbidden)
        return
    }

    var req updatePostRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    if err := h.DB.UpdatePost(id, req.Title, req.Content); err != nil {
        http.Error(w, "could not update post", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *PostsHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
    user := r.Context().Value("user").(*models.User)

    idStr := r.URL.Query().Get("id")
    id, _ := strconv.Atoi(idStr)

    post, err := h.DB.GetPost(id)
    if err != nil {
        http.Error(w, "post not found", http.StatusNotFound)
        return
    }

    if post.UserID != user.ID {
        http.Error(w, "forbidden", http.StatusForbidden)
        return
    }

    if err := h.DB.DeletePost(id); err != nil {
        http.Error(w, "could not delete post", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}
