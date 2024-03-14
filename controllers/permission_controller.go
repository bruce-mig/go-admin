package controllers

import (
	"github.com/bruce-mig/go-admin/db"
	"github.com/bruce-mig/go-admin/models"
	"github.com/gofiber/fiber/v2"
)

func ListPermissions(c *fiber.Ctx) error {
	var permissions []models.Permission

	db.DB.Find(&permissions)

	return c.JSON(permissions)
}

func CreatePermission(c *fiber.Ctx) error {
	var permission models.Permission

	if err := c.BodyParser(&permission); err != nil {
		return err
	}

	db.DB.Create(&permission)

	return c.JSON(permission)
}
