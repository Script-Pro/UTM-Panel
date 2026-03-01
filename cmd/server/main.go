package main

import (
	"fmt"
	"log"
	"utm-panel/config"
	"utm-panel/database"
	"utm-panel/service"
	"utm-panel/web"
)

func main() {
	fmt.Println(">>> Démarrage de UTM Panel...")

	// 1. Charger la Configuration
	config.LoadConfig()

	// 2. Initialiser la Base de Données
	if err := database.InitDB(config.GlobalConfig.DBPath); err != nil {
		log.Fatalf("❌ Erreur DB: %v", err)
	}

	// 3. Démarrer la surveillance
	statsService := service.NewStatsService()
	statsService.StartMonitoring()

	// 4. Démarrer le Web
	// On appelle la fonction simplifiée
	router := web.InitRouter()
	
	port := ":" + config.GlobalConfig.ListenPort
	
	log.Printf("✅ Panel lancé sur http://VOTRE_IP%s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("❌ Erreur Web: %v", err)
	}
}
