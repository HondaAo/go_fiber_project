package controllers

import (
	"ambassador/src/database"
	"ambassador/src/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Link(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var links []model.Link

	database.DB.Where("user_id = ?", id).Find(&links)

	for i, link := range links {
		var orders []model.Order

		database.DB.Where("code = ? and complete = true", link.Code).Find(&orders)

		links[i].Orders = orders
	}

	return c.JSON(links)
}
