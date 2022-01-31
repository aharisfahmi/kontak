package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func main() {
	db, err := connectDb()
	if err != nil {
		log.Println("error connect db:", err)
	}
	db.AutoMigrate(&Kontak{})

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/kontak", getContacts(db))
	e.POST("/kontak", createContact(db))
	e.PUT("/kontak/:id", updateContact(db))
	e.DELETE("/kontak/:id", deleteContact(db))
	e.Logger.Fatal(e.Start(":8080"))
}

type Kontak struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at" gorm:"->"`
	UpdatedAt string `json:"updated_at" gorm:"->"`
}

func (*Kontak) TableName() string {
	return "kontak"
}

func getContacts(db *gorm.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		gender := c.QueryParam("gender")
		var contacts []Kontak

		if gender != "" {
			db = db.Where("gender = ?", gender)
		}
		db.Find(&contacts)
		return c.JSON(http.StatusOK, contacts)
	}
}

func createContact(db *gorm.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		contact := new(Kontak)
		if err := c.Bind(contact); err != nil {
			log.Println("error binding request:", err)
		}
		contact.Id = uuid.NewString()
		log.Println("ct:", contact)
		res := db.Create(&contact)
		if res.Error != nil {
			return c.JSON(http.StatusInternalServerError, "something wrong")
		}

		// tx := db.Begin()
		// defer func() {
		// 	if r := recover(); r != nil {
		// 		tx.Rollback()
		// 	}
		// }()

		// if err := tx.Error; err != nil {
		// 	return err
		// }

		// res := tx.Create(&contact)
		// if res.Error != nil {
		// 	tx.Rollback()
		// 	return c.JSON(http.StatusInternalServerError, "something wrong")
		// }

		// err := tx.Commit().Error
		// if err != nil {
		// 	tx.Rollback()
		// 	return c.JSON(http.StatusInternalServerError, "something wrong")
		// }

		return c.JSON(http.StatusOK, contact)
	}
}

func updateContact(db *gorm.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		id := c.Param("id")
		var contact Kontak
		db.First(&contact, id)
		// contact := new(Kontak)
		// if err := c.Bind(contact); err != nil {
		// 	log.Println("error binding request:", err)
		// }
		log.Println("ct:", contact)
		// db.Find(&contacts)
		// res := db.Create(&contact)
		// if res.Error != nil {
		// 	return c.JSON(http.StatusInternalServerError, "something wrong")
		// }
		return c.JSON(http.StatusOK, contact)
	}
}

func deleteContact(db *gorm.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		id := c.Param("id")
		var contact Kontak
		err := db.First(&contact, id).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "contact not found")
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "something wrong")
		}
		log.Println("ct:", contact)
		err = db.Delete(&contact).Error
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "delete failed")
		}
		return c.JSON(http.StatusOK, contact)
	}
}
