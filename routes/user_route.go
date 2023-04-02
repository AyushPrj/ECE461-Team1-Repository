package routes

import (
	"ECE461-Team1-Repository/controllers" //add this
	"github.com/gin-gonic/gin"
)

func RepoRoute(router *gin.Engine) {
	router.POST("/repo", controllers.CreateRepo())
	router.GET("/repo/:repoId", controllers.GetARepo())
	router.PUT("/repo/:repoId", controllers.EditARepo())
	router.DELETE("/repo/:repoId", controllers.DeleteARepo())
	router.GET("/repos", controllers.GetAllRepos())
	router.GET("/raterepo", controllers.GetMetrics())
}
