package controllers

import (
	"Ecommerce/database"
	"Ecommerce/models"
	"github.com/gofiber/fiber/v2"
)

func Publish(c *fiber.Ctx) error {
	// Getting user from session
	var user models.User
	res := database.ReadSession(c, &user)
	// Only authorised publisher or admin can publish
	if res == "Unauthorised" || user.Role == "Client" {
		return c.Redirect("/")
	}
	return c.Render("publish", fiber.Map{
		"Username": user.Username,
		"UserID":   user.UserID,
	})
}

func PublishConfirm(c *fiber.Ctx) error {
	// Needed variables
	var user models.User
	var data models.Item
	res := database.ReadSession(c, &user)
	if res == "Unauthorised" || user.Role == "Client" {
		return c.Redirect("/")
	}
	// Creating item in database
	c.BodyParser(&data)
	database.DB.Create(&data)
	return c.Redirect("/")
}
