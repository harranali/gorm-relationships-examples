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
	Name  string
	Books []Book
}

// Book represents books model
type Book struct {
	gorm.Model
	Title  string
	UserID uint
}

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/relationships?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(User{}, Book{})

	//////////// 1- create a user with his books ////////////
	var user = User{
		Name: "john",
		Books: []Book{
			{Title: "my first book"},
			{Title: "my second book"},
		},
	}
	db.Create(&user)

	//////////// 2- skip the creation of the books while creating the user ////////////
	db.Omit("Books").Create(&user)
	// or if you want to skip all relationships
	db.Omit(clause.Associations).Create(&user)

	//////////// 3- append to the user's books ////////////
	db.Model(&user).Association("Books").Append([]Book{
		{Title: "my third book"},
		{Title: "my fourth book"},
	})

	//////////// 4- find a user with his books (eager loading) ////////////
	var dbUser User
	db.Preload("Books").Where("name = ?", "john").First(&dbUser)
	fmt.Println(dbUser.Name)
	fmt.Println(dbUser.Books)

	//////////// 5- update user's books ////////////
	dbUser.Books[0].Title = "updated book title" // update the title of the first record in books
	// when updating a relationship you must use `db.Session`
	db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&dbUser)

	//////////// 6- finding a user's books ////////////
	var books []Book
	db.Model(&dbUser).Where("title = ?", "updated book title").Association("Books").Find(&books)
	// print the id of the first record of found records
	fmt.Println(books[0].ID)

	//////////// 7- count a user's books ////////////
	booksCount := db.Model(&dbUser).Association("Books").Count()
	fmt.Println(booksCount)

	//////////// 8- delete a user's books when deleting the user ////////////
	// first let's create a user called `mike`
	var userMike = User{
		Name: "john",
		Books: []Book{
			{Title: "book one"},
			{Title: "book two"},
		},
	}
	db.Create(&userMike)
	// delete user's books when deleting user
	db.Select("Books").Delete(&User{}, userMike.ID)
	// or if you want to delete all relationships
	db.Select(clause.Associations).Delete(&User{}, userMike.ID)
}
