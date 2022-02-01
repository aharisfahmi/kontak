package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func main() {
	db, err := connectDB()
	if err != nil {
		log.Println("error connect db:", err)
	}
	db.AutoMigrate(&Contact{})

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/contacts", getContacts(db))
	e.GET("/contacts/:id", getContacts(db))
	e.POST("/contacts", createContact(db))
	e.PUT("/contacts/:id", updateContact(db))
	e.DELETE("/contacts/:id", deleteContact(db))
	e.Logger.Fatal(e.Start(":8080"))
}

type Contact struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (*Contact) TableName() string {
	return "kontak"
}

func getContacts(db *gorm.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
	  id := c.Param("id")
	  if id != "" {
	    var contact Contact
      db.First(&contact, "id = ?", id)
	    log.Println("contact id:", id)
	    log.Println("contact by id:", contact)
	    return c.JSON(http.StatusOK, contact)
	  }
		gender := c.QueryParam("gender")
		var contacts []Contact

		if gender != "" {
			db = db.Where("gender = ?", gender)
		}
		db.Find(&contacts)
		return c.JSON(http.StatusOK, contacts)
	}
}

func createContact(db *gorm.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		contact := new(Contact)
		if err := c.Bind(contact); err != nil {
			log.Println("error binding request:", err)
		}
		contact.Id = uuid.NewString()
		contact.CreatedAt = time.Now().String()
		contact.UpdatedAt = time.Now().String()
		res := db.Create(&contact)
		if res.Error != nil {
			return c.JSON(http.StatusInternalServerError, "something wrong")
		}
		log.Println("ct:", contact)
		
		db.First(&contact, contact.Id)

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
		var contact Contact
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
		var contact Contact
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
