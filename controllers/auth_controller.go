package controllers

import (
	"github.com/bruce-mig/go-admin/db"
	"github.com/bruce-mig/go-admin/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Id             uint   `gorm:"primary key;autoIncrement" json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `gorm:"unique" json:"email" `
	Password       string `json:"password"`
	VerifyPassword string `json:"verify_password"`
}

func Register(c *fiber.Ctx) error {
	var data registerRequest

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data.Password != data.VerifyPassword {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data.Password), 14)

	user := models.User{
		Id:        data.Id,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  password,
	}

	db.DB.Create(&user)

	return c.JSON(user)
}
