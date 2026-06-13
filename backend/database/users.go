package database

import (
	"database/sql"
	"errors"
	"time"

	"real-time-forum/backend/models"
	"real-time-forum/backend/utils"
)

const sessionDuration = 7 * 24 * time.Hour

var ErrUserNotFound = errors.New("user not found")
var ErrInvalidCredentials = errors.New("invalid email or password")

func (d *Database) CreateUser(username, email, password string) (*models.User, error) {
    hashed, err := utils.HashPassword(password)
    if err != nil {
        return nil, err
    }

    now := time.Now().UTC()
    createdAtStr := now.Format(time.RFC3339)

    res, err := d.DB.Exec(`
        INSERT INTO users (username, email, password_hash, created_at)
        VALUES (?, ?, ?, ?)
    `, username, email, hashed, createdAtStr)
    if err != nil {
        return nil, err
    }

    id, err := res.LastInsertId()
    if err != nil {
        return nil, err
    }

    return &models.User{
        ID:           int(id),
        Username:     username,
        Email:        email,
        PasswordHash: hashed,
        CreatedAt:    now,
    }, nil
}

func (d *Database) GetUserByEmail(email string) (*models.User, error) {
    row := d.DB.QueryRow(`
        SELECT id, username, email, password_hash, created_at
        FROM users
        WHERE email = ?
    `, email)

    var u models.User
    var createdAtStr string

    if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &createdAtStr); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrUserNotFound
        }
        return nil, err
    }

    t, err := time.Parse(time.RFC3339, createdAtStr)
    if err != nil {
        return nil, err
    }
    u.CreatedAt = t

    return &u, nil
}

func (d *Database) GetUserByID(id int) (*models.User, error) {
    row := d.DB.QueryRow(`
        SELECT id, username, email, password_hash, created_at
        FROM users
        WHERE id = ?
    `, id)

    var u models.User
    var createdAtStr string

    if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &createdAtStr); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrUserNotFound
        }
        return nil, err
    }

    t, err := time.Parse(time.RFC3339, createdAtStr)
    if err != nil {
        return nil, err
    }
    u.CreatedAt = t

    return &u, nil
}

func (d *Database) GetUserByUsername(username string) (*models.User, error) {
    row := d.DB.QueryRow(`
        SELECT id, username, email, password_hash, created_at
        FROM users
        WHERE username = ?
    `, username)

    var u models.User
    var createdAtStr string

    if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &createdAtStr); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrUserNotFound
        }
        return nil, err
    }

    t, err := time.Parse(time.RFC3339, createdAtStr)
    if err != nil {
        return nil, err
    }
    u.CreatedAt = t

    return &u, nil
}

func (d *Database) AuthenticateUser(identifier, password string) (*models.User, error) {
    u, err := d.GetUserByEmail(identifier)
    if err != nil {
        if errors.Is(err, ErrUserNotFound) {
            u, err = d.GetUserByUsername(identifier)
            if err != nil {
                return nil, ErrInvalidCredentials
            }
        } else {
            return nil, err
        }
    }

    if !utils.CheckPassword(u.PasswordHash, password) {
        return nil, ErrInvalidCredentials
    }

    return u, nil
}

// Sessions
func (d *Database) CreateSession(userID int) (*models.Session, error) {
    now := time.Now().UTC()
    expiresAt := now.Add(sessionDuration)

    id := utils.NewSessionID()

    _, err := d.DB.Exec(`
        INSERT INTO sessions (id, user_id, created_at, expires_at)
        VALUES (?, ?, ?, ?)
    `, id, userID, now.Format(time.RFC3339), expiresAt.Format(time.RFC3339))
    if err != nil {
        return nil, err
    }

    return &models.Session{
        ID:        id,
        UserID:    userID,
        CreatedAt: now,
        ExpiresAt: expiresAt,
    }, nil
}

func (d *Database) GetSession(id string) (*models.Session, error) {
    row := d.DB.QueryRow(`
        SELECT id, user_id, created_at, expires_at
        FROM sessions
        WHERE id = ?
    `, id)

    var s models.Session
    var createdAtStr, expiresAtStr string

    if err := row.Scan(&s.ID, &s.UserID, &createdAtStr, &expiresAtStr); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrUserNotFound
        }
        return nil, err
    }

    createdAt, err := time.Parse(time.RFC3339, createdAtStr)
    if err != nil {
        return nil, err
    }
    expiresAt, err := time.Parse(time.RFC3339, expiresAtStr)
    if err != nil {
        return nil, err
    }

    if time.Now().UTC().After(expiresAt) {
        _ = d.DeleteSession(s.ID)
        return nil, errors.New("session expired")
    }

    s.CreatedAt = createdAt
    s.ExpiresAt = expiresAt

    return &s, nil
}

func (d *Database) DeleteSession(id string) error {
    _, err := d.DB.Exec(`DELETE FROM sessions WHERE id = ?`, id)
    return err
}

func (d *Database) CleanupExpiredSessions() error {
    nowStr := time.Now().UTC().Format(time.RFC3339)
    _, err := d.DB.Exec(`DELETE FROM sessions WHERE expires_at < ?`, nowStr)
    return err
}
