package main

import (
	"github.com/goombaio/namegenerator"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

type Attendee struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `json:"name" gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type apiHandler struct {
	db            gorm.DB
	nameGenerator namegenerator.Generator
}

func newAPIHandler(db gorm.DB) *apiHandler {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	return &apiHandler{db: db, nameGenerator: nameGenerator}
}

func (h *apiHandler) addAttendee(c echo.Context) error {
	attendee := Attendee{Name: h.nameGenerator.Generate()}
	h.db.Create(&attendee)
	return c.JSON(http.StatusCreated, attendee)
}

func (h *apiHandler) home(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func (h *apiHandler) listAttendees(c echo.Context) error {
	var attendees []Attendee
	h.db.Find(&attendees)
	return c.JSON(http.StatusOK, attendees)
}

func (h *apiHandler) envVarList(c echo.Context) error {
	envVars := os.Environ()
	return c.JSON(http.StatusOK, envVars)
}

func main() {

	e := echo.New()

	dsn := os.Getenv("POSTGRES_URI")
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// Migrate the schema
	db.AutoMigrate(&Attendee{})

	h := newAPIHandler(*db)

	e.GET("/", h.home)
	e.POST("/attendees", h.addAttendee)
	e.GET("/attendees", h.listAttendees)
	e.GET("/env", h.envVarList)

	e.GET("/api/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
