package controllers

import (
	"Ecommerce/database"
	"Ecommerce/models"
	"github.com/gofiber/fiber/v2"
)

func WelcomePage(c *fiber.Ctx) error {
	// Getting user from token
	var user models.User
	res := database.ReadSession(c, &user)
	// Rendering main page
	return c.Render("main", fiber.Map{
		"Message":  res,
		"Username": user.Username,
	})
}
