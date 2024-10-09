package models

type Author struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Bio   string `json:"bio"`
	Email string `json:"email"`

	Book []Book `json:"book" gorm:"foreignKey:AuthorID"`
}
