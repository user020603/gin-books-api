package models

type Review struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	BookID  uint   `json:"book_id" gorm:"not null"`
	UserID  uint   `json:"user_id" gorm:"not null"` 
	Rating  int    `json:"rating" gorm:"check:rating >= 1 AND rating <= 5; not null"`
	Comment string `json:"comment"`

	Book Book `json:"book" gorm:"foreignKey:BookID"` 
	User User `json:"user" gorm:"foreignKey:UserID"`
}