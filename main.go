package main

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name   string
	Orders []Order
}

type Order struct {
	gorm.Model
	Product string
	Amount  float64
	UserID  uint
}

func main() {
	db, err := gorm.Open(sqlite.Open("go_transaction.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	db.AutoMigrate(&User{}, &Order{})

	err = db.Transaction(func(tx *gorm.DB) error {
		user := User{Name: "John"}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		orders := []Order{
			{Product: "Laptop", Amount: 20000.00, UserID: user.ID},
			{Product: "Mouse", Amount: 100.00, UserID: user.ID},
		}

		for _, order := range orders {
			if err := tx.Create(&order).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("Transaction failed: ", err)
	} else {
		fmt.Println("Transaction successful, all changes saved.")
	}

	var users []User
	db.Preload("Orders").Find(&users)
	for _, user := range users {
		fmt.Printf("User: %s\n", user.Name)
		for _, order := range user.Orders {
			fmt.Printf(" - Order: %s, Amount: %.2f\n", order.Product, order.Amount)
		}
	}
}
