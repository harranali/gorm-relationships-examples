# Belongs To

This example shows the relationship between `user` and his `company` and some of the operations that you can perform.
```go
package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Company represents company model
type Company struct {
	gorm.Model
	Name string
}

// User represents users model
type User struct {
	gorm.Model
	Name      string
	// the relationship attribs
	CompanyID uint
	Company   Company
}

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/relationships?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(User{}, Company{})

	//////////// 1- create a user of the company ////////////
	// first create a company
	var company = Company{Name: "AZ Company"}
	db.Create(&company)
	// next ceate a user of the company
	var user = User{Name: "john", CompanyID: company.ID}
	db.Create(&user)

	//////////// 4- find a user and load his company (eager loading) ////////////
	var dbUser User
	db.Preload("Company").Where("name = ?", "john").First(&dbUser)
	fmt.Println(dbUser.Name)
	fmt.Println(dbUser.Company)

	//////////// 5- update user's Company ////////////
	dbUser.Company.Name = "A TO Z Company" // update the name of the company
	// when updating a relationship you must use `db.Session`
	db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&dbUser)

}
```
