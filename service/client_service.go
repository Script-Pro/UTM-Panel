package service

import (
	"errors"
	"utm-panel/core/system"
	"utm-panel/database"
)

type ClientService struct {
	userManager *system.UserManager
	trafficMon  *system.TrafficMonitor
}

func NewClientService() *ClientService {
	return &ClientService{
		userManager: &system.UserManager{},
		trafficMon:  &system.TrafficMonitor{},
	}
}

// AddClient : Ajoute un client dans la DB + Créé l'user Linux
func (s *ClientService) AddClient(client *database.Client) error {
	// 1. Vérifier si le nom d'utilisateur est vide
	if client.Username == "" || client.Password == "" {
		return errors.New("username et password requis")
	}

	// 2. Sauvegarder dans la Base de Données d'abord
	// (Si le pseudo existe déjà, la DB renverra une erreur ici)
	if err := database.GetDB().Create(client).Error; err != nil {
		return err
	}

	// 3. Créer l'utilisateur Système Linux
	err := s.userManager.CreateUser(client.Username, client.Password, client.ExpiryTime)
	if err != nil {
		// OUPS ! Échec création Linux -> On supprime l'entrée DB pour rester propre
		database.GetDB().Delete(client)
		return err
	}

	// 4. Initialiser la surveillance trafic (IPTables)
	s.trafficMon.ResetUserUsage(client.Username)

	return nil
}

// UpdateClient : Met à jour mot de passe ou expiration
func (s *ClientService) UpdateClient(id int, newPass string, newExpiry int64, enable bool) error {
	var client database.Client
	db := database.GetDB()

	// Récupérer le client
	if err := db.First(&client, id).Error; err != nil {
		return errors.New("client introuvable")
	}

	// Mise à jour Linux : Mot de passe
	if newPass != "" && newPass != client.Password {
		if err := s.userManager.UpdatePassword(client.Username, newPass); err != nil {
			return err
		}
		client.Password = newPass
	}

	// Mise à jour Linux : Expiration
	if newExpiry != client.ExpiryTime {
		if err := s.userManager.UpdateExpiry(client.Username, newExpiry); err != nil {
			return err
		}
		client.ExpiryTime = newExpiry
	}

	// Mise à jour Linux : Activation / Désactivation
	if enable != client.Enable {
		if enable {
			s.userManager.UnlockUser(client.Username)
		} else {
			s.userManager.LockUser(client.Username)
		}
		client.Enable = enable
	}

	// Sauvegarder les changements dans la DB
	return db.Save(&client).Error
}

// DeleteClient : Supprime de la DB et du Système
func (s *ClientService) DeleteClient(id int) error {
	var client database.Client
	db := database.GetDB()

	if err := db.First(&client, id).Error; err != nil {
		return errors.New("client introuvable")
	}

	// 1. Supprimer User Linux
	s.userManager.DeleteUser(client.Username)
	
	// 2. Supprimer Règle IPTables (Nettoyage)
	// On fait un "Reset" qui supprime la règle, mais on ne la recrée pas car l'user n'existe plus
	// (Note: ResetUserUsage supprime puis ajoute. Ici userdel a déjà nettoyé, mais on assure le coup)
	
	// 3. Supprimer de la DB
	return db.Delete(&client).Error
}

// GetAllClients : Récupère la liste pour l'affichage
func (s *ClientService) GetAllClients() ([]database.Client, error) {
	var clients []database.Client
	result := database.GetDB().Find(&clients)
	return clients, result.Error
}

// ResetTraffic : Remet le compteur à 0 (DB + IPTables)
func (s *ClientService) ResetTraffic(id int) error {
	var client database.Client
	db := database.GetDB()

	if err := db.First(&client, id).Error; err != nil {
		return err
	}

	// Reset DB
	client.Up = 0
	client.Down = 0
	db.Save(&client)

	// Reset IPTables
	return s.trafficMon.ResetUserUsage(client.Username)
}
