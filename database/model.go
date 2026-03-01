package database

import (
	"gorm.io/gorm"
)

// Client : C'est la structure qui représente un utilisateur dans ton panel
type Client struct {
	ID           int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string `gorm:"unique;not null" json:"username"`
	Password     string `json:"password"`
	
	// Gestion de l'expiration et du quota
	ExpiryTime   int64  `json:"expiry_time"` // Date d'expiration (Timestamp Unix)
	Total        int64  `json:"total"`       // Quota total en Octets (Bytes)
	Up           int64  `json:"up"`          // Upload consommé
	Down         int64  `json:"down"`        // Download consommé
	
	// État
	Enable       bool   `json:"enable"`      // Le compte est-il actif ?
	
	// Permissions des protocoles (Cocher les cases dans le panel)
	UDPCustom    bool   `json:"udp_custom"`
	ZiVPN        bool   `json:"zivpn"`
	SlowDNS      bool   `json:"slowdns"`
}

// Setting : Pour sauvegarder les ports et configurations globales
type Setting struct {
	ID    int    `gorm:"primaryKey" json:"id"`
	Key   string `gorm:"unique" json:"key"`   // Ex: "udp_port", "zivpn_port"
	Value string `json:"value"`               // Ex: "36712", "5667"
}
