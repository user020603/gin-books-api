package models

import "time"

type BorrowedBook struct {
    ID         uint      `json:"id" gorm:"primaryKey"`
    BookID     uint      `json:"book_id" gorm:"not null"`
    UserID     uint      `json:"user_id" gorm:"not null"`
    BorrowedAt time.Time `json:"borrowed_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
    DueDate    time.Time `json:"due_date" gorm:"not null"`

    Book Book `json:"book" gorm:"foreignKey:BookID"`
    User User `json:"user" gorm:"foreignKey:UserID"`
}
