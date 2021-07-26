// package main

// import (
// 	"ambassador/src/database"
// 	"ambassador/src/model"
// 	"math/rand"

// 	"github.com/bxcodec/faker/v3"
// )

// func main() {
// 	database.Connect()
// 	for i := 0; i < 30; i++ {
// 		product := model.Product{
// 			Title:       faker.Username(),
// 			Description: faker.Username(),
// 			Image:       faker.URL(),
// 			Price:       float64(rand.Intn(90) + 10),
// 		}

// 		database.DB.Create(&product)
// 	}
// }
