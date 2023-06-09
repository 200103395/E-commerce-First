package controllers

import (
	"Ecommerce/database"
	"Ecommerce/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"math"
	"strconv"
	"unicode"
)

func ItemComment(c *fiber.Ctx) error {
	// Needed variables
	var user models.User
	var Comments []models.Comment
	var Item models.Item
	var curRating models.Rating
	database.ReadSession(c, &user)
	// Getting item_id from path
	validId := true
	for _, char := range c.Path()[6:] {
		if unicode.IsDigit(char) != true {
			validId = false
		}
	}
	if validId == false || len(c.Path()[6:]) == 0 {
		return c.Redirect("/shop")
	}
	id := []byte(c.Path())[6:]
	// Getting item from database and converting for output
	database.DB.Where("item_id = ?", id).First(&Item)
	if Item.ItemID == 0 {
		return c.Redirect("/shop")
	}
	var ItemSt = models.ItemStock{
		ItemID:    Item.ItemID,
		Rating:    math.Floor(float64(Item.Rating)/(float64(Item.RatingNumber))*10) / 10,
		Review:    Item.Review,
		Name:      Item.Name,
		Price:     Item.Price,
		Publisher: Item.PublisherName,
	}
	if Item.RatingNumber == 0 {
		ItemSt.Rating = 0
	}
	// Getting comments and current user rating for given item
	database.DB.Table("comments").Where("item_id = ?", id).Scan(&Comments)
	database.DB.Table("ratings").Where("user_id = ? and item_id = ?", user.UserID, Item.ItemID).First(&curRating)
	// val is needed for showing user their previous rating for given item
	val := "Is" + strconv.Itoa(int(curRating.Rating))
	for i, comment := range Comments {
		if comment.UserID == user.UserID || user.Role == "Admin" {
			Comments[i].Deletable = true
		}
	}
	return c.Render("item", fiber.Map{
		"Item":     ItemSt,
		"Comments": Comments,
		"Username": user.Username,
		"Role":     user.Role,
		val:        true,
	})
}

func CommentConfirm(c *fiber.Ctx) error {
	// Needed variables
	var user models.User
	var data struct {
		Comment  string
		ItemID   uint
		Username string
	}
	c.BodyParser(&data)
	// Getting user from database
	res := database.DB.Where("username = ?", data.Username).First(&user)
	if res.Error == gorm.ErrRecordNotFound {
		return c.Redirect("/")
	}
	// Creating comment in database
	var comment = models.Comment{
		Comment:  data.Comment,
		UserID:   user.UserID,
		Username: user.Username,
		ItemID:   data.ItemID,
	}
	database.DB.Create(&comment)
	// Redirecting to the given item page
	id := int(data.ItemID)
	url := ("/item/") + strconv.Itoa(id)
	return c.Redirect(url)
}

func RatingConfirm(c *fiber.Ctx) error {
	// Needed variables
	var user models.User
	var rating models.Rating
	var prevRating uint
	var item models.Item
	var data struct {
		Username string
		ItemID   uint
		Rating   uint
	}
	c.BodyParser(&data)
	// Getting user and rating from database
	res := database.DB.Where("username = ?", data.Username).First(&user)
	if res.Error == gorm.ErrRecordNotFound {
		return c.Redirect("/shop")
	}
	res = database.DB.Where("user_id = ? and item_id = ?", user.UserID, data.ItemID).First(&rating)
	database.DB.Where("item_id = ?", data.ItemID).First(&item)
	// Changing item's rating
	if res.Error == gorm.ErrRecordNotFound {
		prevRating = 0
		item.RatingNumber++
	} else {
		prevRating = rating.Rating
	}
	item.Rating = item.Rating + data.Rating - prevRating
	database.DB.Save(&item)
	// Creating/changing rating in database
	rating = models.Rating{
		UserID: user.UserID,
		ItemID: data.ItemID,
		Rating: data.Rating,
	}
	if res.Error != gorm.ErrRecordNotFound {
		database.DB.Save(&rating)
	} else {
		database.DB.Create(&rating)
	}
	// Redirecting to the given item page
	url := ("/item/") + strconv.Itoa(int(data.ItemID))
	return c.Redirect(url)
}

func CommentDelete(c *fiber.Ctx) error {
	// Deleting comment from database using comment_id
	var data struct {
		CommentID int
	}
	c.BodyParser(&data)
	var comment models.Comment
	database.DB.Where("comment_id = ?", data.CommentID).First(&comment)
	database.DB.Where("comment_id = ?", data.CommentID).Delete(&models.Comment{})
	return c.Redirect("/item/" + strconv.Itoa(int(comment.ItemID)))
}
