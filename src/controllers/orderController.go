package controllers

import (
	"ambassador/src/database"
	"ambassador/src/model"

	"github.com/gofiber/fiber/v2"
)

func Orders(c *fiber.Ctx) error {
	var orders []model.Order

	database.DB.Preload("OrderItems").Find(&orders)

	for i, order := range orders {
		orders[i].Name = order.FullName()
	}

	return c.JSON(orders)
}
