package main

import (
	"Ecommerce/database"
	"Ecommerce/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func main() {
	//Setting up database, templates, statics, routes
	database.Setup()
	engine := html.New("./templates/", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/static", "./static")
	routes.Setup(app)
	//Serving http request from given port
	app.Listen(":8000")
}
