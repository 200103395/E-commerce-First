package controllers

import (
	"Ecommerce/database"
	"Ecommerce/models"
	"github.com/gofiber/fiber/v2"
)

func Shop(c *fiber.Ctx) error {
	// Needed variables
	var Items []models.Item
	var ItemsStock []models.ItemStock
	var cart []models.Cart
	var user models.User
	database.ReadSession(c, &user)
	// Getting all items from database
	database.DB.Find(&Items)
	models.Convert(Items, &ItemsStock)
	// Getting cart for given user
	price := 0
	database.DB.Table("carts").Where("user_id = ?", user.UserID).Scan(&cart)
	for _, i := range cart {
		price += int(i.Price * i.Quantity)
	}
	return c.Render("shop", fiber.Map{
		"Message":  "Success",
		"Items":    ItemsStock,
		"Username": user.Username,
		"Price":    price,
	})
}

func Filter(c *fiber.Ctx) error {
	// Needed variables
	var user models.User
	var data models.Filter
	var Items []models.Item
	var ItemsStock []models.ItemStock
	var cart []models.Cart
	database.ReadSession(c, &user)
	// Getting filtered items from database
	c.BodyParser(&data)
	database.DB.Table("items").Find(&Items)
	models.ConvertAndFilter(Items, &ItemsStock, data)
	// Getting cart for given user
	price := 0
	database.DB.Table("carts").Where("user_id = ?", user.UserID).Scan(&cart)
	for _, i := range cart {
		price += int(i.Price * i.Quantity)
	}
	return c.Render("filter", fiber.Map{
		"Username":   user.Username,
		"Items":      ItemsStock,
		"SearchData": data.Name,
		"MinRating":  data.MinRating,
		"MaxPrice":   data.MaxPrice,
		"MinPrice":   data.MinPrice,
		"Price":      price,
	})
}
