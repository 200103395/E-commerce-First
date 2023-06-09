package controllers

import (
	"Ecommerce/database"
	"Ecommerce/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"unicode"
)

func UserPage(c *fiber.Ctx) error {
	// Getting user and checking if admin
	var user models.User
	res := database.ReadSession(c, &user)
	if res == "Unauthorised" {
		return c.Redirect("/")
	}
	role := false
	if user.Role == "Admin" {
		role = true
	}
	// Rendering user page with several
	return c.Render("user_page", fiber.Map{
		"Username": user.Username,
		"UserRole": user.Role,
		"Tokens":   user.Tokens,
		"Role":     role,
	})
}

func ChangePassword(c *fiber.Ctx) error {
	// Getting user
	var user models.User
	res := database.ReadSession(c, &user)
	if res == "Unauthorised" {
		return c.Redirect("/")
	}
	return c.Render("change_password", fiber.Map{
		"Username": user.Username,
	})
}

func ChangePasswordConfirm(c *fiber.Ctx) error {
	// Getting data and user
	var data struct {
		Username     string
		OldPassword  string
		NewPassword  string
		NewPassword2 string
	}
	c.BodyParser(&data)
	var user models.User
	database.DB.Where("username = ?", data.Username).First(&user)
	// Comparing hash and given old password
	err := bcrypt.CompareHashAndPassword(user.Password, []byte(data.OldPassword))
	// Checking if new password is valid
	var validPchars, validPlower, validPupper, validPdigit bool
	validPchars = true
	for _, char := range data.NewPassword {
		if unicode.IsUpper(char) {
			validPupper = true
		}
		if unicode.IsLower(char) {
			validPlower = true
		}
		if unicode.IsDigit(char) {
			validPdigit = true
		}
		if !unicode.IsDigit(char) && !unicode.IsLetter(char) {
			validPchars = false
		}
	}
	// Checking in this order: incorrect old password, checking both new passwords, validation of new password
	if err != nil {
		return c.Render("change_password", fiber.Map{
			"Username": user.Username,
			"Error":    "Incorrect old password",
		})
	} else if data.NewPassword != data.NewPassword2 {
		return c.Render("change_password", fiber.Map{
			"Username": user.Username,
			"Error":    "New passwords not the same",
		})
	} else if validPlower == false || validPchars == false || validPupper == false || validPdigit == false || len(data.NewPassword) > 20 || len(data.NewPassword) < 8 {
		return c.Render("change_password", fiber.Map{
			"Username": user.Username,
			"Error":    "Invalid new password",
		})
	} else {
		password, _ := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
		user.Password = password
		database.DB.Save(&user)
		return c.Redirect("/user")
	}
}

func ChangeRole(c *fiber.Ctx) error {
	// Getting user from session
	var user models.User
	res := database.ReadSession(c, &user)
	if res == "Unauthorised" {
		return c.Redirect("/")
	}
	return c.Render("change_role", fiber.Map{
		"Username": user.Username,
		"Role":     user.Role,
	})
}

func ChangeRoleConfirm(c *fiber.Ctx) error {
	// Getting role request and user
	var data struct {
		Username string
		Role     string
	}
	c.BodyParser(&data)
	var user models.User
	res := database.DB.Where("username = ?", data.Username).First(&user)
	if res.Error == gorm.ErrRecordNotFound {
		return c.Redirect("/")
	}
	if data.Role != user.Role {
		// Creating role request in database
		var request models.RoleRequest
		res = database.DB.Where("user_id = ? and role = ?", user.UserID, data.Role).First(&request)
		if res.Error == gorm.ErrRecordNotFound {
			fmt.Println("Got your stuff")
			request = models.RoleRequest{
				Role:     data.Role,
				Username: data.Username,
				UserID:   user.UserID,
			}
			database.DB.Create(&request)
		}
	}
	return c.Redirect("/")
}

