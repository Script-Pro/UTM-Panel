package web

import (
	"net/http"
	"utm-panel/web/controller"

	"github.com/gin-gonic/gin"
)

// InitRouter : Version simplifiée qui lit directement sur le disque
func InitRouter() *gin.Engine {
	router := gin.Default()

	// 1. Gestion des fichiers statiques
	router.Static("/assets", "./web/assets")

	// 2. Page d'accueil
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/panel")
	})

	router.GET("/login", func(c *gin.Context) {
		c.File("./web/html/login.html")
	})

	// Login Action
	router.POST("/login", controller.NewClientController().Login)

	// 3. Groupe "Panel"
	panelGroup := router.Group("/panel")
	{
		panelGroup.GET("/", func(c *gin.Context) {
			c.File("./web/html/index.html")
		})

		panelGroup.GET("/inbounds", func(c *gin.Context) {
			c.File("./web/html/inbounds.html")
		})

		panelGroup.GET("/settings", func(c *gin.Context) {
			c.File("./web/html/settings.html")
		})

		// API
		api := panelGroup.Group("/api/inbound")
		{
			clientHandler := controller.NewClientController()
			api.POST("/add", clientHandler.AddClient)
			api.POST("/update", clientHandler.UpdateClient)
			api.POST("/del/:id", clientHandler.DeleteClient)
			api.GET("/list", clientHandler.ListClients)
			api.POST("/reset/:id", clientHandler.ResetClient)
		}
	}

	return router
}
