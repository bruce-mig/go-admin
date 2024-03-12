package controllers

import (
	"strconv"
	"time"

	"github.com/bruce-mig/go-admin/db"
	"github.com/bruce-mig/go-admin/models"
	"github.com/dgrijalva/jwt-go/v4"
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

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
		// Id:        data.Id,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  password,
	}

	db.DB.Create(&user)

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data loginRequest

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	db.DB.Where("email = ?", data.Email).First(&user)

	if user.Id == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "not found",
		})
	}

	err := bcrypt.CompareHashAndPassword(user.Password, []byte(data.Password))
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "incorrect creedentials",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer: strconv.Itoa(int(user.Id)),
		ExpiresAt: &jwt.Time{
			time.Now().Add(time.Hour * 24),
		},
	})

	token, err := claims.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "successfully logged in",
	})
}
