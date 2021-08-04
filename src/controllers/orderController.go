package controllers

import (
	"ambassador/src/database"
	"ambassador/src/model"
	"context"
	"fmt"
	"net/smtp"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
)

func Orders(c *fiber.Ctx) error {
	var orders []model.Order

	database.DB.Preload("OrderItems").Find(&orders)

	for i, order := range orders {
		orders[i].Name = order.FullName()
	}

	return c.JSON(orders)
}

type OrderRequest struct {
	Code      string
	FirstName string
	LastName  string
	Email     string
	Address   string
	Country   string
	City      string
	Zip       string
	Products  []map[string]int
}

func CreateOrders(c *fiber.Ctx) error {
	var request OrderRequest

	if err := c.BodyParser(&request); err != nil {
		return err
	}

	link := model.Link{
		Code: request.Code,
	}

	database.DB.Preload("User").First(&link)

	if link.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid link!",
		})
	}

	order := model.Order{
		Code:            link.Code,
		UserId:          link.UserId,
		AmbassadorEmail: link.User.Email,
		Firstname:       request.FirstName,
		Lastname:        request.LastName,
		Email:           request.Email,
		Address:         request.Address,
		Country:         request.Country,
		City:            request.City,
		Zip:             request.Zip,
	}
	tx := database.DB.Begin()

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Error",
		})
	}

	var lineItem []*stripe.CheckoutSessionLineItemParams

	for _, requestProduct := range request.Products {
		product := model.Product{}
		product.Id = uint(requestProduct["product_id"])
		database.DB.First(&product)

		total := product.Price * float64(requestProduct["quantity"])

		item := model.OrderItem{
			OrderId:           order.Id,
			ProductTitle:      product.Title,
			Price:             product.Price,
			Quantity:          uint(requestProduct["quantity"]),
			AmbassadorRevenue: 0.1 * total,
			AdminRevenue:      0.9 * total,
		}
		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": "Error",
			})
		}

		lineItem = append(lineItem, &stripe.CheckoutSessionLineItemParams{
			Name:        stripe.String(product.Title),
			Description: stripe.String(product.Description),
			Images:      []*string{stripe.String(product.Image)},
			Amount:      stripe.Int64(100 * int64(product.Price)),
			Currency:    stripe.String("usd"),
			Quantity:    stripe.Int64(int64(requestProduct["quantity"])),
		})
	}
	stripe.Key = "sk_test_51HbGBEE3hSvEQEMsvqACtuuzfeFsiYmw6Q6ttXxUpd26evoT09HB8AphtyH6R3qjEeMqA4AtAmrDZdv8UpfMOFUG00tgsURdsV"

	params := stripe.CheckoutSessionParams{
		SuccessURL:         stripe.String("http://localhost:5000/success?source={CHECKOUT_SESSION_ID}"),
		CancelURL:          stripe.String("http://localhost:5000/error"),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItem,
	}

	source, err := session.New(&params)

	if err != nil {
		tx.Rollback()
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Error",
		})
	}

	order.TransctionId = source.ID

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Error",
		})
	}

	tx.Commit()
	return c.JSON(source)
}

func CompleteOrder(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	order := model.Order{}

	database.DB.Preload("OrderItem").First(&order, model.Order{
		TransctionId: data["source"],
	})

	if order.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "error",
		})
	}

	order.Complete = true
	database.DB.Save(&order)

	go func(order model.Order) {
		AmbassadorRevenue := 0
		AdminRevenue := 0

		for _, item := range order.OrderItems {
			AmbassadorRevenue += int(item.AmbassadorRevenue)
			AdminRevenue += int(item.AdminRevenue)
		}

		user := model.User{}
		user.Id = order.UserId

		database.DB.First(&user)
		database.Cache.ZIncrBy(context.Background(), "rankings", float64(AmbassadorRevenue), user.Name())

		ambassadorMessage := []byte(fmt.Sprintf("You earned $%s from the link #%s", AmbassadorRevenue, order.Code))

		smtp.SendMail("host.docker.internal:1025", nil, "no-replayemail.com", []string{order.AmbassadorEmail}, ambassadorMessage)

		adminMessage := []byte(fmt.Sprintf("You earned $%s from the link #%s", AdminRevenue, order.Code))

		smtp.SendMail("host.docker.internal:1025", nil, "no-replayemail.com", []string{"admin@admin.com"}, adminMessage)

	}(order)

	return c.JSON(fiber.Map{
		"message": "success",
	})

}
