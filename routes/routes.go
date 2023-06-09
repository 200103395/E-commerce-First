package routes

import (
	"Ecommerce/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Get("/", controllers.WelcomePage)
	//Authorisation
	app.Get("/register", controllers.Register)
	app.Post("/register/confirm", controllers.RegisterConfirm)
	app.Get("/login", controllers.Login)
	app.Post("/login/confirm", controllers.LoginConfirm)
	app.Get("/logout", controllers.Logout)
	//Shop
	app.Get("/shop", controllers.Shop)
	app.Post("/filter", controllers.Filter)
	//Item_info
	app.Use("/item/", controllers.ItemComment)
	app.Post("/comment/confirm", controllers.CommentConfirm)
	app.Post("/comment/delete", controllers.CommentDelete)
	app.Post("/rating/confirm", controllers.RatingConfirm)
	//User_info
	app.Get("/user", controllers.UserPage)
	app.Get("/change_password", controllers.ChangePassword)
	app.Post("/change_password/confirm", controllers.ChangePasswordConfirm)
	app.Get("/change_role", controllers.ChangeRole)
	app.Post("/change_role/confirm", controllers.ChangeRoleConfirm)
	app.Get("/role_approve", controllers.RoleApprove)
	app.Post("/role_approve/confirm", controllers.RoleApproveConfirm)
	app.Post("/role/decline", controllers.RoleDecline)
	app.Get("/token_new", controllers.TokenRequest)
	app.Post("/token_new/confirm", controllers.TokenRequestConfirm)
	app.Get("/token_approve", controllers.TokenAdd)
	app.Post("/token_approve/confirm", controllers.TokenAddConfirm)
	app.Post("/token_approve/decline", controllers.TokenAddDecline)
	//Publish
	app.Get("/publish", controllers.Publish)
	app.Post("/publish/confirm", controllers.PublishConfirm)
	//Purchase
	app.Get("/purchase", controllers.Purchase)
	app.Use("/add/", controllers.AddToCart)
	app.Post("/cart/add", controllers.AddToCartConfirm)
	app.Post("/cart/edit", controllers.EditCart)
	app.Post("/cart/edit/confirm", controllers.EditCartConfirm)
	app.Post("/cart/delete", controllers.DeleteFromCart)
	app.Get("/purchase/confirm", controllers.PurchaseConfirm)
}
