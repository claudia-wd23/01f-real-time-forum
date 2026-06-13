package models

import "time"

type User struct {
    ID           int
    Username     string
    Email        string
    Password     string      // keep if some parts of the code still use it
    PasswordHash string      // used by database/users.go
    CreatedAt    time.Time
}

type Post struct {
    ID        int
    UserID    int
    Title     string
    Content   string
    CreatedAt time.Time
    UpdatedAt time.Time
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
    SenderID   int        // expected by database/message.go
    ReceiverID int        // expected by database/message.go
    Content    string
    CreatedAt  time.Time
    ReadAt     *time.Time // can be NULL in DB
}

type Session struct {
    ID        string
    UserID    int
    Token     string
    CreatedAt time.Time
    ExpiresAt time.Time
}

/*type Session struct {
    ID        int
    UserID    int
    Token     string
    ExpiresAt time.Time
}*/

/*package models

import "time"

type User struct {
    ID        int
    Username  string
    Email     string
    Password  string
    CreatedAt time.Time
}

type Post struct {
    ID        int
    UserID    int
    Title     string
    Content   string
    CreatedAt time.Time
    UpdatedAt time.Time
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
    Content    string
    CreatedAt  time.Time
}*/

/*package models

type User struct {
    ID       int
    Username string
    Email    string
    Password string
}

type Post struct {
    ID      int
    UserID  int
    Title   string
    Content string
    CreatedAt string
}

type Comment struct {
    ID     int
    PostID int
    UserID int
    Content string
    CreatedAt string
}

type Message struct {
    ID int
    FromUserID int
    ToUserID int
    Content string
    CreatedAT string
}*/
