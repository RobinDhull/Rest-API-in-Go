package main

import (
	// "fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	// "golang.org/x/exp/slices"
	"errors"
	"strconv"
)

type Book struct {
	ID string `json:"id"`
	Title string `json:"name"`
	Author string `json:"author"`
	Quantity int `json:"quantity"`
}

var books = []Book {
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func GetBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func CreateBook(c *gin.Context) {
	var newBook Book
	var err error
	err = c.BindJSON(&newBook)
	if err != nil {
		return
	}
	for _, book := range books {
		if (book.ID == newBook.ID) {
			c.String(http.StatusOK, "Book is already in the list")
			return
		}
	}
	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func RemoveBook (c *gin.Context) {
	bookId := c.Query("id")
	for index, book := range books {
		if (book.ID == bookId) {
			book = books[index]
			books = append(books[:index], books[index+1:]...)
			c.String(http.StatusOK, book.Title + " is removed from the list")
			return
		}
	}
	c.String(http.StatusOK, "Book is not present in the list")
}

func CheckoutBook(c *gin.Context ) {
	id := c.Param("id")
	for index, book := range books {
		if book.ID == id {
			if (books[index].Quantity > 0) {
				books[index].Quantity--
				c.IndentedJSON(http.StatusOK, book)
			} else {
				c.IndentedJSON(http.StatusOK, gin.H{"message" : "Book is out of stock."})
			}
			return
		}
	}
}

func UpdateBook(c *gin.Context) {
	id := c.Query("id")
	quantity, err := strconv.Atoi(c.Query("quantity"))
	if (err != nil || quantity <= 0) {
		c.IndentedJSON(http.StatusOK, gin.H{"message" : "Invalid format"})
	} else {
		// for index, book := range books {
		// 	if (id == book.ID) {
		// 		books[index].Quantity += quantity
		// 		c.IndentedJSON(http.StatusOK, books[index])
		// 		return
		// 	}
		// }
		book, err := getBookByID(id)
		if (err == nil) {
			book.Quantity += quantity
			c.IndentedJSON(http.StatusOK, book)
			return
		}
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message":"Book not found"})
}

func getBookByID(id string) (*Book, error) {
	for index, book := range books {
		if book.ID == id {
			return &books[index], nil
		}
	}
	return nil, errors.New("Book not found")
}

func main() {
	router := gin.Default();
	router.GET("/books", GetBooks)
	router.POST("/books", CreateBook)
	router.DELETE("/books", RemoveBook)
	router.PATCH("/books/:id", CheckoutBook)
	router.PATCH("/books", UpdateBook)
	router.Run("localhost:3000")
}
