package main

import (
	routes "awesomeProject/api"
	"awesomeProject/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))

	packageContorller := controllers.NewPacksController()
	routes.CreateRoutes(r, packageContorller)

	err := r.Run(":8080")
	if err != nil {
		log.Println(err)
		return
	}
}
