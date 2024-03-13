package controllers

import (
	"strconv"
	"time"

	"github.com/bruce-mig/go-admin/db"
	"github.com/bruce-mig/go-admin/models"
	"github.com/bruce-mig/go-admin/util"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gofiber/fiber/v2"
)

type (
	registerRequest struct {
		Id             uint   `gorm:"primary key;autoIncrement" json:"id"`
		FirstName      string `json:"first_name"`
		LastName       string `json:"last_name"`
		Email          string `gorm:"unique" json:"email" `
		Password       string `json:"password"`
		VerifyPassword string `json:"verify_password"`
	}

	loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	Claims struct {
		jwt.StandardClaims
	}
)

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

	user := models.User{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		RoleId:    1,
	}

	user.SetPassword(data.Password)
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

	err := user.ComparePassword(data.Password)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "incorrect creedentials",
		})
	}

	token, err := util.GenerateJwt(strconv.Itoa(int(user.Id)))
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

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	id, _ := util.ParseJwt(cookie)

	var user models.User
	db.DB.Where("id = ?", id).First(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	// c.ClearCookie("jwt")

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "successfully logged out",
	})
}
