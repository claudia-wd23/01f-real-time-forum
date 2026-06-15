package models

import "time"

type User struct {
    ID           int
    Username     string
    Email        string
    Password     string
    PasswordHash string
    CreatedAt    time.Time
}

type Category struct {
    ID   int    `json:"ID"`
    Name string `json:"Name"`
}

type Post struct {
    ID         int      `json:"ID"`
    UserID     int      `json:"UserID"`
    Username   string   `json:"Username"`
    Title      string   `json:"Title"`
    Content    string   `json:"Content"`
    Categories []string `json:"Categories"`
    CreatedAt  time.Time `json:"CreatedAt"`
    UpdatedAt  time.Time `json:"UpdatedAt"`
}

type Comment struct {
    ID        int
    PostID    int
    UserID    int
    Content   string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Message struct {
    ID         int
    FromUserID int
    ToUserID   int
    SenderID   int
    ReceiverID int
    Content    string
    CreatedAt  time.Time
    ReadAt     *time.Time
}

type Session struct {
    ID        string
    UserID    int
    Token     string
    CreatedAt time.Time
    ExpiresAt time.Time
}
