package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func main() {
	db, err := connectDb()
	if err != nil {
		log.Println("error connect db:", err)
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/kontak", getContacts(db))
	e.Logger.Fatal(e.Start(":8080"))
}

type Kontak struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (*Kontak) TableName() string {
	return "kontak"
}

func getContacts(db *gorm.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		var contacts []Kontak
		db.Find(&contacts)
		return c.JSON(http.StatusOK, contacts)
	}
}
