package controllers

import (
	"ambassador/src/database"
	"ambassador/src/model"

	"github.com/gofiber/fiber/v2"
)

func Ambassador(c *fiber.Ctx) error {
	var users []model.User

	database.DB.Where("is_ambassador = true").Find(&users)

	return c.JSON(users)
}
