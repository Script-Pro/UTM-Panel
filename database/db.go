package database

import (
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB : Lance la connexion à la base de données
func InitDB(dbPath string) error {
	// Créer le dossier si nécessaire
	dir := filepath.Dir(dbPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	// AutoMigrate crée les tables automatiquement selon le fichier model.go
	err = DB.AutoMigrate(&Client{}, &Setting{})
	if err != nil {
		return err
	}

	log.Println("Base de données initialisée avec succès !")
	return nil
}

// GetDB : Récupère l'instance de la DB
func GetDB() *gorm.DB {
	return DB
}
