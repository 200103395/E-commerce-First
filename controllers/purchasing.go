package controllers

import (
	"Ecommerce/database"
	"Ecommerce/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"math"
	"unicode"
)

func AddToCart(c *fiber.Ctx) error {
	var Item models.Item
	var user models.User
	var ItemSt models.ItemStock
	idValid := (c.Path())[5:]
	isValid := true
	for _, char := range idValid {
		isValid = unicode.IsDigit(char) && isValid
	}
	if isValid == false {
		return c.Redirect("/")
	}
	id := []byte(c.Path())[5:]
	database.DB.Where("item_id = ?", id).First(&Item)
	ItemSt = models.ItemStock{
		ItemID: Item.ItemID,
		Rating: math.Floor(float64(Item.Rating)/(float64(Item.RatingNumber))*10) / 10,
		Review: Item.Review,
		Name:   Item.Name,
		Price:  Item.Price,
	}
	if Item.RatingNumber == 0 {
		ItemSt.Rating = 0
	}
	database.ReadSession(c, &user)
	return c.Render("add", fiber.Map{
		"Username": user.Username,
		"Item":     ItemSt,
		"UserID":   user.UserID,
	})
}

func AddToCartConfirm(c *fiber.Ctx) error {
	var data models.Cart
	var oldData models.Cart
	c.BodyParser(&data)
	res := database.DB.Where("user_id = ? and item_id = ?", data.UserID, data.ItemID).First(&oldData)
	if res.Error == gorm.ErrRecordNotFound {
		database.DB.Create(&data)
	} else {
		oldData.Price = data.Price
		oldData.Quantity += data.Quantity
		database.DB.Save(&oldData)
	}
	return c.Redirect("/shop")
}

func EditCart(c *fiber.Ctx) error {
	var user models.User
	var cart models.Cart
	c.BodyParser(&cart)
	if res := database.ReadSession(c, &user); res == "Unauthorised" || user.UserID != cart.UserID {
		return c.Redirect("/")
	}
	return c.Render("cart_edit", fiber.Map{
		"Username": user.Username,
		"Cart":     cart,
	})
}

func EditCartConfirm(c *fiber.Ctx) error {
	var cart models.Cart
	var user models.User
	c.BodyParser(&cart)
	res := database.ReadSession(c, &user)
	if res == "Unauthorised" || cart.UserID != user.UserID {
		c.Redirect("/")
	}
	database.DB.Save(&cart)
	return c.Redirect("/purchase")
}

func DeleteFromCart(c *fiber.Ctx) error {
	var user models.User
	var cart models.Cart
	c.BodyParser(&cart)
	if res := database.ReadSession(c, &user); res == "Unauthorised" || user.UserID != cart.UserID {
		return c.Redirect("/")
	}
	database.DB.Delete(&cart)
	return c.Redirect("/purchase")
}

func Purchase(c *fiber.Ctx) error {
	var user models.User
	var cart []models.Cart
	var price uint
	res := database.ReadSession(c, &user)
	if res == "Unauthorised" {
		return c.Redirect("/")
	}
	database.DB.Table("carts").Where("user_id = ?", user.UserID).Scan(&cart)
	for _, cartI := range cart {
		price += cartI.Price * cartI.Quantity
	}
	return c.Render("purchase", fiber.Map{
		"Username": user.Username,
		"Cart":     cart,
		"Price":    price,
	})
}

func PurchaseConfirm(c *fiber.Ctx) error {
	var user models.User
	var cart []models.Cart
	res := database.ReadSession(c, &user)
	if res == "Unauthorised" {
		return c.Redirect("/")
	}
	var overallPrice uint
	database.DB.Table("carts").Where("user_id = ?", user.UserID).Scan(&cart)
	for _, i := range cart {
		overallPrice += i.Price * i.Quantity
	}
	if user.Tokens >= overallPrice {
		user.Tokens -= overallPrice
		database.DB.Save(&user)
		for _, i := range cart {
			database.DB.Delete(&i)
		}
	}
	return c.Redirect("/")
}
