package database

import (
	"database/sql"
	"errors"
	"time"

	"real-time-forum/backend/models"
)

var ErrMessageNotFound = errors.New("message not found")

func (d *Database) SendMessage(senderID, receiverID int, content string) (*models.Message, error) {
    now := time.Now().UTC()
    ts := now.Format(time.RFC3339)

    res, err := d.DB.Exec(`
        INSERT INTO messages (sender_id, receiver_id, content, created_at)
        VALUES (?, ?, ?, ?)
    `, senderID, receiverID, content, ts)
    if err != nil {
        return nil, err
    }

    id, err := res.LastInsertId()
    if err != nil {
        return nil, err
    }

    return &models.Message{
        ID:         int(id),
        SenderID:   senderID,
        ReceiverID: receiverID,
        Content:    content,
        CreatedAt:  now,
        ReadAt:     nil,
    }, nil
}

func (d *Database) ListMessagesBetween(a, b int) ([]models.Message, error) {
    rows, err := d.DB.Query(`
        SELECT id, sender_id, receiver_id, content, created_at, read_at
        FROM messages
        WHERE (sender_id = ? AND receiver_id = ?)
           OR (sender_id = ? AND receiver_id = ?)
        ORDER BY created_at ASC
    `, a, b, b, a)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var msgs []models.Message

    for rows.Next() {
        var m models.Message
        var createdAtStr string
        var readAtStr sql.NullString

        if err := rows.Scan(&m.ID, &m.SenderID, &m.ReceiverID, &m.Content, &createdAtStr, &readAtStr); err != nil {
            return nil, err
        }

        m.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)

        if readAtStr.Valid {
            t, _ := time.Parse(time.RFC3339, readAtStr.String)
            m.ReadAt = &t
        }

        msgs = append(msgs, m)
    }

    return msgs, nil
}

func (d *Database) MarkMessageRead(id int) error {
    now := time.Now().UTC().Format(time.RFC3339)

    res, err := d.DB.Exec(`
        UPDATE messages
        SET read_at = ?
        WHERE id = ?
    `, now, id)
    if err != nil {
        return err
    }

    affected, _ := res.RowsAffected()
    if affected == 0 {
        return ErrMessageNotFound
    }

    return nil
}
    