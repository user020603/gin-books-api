package models

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique;not null"`
	Email    string `json:"email" gorm:"unique; not null"`
	Password string `json:"password" gorm:"not null"`
	Active bool `json:"active" gorm:"default:true"`

	Reviews []Review `json:"-" gorm:"foreignKey:UserID"`
	BorrowedBooks []BorrowedBook `json:"-" gorm:"foreignKey:UserID"`
}
