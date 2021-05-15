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
	Name      string
	Languages []Language `gorm:"many2many:user_languages;"` // the relationship attribute
}

// Language represents languages model
type Language struct {
	gorm.Model
	Name string
}

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/relationships?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(User{}, Language{})

	//////////// 1- assign languages to user ////////////
	// first create the languages
	var languages = []Language{
		{Name: "english"},
		{Name: "chinese"},
	}
	db.Create(&languages)
	// next create a user
	var user = User{Name: "john"}
	db.Create(&user)
	// next assign the languages to the user
	db.Model(&user).Association("Languages").Append(&languages)

	//////////// 2- create a user, create languages, and assign them to the user at the same time ////////////
	var tom = User{
		Name: "tom",
		Languages: []Language{
			{Name: "arabic"},
			{Name: "russian"},
		},
	}
	db.Create(&tom)

	//////////// 3- skip the creation of the languages while creating the user ////////////
	var james = User{
		Name: "James",
		Languages: []Language{
			{Name: "arabic"},
			{Name: "russian"},
		},
	}
	db.Omit("Languages").Create(&james)
	// or if you want to skip all relationships
	db.Omit(clause.Associations).Create(&james)

	//////////// 4- find a user with his languages (eager loading) ////////////
	var dbUser User
	db.Preload("Languages").Where("name = ?", "john").First(&dbUser)
	fmt.Println(dbUser.Name)
	fmt.Println(dbUser.Languages)

	//////////// 5- update user's languages ////////////
	dbUser.Languages[0].Name = "updated language" // update the title of the first record in books
	// when updating a relationship you must use `db.Session`
	db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&dbUser)

	//////////// 6- finding a user's languages ////////////
	var userLanguages []Language
	db.Model(&dbUser).Where("name = ?", "chinese").Association("Languages").Find(&userLanguages)
	fmt.Println(userLanguages)

	//////////// 7- count a user's languages ////////////
	count := db.Model(&dbUser).Association("Languages").Count()
	fmt.Println(count)

	//////////// 8- unlink user's languages ////////////
	// first get the languages
	var langs []Language
	db.Model(Language{}).Where("name = ?", "chinese").Find(&langs)

	db.Model(&user).Association("Languages").Delete(langs)
}
