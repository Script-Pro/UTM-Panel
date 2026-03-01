package system

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// UserManager gère les commandes système Linux
type UserManager struct{}

// CreateUser : Crée un utilisateur système sans accès shell (/bin/false)
func (m *UserManager) CreateUser(username string, password string, expiryTime int64) error {
	// 1. Vérifier si l'utilisateur existe déjà
	if _, err := exec.Command("id", username).CombinedOutput(); err == nil {
		return fmt.Errorf("l'utilisateur %s existe déjà", username)
	}

	// 2. Formater la date d'expiration pour Linux (YYYY-MM-DD)
	var expiryDate string
	if expiryTime > 0 {
		expiryDate = time.Unix(expiryTime, 0).Format("2006-01-02")
	}

	// 3. Préparer la commande useradd
	// -M : Pas de dossier home (pour garder le VPS propre)
	// -s /bin/false : Pas de connexion SSH shell (juste tunnel)
	args := []string{"-M", "-s", "/bin/false", username}
	
	if expiryDate != "" {
		args = append(args, "-e", expiryDate)
	}

	// Exécuter useradd
	cmd := exec.Command("useradd", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("erreur useradd: %s", string(out))
	}

	// 4. Définir le mot de passe
	if err := m.UpdatePassword(username, password); err != nil {
		// Nettoyage si le mot de passe échoue
		m.DeleteUser(username)
		return err
	}

	return nil
}

// UpdatePassword : Change le mot de passe Linux
func (m *UserManager) UpdatePassword(username string, password string) error {
	cmd := exec.Command("chpasswd")
	// On injecte "user:pass" dans l'entrée standard de la commande
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%s:%s", username, password))
	
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("erreur chpasswd: %s", string(out))
	}
	return nil
}

// DeleteUser : Supprime un utilisateur complètement
func (m *UserManager) DeleteUser(username string) error {
	// -f : Force la suppression même si connecté
	cmd := exec.Command("userdel", "-f", username)
	if out, err := cmd.CombinedOutput(); err != nil {
		// On ignore l'erreur si l'user n'existe pas déjà
		if strings.Contains(string(out), "does not exist") {
			return nil
		}
		return fmt.Errorf("erreur userdel: %s", string(out))
	}
	return nil
}

// UpdateExpiry : Change la date d'expiration
func (m *UserManager) UpdateExpiry(username string, expiryTime int64) error {
	dateStr := ""
	if expiryTime > 0 {
		dateStr = time.Unix(expiryTime, 0).Format("2006-01-02")
	}

	// Si dateStr est vide, cela enlève l'expiration (expire jamais)
	args := []string{"-E", dateStr, username}
	if dateStr == "" {
		args = []string{"-E", "-1", username}
	}

	cmd := exec.Command("chage", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("erreur chage: %s", string(out))
	}
	return nil
}

// LockUser : Bloque un utilisateur (quand quota dépassé)
func (m *UserManager) LockUser(username string) error {
	return exec.Command("usermod", "-L", username).Run()
}

// UnlockUser : Débloque un utilisateur
func (m *UserManager) UnlockUser(username string) error {
	return exec.Command("usermod", "-U", username).Run()
}
