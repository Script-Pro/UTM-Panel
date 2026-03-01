package main

import (
	"embed"
	"fmt"
	"log"
	"utm-panel/config"
	"utm-panel/database"
	"utm-panel/service"
	"utm-panel/web"
)

// Cette ligne magique permet d'inclure les fichiers HTML/CSS dans l'exécutable final
// (Même si on utilise des CDN pour l'instant, c'est une bonne pratique Go)
//
//go:embed all:web/assets all:web/html
var assets embed.FS

func main() {
	fmt.Println(`
	 _   _ _____ __  __      _____   _    _   _ 
	| | | |_   _|  \/  |    |  __ \ (_)  | | (_)
	| | | | | | | .  . |    | |__) | _ __| |_ _ 
	| | | | | | | |\/| |    |  ___/ | | _| __| |
	| |_| | | | | |  | |    | |     | | | | |_| |
	 \___/  |_| |_|  |_|    |_|     |_|_|  \__|_|
	                                             
	>>> Démarrage du contrôleur Universel...
	`)

	// 1. Charger la Configuration
	config.LoadConfig()
	log.Println("Configuration chargée.")

	// 2. Initialiser la Base de Données
	// On utilise le chemin défini dans la config
	if err := database.InitDB(config.GlobalConfig.DBPath); err != nil {
		log.Fatalf("❌ Erreur fatale Base de Données: %v", err)
	}

	// 3. Démarrer le Service de Statistiques (Surveillance Quota/Expiration)
	// Il tourne en arrière-plan (goroutine)
	statsService := service.NewStatsService()
	statsService.StartMonitoring()
	log.Println("✅ Surveillance du trafic (IPTables) active.")

	// 4. Démarrer le Serveur Web (Interface Graphique)
	router := web.InitRouter(assets)
	
	port := ":" + config.GlobalConfig.ListenPort
	log.Printf("✅ Panel accessible sur http://VOTRE_IP%s", port)
	log.Printf("👤 Login par défaut: %s / %s", config.GlobalConfig.AdminUser, config.GlobalConfig.AdminPass)

	// Lancement bloquant (reste allumé tant qu'on ne coupe pas)
	if err := router.Run(port); err != nil {
		log.Fatalf("❌ Erreur démarrage Web: %v", err)
	}
}
