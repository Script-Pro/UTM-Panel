package protocols

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ZiVPNConfig structure (basé sur ton fichier module/zvpn.json)
type ZiVPNConfig struct {
	Listen string      `json:"listen"`
	Cert   string      `json:"cert"`
	Key    string      `json:"key"`
	Obfs   string      `json:"obfs"`
	Auth   ZiVPNAuth   `json:"auth"`
}

type ZiVPNAuth struct {
	Mode   string   `json:"mode"`
	Config []string `json:"config"`
}

type ZiVPNService struct {
	BinPath    string
	ConfigPath string
	CertPath   string
	KeyPath    string
}

func NewZiVPNService() *ZiVPNService {
	cwd, _ := os.Getwd()
	// On stocke les certs dans le dossier bin/certs du projet pour rester propre
	certDir := filepath.Join(cwd, "bin", "certs")
	os.MkdirAll(certDir, 0755)

	return &ZiVPNService{
		BinPath:    filepath.Join(cwd, "bin", "zivpn"),
		ConfigPath: filepath.Join(cwd, "bin", "zvpn.json"),
		CertPath:   filepath.Join(certDir, "zivpn.crt"),
		KeyPath:    filepath.Join(certDir, "zivpn.key"),
	}
}

// Start : Installe, Génère SSL, Configure et Démarre ZiVPN
func (s *ZiVPNService) Start(port string) error {
	// 1. Télécharger le binaire si absent
	if _, err := os.Stat(s.BinPath); os.IsNotExist(err) {
		fmt.Println("Téléchargement de ZiVPN...")
		// Lien du binaire ZiVPN (Source standard compatible)
		cmd := exec.Command("wget", "-O", s.BinPath, "https://raw.githubusercontent.com/TARAPRO/ZIVPN-TUNNEL/main/zivpn")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("échec téléchargement zivpn: %v", err)
		}
		os.Chmod(s.BinPath, 0755)
	}

	// 2. Générer les certificats SSL (OBLIGATOIRE pour ZiVPN)
	if _, err := os.Stat(s.CertPath); os.IsNotExist(err) {
		fmt.Println("Génération du certificat SSL pour ZiVPN...")
		// Commande openssl pour créer un certificat auto-signé valide 10 ans
		cmd := exec.Command("openssl", "req", "-new", "-newkey", "rsa:2048", "-days", "3650", "-nodes", "-x509",
			"-subj", "/C=CM/ST=Center/L=Yaounde/O=UTM-Panel/CN=zivpn.com",
			"-keyout", s.KeyPath,
			"-out", s.CertPath)
		
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("échec génération SSL: %v", err)
		}
	}

	// 3. Créer le fichier zvpn.json
	config := ZiVPNConfig{
		Listen: ":" + port,
		Cert:   s.CertPath,
		Key:    s.KeyPath,
		Obfs:   "zivpn",
		Auth: ZiVPNAuth{
			Mode:   "passwords",
			Config: []string{"DOTYCAT"}, // Tag de connexion (identique à ton fichier d'origine)
		},
	}

	file, _ := json.MarshalIndent(config, "", "  ")
	if err := os.WriteFile(s.ConfigPath, file, 0644); err != nil {
		return fmt.Errorf("échec écriture config: %v", err)
	}

	// 4. Créer le service Systemd
	serviceContent := fmt.Sprintf(`[Unit]
Description=ZiVPN Tunnel Service
After=network.target

[Service]
User=root
WorkingDirectory=%s
ExecStart=%s server -c %s
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target`, filepath.Dir(s.BinPath), s.BinPath, s.ConfigPath)

	os.WriteFile("/etc/systemd/system/zivpn.service", []byte(serviceContent), 0644)

	// 5. Démarrer le service
	exec.Command("systemctl", "daemon-reload").Run()
	exec.Command("systemctl", "enable", "zivpn").Run()
	err := exec.Command("systemctl", "restart", "zivpn").Run()

	if err != nil {
		return fmt.Errorf("échec démarrage service: %v", err)
	}

	return nil
}

// Stop : Arrête le service
func (s *ZiVPNService) Stop() error {
	return exec.Command("systemctl", "stop", "zivpn").Run()
}
