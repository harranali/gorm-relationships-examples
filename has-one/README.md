# Has one

This example shows the relationship between `user` and his `books` and all the operations that you can perform.
```go
package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// User represents users model
type User struct {
	gorm.Model
	Name       string
	CreditCard CreditCard // the relationship attribute-1
}

// CreditCard represents credit card model
type CreditCard struct {
	gorm.Model
	Number uint
	UserID uint // the relationship attribute-2
}

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/relationships?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(User{}, CreditCard{})

	//////////// 1- create a user with his credit card ////////////
	var user = User{
		Name:       "john",
		CreditCard: CreditCard{Number: 4242424242424242},
	}
	db.Create(&user)

	//////////// 2- skip the creation of the credit card while creating the user ////////////
	db.Omit("CreditCard").Create(&user)
	// or if you want to skip all relationships
	db.Omit(clause.Associations).Create(&user)

	//////////// 3- find a user with his credit card (eager loading) ////////////
	var dbUser User
	db.Preload("CreditCard").Where("name = ?", "john").First(&dbUser)
	fmt.Println(dbUser.Name)
	fmt.Println(dbUser.CreditCard)

	//////////// 4- update user's credit card ////////////
	dbUser.CreditCard.Number = 555555555555444
	// when updating a relationship you must use `db.Session`
	db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&dbUser)

	//////////// 5- finding a user's credit card ////////////
	var creditCard CreditCard
	db.Model(&dbUser).Where("number = ?", 5555555555554444).Association("CreditCard").Find(&creditCard)
	fmt.Println(creditCard.ID)

	//////////// 6- delete a user's credit card when deleting the user ////////////
	// first let's create a user called `mike`
	var userMike = User{
		Name:       "mike",
		CreditCard: CreditCard{Number: 4111111111111111},
	}
	db.Create(&userMike)
	// delete user's credit card when deleting user
	db.Select("CreditCard").Delete(&userMike)
	// or if you want to delete all relationships
	db.Select(clause.Associations).Delete(&userMike)
}
```
