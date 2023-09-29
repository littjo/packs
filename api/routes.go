package routes

import (
	"awesomeProject/controllers"
	"github.com/gin-gonic/gin"
)

func CreateRoutes(r *gin.Engine, pc *controllers.PacksController) {

	r.POST("/packs", pc.WritePacksHandler)
	r.GET("/packs", pc.ReadPacksHandler)
	r.GET("/order/:items", pc.CalculatePacksHandler)
	r.Static("/web", "./public")
}
