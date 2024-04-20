package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4/middleware"
)

type book struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

var (
	books     = make([]book, 0)
	runningID = 1
	lock      = sync.Mutex{}
)

func createBook(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	b := book{
		ID: runningID,
	}
	if err := c.Bind(b); err != nil {
		return err
	}
	books = append(books, b)
	runningID++
	return c.JSON(http.StatusCreated, b)
}

func getBook(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	id, _ := strconv.Atoi(c.Param("id"))
	var b book
	for _, book := range books {
		if book.ID == id {
			b = book
			break
		}
	}
	return c.JSON(http.StatusOK, b)
}

func updateBook(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	u := new(book)
	if err := c.Bind(u); err != nil {
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var b book
	for i, book := range books {
		if book.ID == id {
			book.Title = u.Title
			b = book
			books[i] = book
			break
		}

	}
	return c.JSON(http.StatusOK, b)
}

func deleteBook(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	id, _ := strconv.Atoi(c.Param("id"))
	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
			break
		}
	}
	return c.NoContent(http.StatusNoContent)
}

func getAllBooks(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	return c.JSON(http.StatusOK, books)
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/books", getAllBooks)
	e.POST("/books", createBook)
	e.GET("/books/:id", getBook)
	e.PUT("/books/:id", updateBook)
	e.DELETE("/books/:id", deleteBook)

	// Start server
	e.Logger.Fatal(e.Start(":1324"))
}
