package config

import (
	"os"
)

// AppConfig contient la configuration globale du panel
type AppConfig struct {
	ListenPort   string
	AdminUser    string
	AdminPass    string
	DBPath       string
}

// GlobalConfig stocke la config chargée
var GlobalConfig *AppConfig

// LoadConfig : Charge la configuration (Variables d'env ou Défaut)
func LoadConfig() {
	// Par défaut, on utilise ces valeurs si rien n'est configuré
	GlobalConfig = &AppConfig{
		ListenPort: getEnv("PANEL_PORT", "8080"),      // Le panel sera sur le port 8080
		AdminUser:  getEnv("ADMIN_USER", "admin"),     // Login par défaut
		AdminPass:  getEnv("ADMIN_PASS", "admin"),     // Mot de passe par défaut
		DBPath:     getEnv("DB_PATH", "/etc/utm-panel/utm.db"), // Où est stockée la DB
	}
}

// Petite fonction utilitaire pour lire les variables d'environnement
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
