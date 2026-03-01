package protocols

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Config JSON structure pour UDP Custom
type UDPConfig struct {
	Listen        string `json:"listen"`
	StreamBuffer  int    `json:"stream_buffer"`
	ReceiveBuffer int    `json:"receive_buffer"`
	Auth          Auth   `json:"auth"`
}

type Auth struct {
	Mode string `json:"mode"`
}

type UDPCustomService struct {
	BinPath    string
	ConfigPath string
}

func NewUDPCustomService() *UDPCustomService {
	// On définit où seront stockés les fichiers
	cwd, _ := os.Getwd()
	return &UDPCustomService{
		BinPath:    filepath.Join(cwd, "bin", "udp-custom"),
		ConfigPath: filepath.Join(cwd, "bin", "config.json"),
	}
}

// Start : Installe, Configure et Démarre UDP Custom
func (s *UDPCustomService) Start(port string) error {
	// 1. Vérifier si le binaire existe, sinon le télécharger
	if _, err := os.Stat(s.BinPath); os.IsNotExist(err) {
		fmt.Println("Téléchargement de UDP Custom...")
		// Lien direct vers le binaire compatible (basé sur ton ancien script)
		cmd := exec.Command("wget", "-O", s.BinPath, "https://github.com/mizolinetech/udp-custom/raw/main/bin/udp-custom-linux-amd64")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("échec du téléchargement: %v", err)
		}
		os.Chmod(s.BinPath, 0755)
	}

	// 2. Créer le fichier config.json
	config := UDPConfig{
		Listen:        ":" + port,
		StreamBuffer:  33554432,
		ReceiveBuffer: 83886080,
		Auth: Auth{
			Mode: "passwords", // Utilise les utilisateurs Linux qu'on a créés
		},
	}

	file, _ := json.MarshalIndent(config, "", "  ")
	if err := os.WriteFile(s.ConfigPath, file, 0644); err != nil {
		return fmt.Errorf("échec écriture config: %v", err)
	}

	// 3. Créer le service Systemd (pour qu'il tourne en arrière-plan)
	serviceContent := fmt.Sprintf(`[Unit]
Description=UDP Custom Service
After=network.target

[Service]
User=root
WorkingDirectory=%s
ExecStart=%s server -config %s
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target`, filepath.Dir(s.BinPath), s.BinPath, s.ConfigPath)

	os.WriteFile("/etc/systemd/system/udp-custom.service", []byte(serviceContent), 0644)

	// 4. Démarrer le service
	exec.Command("systemctl", "daemon-reload").Run()
	exec.Command("systemctl", "enable", "udp-custom").Run()
	err := exec.Command("systemctl", "restart", "udp-custom").Run()
	
	if err != nil {
		return fmt.Errorf("échec démarrage service: %v", err)
	}
	
	return nil
}

// Stop : Arrête le service
func (s *UDPCustomService) Stop() error {
	return exec.Command("systemctl", "stop", "udp-custom").Run()
}
