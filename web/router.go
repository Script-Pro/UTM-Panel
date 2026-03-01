package web

import (
	"embed"
	"net/http"
	"utm-panel/web/controller"

	"github.com/gin-gonic/gin"
)

// InitRouter : Configure les URLs du site
func InitRouter(assets embed.FS) *gin.Engine {
	router := gin.Default()

	// 1. Gestion des fichiers statiques (CSS, JS, Images)
	// On sert le dossier "assets" pour que le site s'affiche correctement
	router.Static("/assets", "./web/assets")

	// 2. Page d'accueil (Redirection vers le panel ou login)
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/panel")
	})

	// Page de connexion (Affichage HTML)
	router.GET("/login", func(c *gin.Context) {
		c.File("./web/html/login.html")
	})

	// === CORRECTION ICI ===
	// Route POST pour traiter la connexion (quand on clique sur "Login")
	router.POST("/login", controller.NewClientController().Login)

	// 3. Groupe "Panel" (Zone protégée)
	panelGroup := router.Group("/panel")
	{
		// Page principale (Tableau de bord)
		panelGroup.GET("/", func(c *gin.Context) {
			c.File("./web/html/index.html")
		})

		// Page des clients (Inbounds)
		panelGroup.GET("/inbounds", func(c *gin.Context) {
			c.File("./web/html/inbounds.html")
		})

		// Page des paramètres
		panelGroup.GET("/settings", func(c *gin.Context) {
			c.File("./web/html/settings.html")
		})

		// --- API (Les commandes invisibles exécutées par le Javascript) ---
		api := panelGroup.Group("/api/inbound")
		{
			// Créer une instance du contrôleur
			clientHandler := controller.NewClientController()

			api.POST("/add", clientHandler.AddClient)        // Ajouter
			api.POST("/update", clientHandler.UpdateClient)  // Modifier
			api.POST("/del/:id", clientHandler.DeleteClient) // Supprimer
			api.GET("/list", clientHandler.ListClients)      // Lister
			api.POST("/reset/:id", clientHandler.ResetClient)// Reset Quota
		}
	}

	return router
}