func RoleApprove(c *fiber.Ctx) error {
	// Getting user
	var user models.User
	res := database.ReadSession(c, &user)
	// If user is admin get roles requests
	if user.Role != "Admin" || res == "Unauthorised" {
		return c.Redirect("/")
	}
	var Requests []models.RoleRequest
	database.DB.Table("role_requests").Scan(&Requests)
	return c.Render("roles_approve", fiber.Map{
		"Username": user.Username,
		"Requests": Requests,
	})
}

func RoleApproveConfirm(c *fiber.Ctx) error {
	// Getting user and checking if Admin
	var user models.User
	res := database.ReadSession(c, &user)
	if user.Role != "Admin" || res == "Unauthorised" {
		return c.Redirect("/")
	}
	var data struct {
		UserID   uint
		Username string
		Role     string
	}
	c.BodyParser(&data)
	// Changing role for given user
	var changingUser models.User
	database.DB.Where("user_id = ?", data.UserID).First(&changingUser)
	changingUser.Role = data.Role
	database.DB.Save(&changingUser)
	// Deleting the request
	var request = models.RoleRequest{
		UserID:   data.UserID,
		Username: data.Username,
		Role:     data.Role,
	}
	database.DB.Delete(&request)
	return c.Redirect("/")
}

func RoleDecline(c *fiber.Ctx) error {
	// Getting user and checking if Admin
	var user models.User
	res := database.ReadSession(c, &user)
	if user.Role != "Admin" || res == "Unauthorised" {
		return c.Redirect("/")
	}
	var data struct {
		UserID   uint
		Username string
		Role     string
	}
	// Deleting request
	c.BodyParser(&data)
	var request = models.RoleRequest{
		UserID:   data.UserID,
		Username: data.Username,
		Role:     data.Role,
	}
	database.DB.Delete(&request)
	return c.Redirect("/")
}

func TokenRequest(c *fiber.Ctx) error {
	// Getting user from session
	var user models.User
	res := database.ReadSession(c, &user)
	if res == "Unauthorised" {
		return c.Redirect("/")
	}
	return c.Render("token_new", fiber.Map{
		"Username": user.Username,
		"UserID":   user.UserID,
	})
}

func TokenRequestConfirm(c *fiber.Ctx) error {
	var user models.User
	var request models.TokenRequest
	//var oldRequest models.TokenRequest
	res := database.ReadSession(c, &user)
	c.BodyParser(&request)
	if res == "Unauthorised" || user.UserID != request.UserID {
		return c.Redirect("/")
	}
	if request.Token == 0 {
		c.Redirect("/user")
	}
	database.DB.Create(&request)
	return c.Redirect("/user")
}

func TokenAdd(c *fiber.Ctx) error {
	var user models.User
	var requests []models.TokenRequest
	res := database.ReadSession(c, &user)
	if res == "Unauthorised" || user.Role != "Admin" {
		return c.Redirect("/")
	}
	database.DB.Table("token_requests").Scan(&requests)
	return c.Render("tokens_approve", fiber.Map{
		"Username": user.Username,
		"Requests": requests,
	})
}

func TokenAddConfirm(c *fiber.Ctx) error {
	var request models.TokenRequest
	var user models.User
	var changingUser models.User
	res := database.ReadSession(c, &user)
	if res == "Unauthorised" || user.Role != "Admin" {
		return c.Redirect("/")
	}
	c.BodyParser(&request)
	database.DB.Where("user_id = ?", request.UserID).First(&changingUser)
	changingUser.Tokens += request.Token
	database.DB.Delete(&request)
	database.DB.Save(&changingUser)
	return c.Redirect("/user")
}

func TokenAddDecline(c *fiber.Ctx) error {
	var request models.TokenRequest
	var user models.User
	res := database.ReadSession(c, &user)
	if res == "Unauthorised" || user.Role != "Admin" {
		return c.Redirect("/")
	}
	c.BodyParser(&request)
	database.DB.Delete(&request)
	return c.Redirect("/user")
}
