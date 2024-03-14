package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id        uint   `gorm:"primary key;autoIncrement" json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `gorm:"unique" json:"email" `
	Password  []byte `json:"-"`
	RoleId    uint   `json:"role_id"`
	Role      Role   `json:"role" gorm:"foreignKey:RoleId"`
}

func (user *User) SetPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user.Password = hashedPassword
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}

func (user *User) Count(db *gorm.DB) int64 {
	var total int64
	db.Model(&User{}).Count(&total)
	return total
}

func (user *User) Take(db *gorm.DB, limit int, offset int) any {
	var users []User
	db.Preload("Role").Offset(offset).Limit(limit).Find(&users)

	return users
}
