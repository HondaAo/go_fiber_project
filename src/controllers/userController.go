package controllers

import (
	"ambassador/src/database"
	"ambassador/src/model"
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func Ambassador(c *fiber.Ctx) error {
	var users []model.User

	database.DB.Where("is_ambassador = true").Find(&users)

	return c.JSON(users)
}

func Ranking(c *fiber.Ctx) error {
	rankings, err := database.Cache.ZRevRangeByScoreWithScores(context.Background(), "rankings", &redis.ZRangeBy{
		Min: "inf",
		Max: "inf",
	}).Result()

	if err != nil {
		return err
	}

	result := make(map[string]float64)

	for _, ranking := range rankings {
		result[ranking.Member.(string)] = ranking.Score
	}
	return c.JSON(result)
}

func GetLinks(c *fiber.Ctx) error {
	code := c.Params("code")

	link := model.Link{
		Code: code,
	}

	database.DB.Preload("User").Preload("Products").First(&link)

	return c.JSON(link)
}
