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
    Title      string   `json:"title"`
    Content    string   `json:"content"`
    Categories []string `json:"categories"`
}

func (h *PostsHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
    user := r.Context().Value("user").(*models.User)

    var req createPostRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    post, err := h.DB.CreatePost(user.ID, req.Title, req.Content, req.Categories)
    if err != nil {
        http.Error(w, "could not create post", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
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

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(post)
}

func (h *PostsHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
    category := r.URL.Query().Get("category")

    posts, err := h.DB.ListPosts(category)
    if err != nil {
        http.Error(w, "could not list posts", http.StatusInternalServerError)
        return
    }

    if posts == nil {
        posts = []models.Post{}
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}

func (h *PostsHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
    cats, err := h.DB.ListCategories()
    if err != nil {
        http.Error(w, "could not list categories", http.StatusInternalServerError)
        return
    }

    if cats == nil {
        cats = []models.Category{}
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(cats)
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
