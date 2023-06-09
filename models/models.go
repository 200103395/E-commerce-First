package models

import (
	"math"
	"strings"
)

// Items that we save in database
type Item struct {
	ItemID        uint   `json:"item_id" gorm:"primaryKey" gorm:"autoIncrement"`
	Name          string `json:"name"`
	Price         uint   `json:"price"`
	Review        string `json:"review"`
	Rating        uint   `json:"rating " gorm:"default:0"`
	RatingNumber  uint   `json:"rating-number" gorm:"default:0"`
	PublisherName string `json:"publisher_name"`
}

// Same items, but edited for output
type ItemStock struct {
	ItemID    uint    `json:"item_id"`
	Name      string  `json:"name"`
	Price     uint    `json:"price"`
	Review    string  `json:"review"`
	Rating    float64 `json:"rating"`
	Publisher string  `json:"publisher"`
}

// User model
type User struct {
	UserID   uint   `json:"user_id" gorm:"primaryKey" gorm:"autoIncrement"`
	Username string `json:"username" gorm:"unique"`
	Password []byte `json:"password"'`
	Role     string `json:"role" gorm:"default:'Client'"`
	Tokens   uint   `json:"tokens"`
}

// Rating model (rating from particular user to one item)
type Rating struct {
	UserID uint `json:"user_id"  gorm:"primaryKey"`
	ItemID uint `json:"item_id"  gorm:"primaryKey"`
	Rating uint `json:"rating"`
}

// Comment model (comment from particular user to one item)
type Comment struct {
	CommentID uint   `json:"comment_id" gorm:"primaryKey" gorm:"autoIncrement"`
	UserID    uint   `json:"user_id"`
	ItemID    uint   `json:"item_id"`
	Comment   string `json:"comment"`
	Username  string `json:"username"`
	Deletable bool   `gorm:"default:false"`
}

// Request from user to change role
type RoleRequest struct {
	UserID   uint   `json:"user_id" gorm:"primaryKey"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// Request from user to get tokens
type TokenRequest struct {
	UserID    uint   `json:"user_id"`
	RequestID uint   `json:"request_id" gorm:"primaryKey" gorm:"autoIncrement"`
	Username  string `json:"username"`
	Token     uint   `json:"token"`
}

// Item in cart (one item that is in particular user's cart)
type Cart struct {
	UserID   uint   `json:"user_id" gorm:"primaryKey"`
	ItemID   uint   `json:"item_id" gorm:"primaryKey"`
	ItemName string `json:"item_name"`
	Quantity uint   `json:"quantity"`
	Price    uint   `json:"price"`
}

// Filter for searching
type Filter struct {
	Name      string
	MinPrice  uint
	MaxPrice  uint
	MinRating float64
}

func Convert(items []Item, stock *[]ItemStock) {
	// Converting items getting from database to output form
	for _, item := range items {
		var val float64
		if item.RatingNumber == 0 {
			val = 0
		} else {
			val = math.Floor(float64(item.Rating)/float64(item.RatingNumber)*10) / 10
		}
		var itst = ItemStock{item.ItemID, item.Name, item.Price, item.Review, val, item.PublisherName}
		*stock = append(*stock, itst)
	}
}

func ConvertAndFilter(items []Item, stock *[]ItemStock, data Filter) {
	// Converting and filtering items getting from database to output form
	for _, item := range items {
		var val float64
		if item.RatingNumber == 0 {
			val = 0
		} else {
			val = math.Floor(float64(item.Rating)/float64(item.RatingNumber)*10) / 10
		}
		if val >= data.MinRating && item.Price >= data.MinPrice && item.Price <= data.MaxPrice && strings.Contains(strings.ToLower(item.Name), strings.ToLower(data.Name)) {
			var itst = ItemStock{item.ItemID, item.Name, item.Price, item.Review, val, item.PublisherName}
			*stock = append(*stock, itst)
		}
	}
}
