package service

import (
	"log"
	"time"
	"utm-panel/core/system"
	"utm-panel/database"
)

type StatsService struct {
	trafficMon  *system.TrafficMonitor
	userManager *system.UserManager
}

func NewStatsService() *StatsService {
	return &StatsService{
		trafficMon:  &system.TrafficMonitor{},
		userManager: &system.UserManager{},
	}
}

// StartMonitoring : Lance la boucle infinie de surveillance
func (s *StatsService) StartMonitoring() {
	// Créer un "Ticker" qui tape toutes les 10 secondes
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for {
			// Attendre le prochain "tic"
			<-ticker.C
			s.SyncTraffic()
		}
	}()
}

// SyncTraffic : Vérifie tout le monde
func (s *StatsService) SyncTraffic() {
	db := database.GetDB()
	var clients []database.Client

	// 1. Récupérer tous les clients actifs
	if err := db.Where("enable = ?", true).Find(&clients).Error; err != nil {
		log.Println("Erreur lecture DB stats:", err)
		return
	}

	for _, client := range clients {
		// 2. Demander à IPTables la consommation actuelle (en octets)
		usage, err := s.trafficMon.GetUserUsage(client.Username)
		if err != nil {
			continue
		}

		// Si pas de changement, on passe au suivant
		// (Note : IPTables donne le total cumulé, donc on met à jour la valeur brute)
		if usage == 0 && client.Down == 0 {
			continue
		}

		// 3. Mettre à jour la DB
		// On considère ici que le trafic sortant du serveur = Download du client
		client.Down = usage
		// On pourrait ajouter l'Upload si on surveillait la chaine INPUT aussi
		
		db.Save(&client)

		// 4. VÉRIFICATION DES LIMITES (Le "Coupe-Circuit")
		shouldLock := false

		// A. Vérification Date d'Expiration
		if client.ExpiryTime > 0 && time.Now().Unix() > client.ExpiryTime {
			log.Printf("Client %s a expiré (Date). Verrouillage...", client.Username)
			shouldLock = true
		}

		// B. Vérification Quota (Si Total > 0)
		if client.Total > 0 && (client.Up+client.Down) >= client.Total {
			log.Printf("Client %s a dépassé son quota. Verrouillage...", client.Username)
			shouldLock = true
		}

		// Si une limite est atteinte -> BLOQUER
		if shouldLock {
			s.userManager.LockUser(client.Username) // Bloque Linux
			client.Enable = false                   // Bloque DB
			db.Save(&client)
		}
	}
}
