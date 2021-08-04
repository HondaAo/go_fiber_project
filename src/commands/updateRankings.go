package main

import (
	"ambassador/src/database"
	"ambassador/src/model"
	"context"

	"github.com/go-redis/redis/v8"
)

func main() {
	database.Connect()
	database.SetupRedis()

	ctx := context.Background()

	var users []model.User

	database.DB.Find(&users, model.User{
		IsAmbassador: true,
	})

	for _, user := range users {
		ambassador := model.Ambassador(user)
		ambassador.CalculateRevenue(database.DB)

		database.Cache.ZAdd(ctx, "rankings", &redis.Z{
			Score:  *ambassador.Revenue,
			Member: user.Name(),
		})
	}
}
