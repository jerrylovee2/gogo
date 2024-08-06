package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jerrylovee2/gogo/docs"
	handlers "github.com/jerrylovee2/gogo/handler"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				log.Println(e.Err)
			}
		}
	})

	r.POST("/books/create", handlers.CreateBookHandler)
	r.DELETE("/books/delete", handlers.DeleteBookHandler)
	r.GET("/books/all", handlers.GetAllBooksHandler)
	r.GET("/books/search", handlers.SearchBooksHandler)

	r.POST("/members/create", handlers.CreateMemberHandler)
	r.GET("/members/get", handlers.GetMemberByIDHandler)
	r.DELETE("/members/delete", handlers.DeleteMemberByIDHandler)

	r.POST("/borrowers/create", handlers.CreateBorrowerHandler)
	r.GET("/borrowers/get", handlers.GetBorrowerByIDHandler)
	r.DELETE("/borrowers/delete", handlers.DeleteBorrowerByIDHandler)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(r.Run(":" + port))
}
