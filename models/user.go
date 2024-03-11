package models

type User struct {
	Id        uint   `gorm:"primary key;autoIncrement" json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `gorm:"unique" json:"email" `
	Password  []byte `json:"password"`
}
