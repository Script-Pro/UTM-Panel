package protocols

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type SlowDNSService struct {
	BinPath     string
	PrivKeyPath string
	PubKeyPath  string
}

func NewSlowDNSService() *SlowDNSService {
	cwd, _ := os.Getwd()
	return &SlowDNSService{
		BinPath:     filepath.Join(cwd, "bin", "dnstt-server"),
		PrivKeyPath: filepath.Join(cwd, "bin", "dnstt.key"),
		PubKeyPath:  filepath.Join(cwd, "bin", "dnstt.pub"),
	}
}

// Start : Configure et lance le serveur SlowDNS
// nameserver : Le sous-domaine NS (ex: ns1.tondomaine.com)
func (s *SlowDNSService) Start(nameserver string) error {
	if nameserver == "" {
		return fmt.Errorf("le nameserver est vide")
	}

	// 1. Télécharger le binaire DNSTT si absent
	if _, err := os.Stat(s.BinPath); os.IsNotExist(err) {
		fmt.Println("Téléchargement de DNSTT Server...")
		// On récupère un binaire compatible Linux AMD64
		cmd := exec.Command("wget", "-O", s.BinPath, "https://github.com/mizolinetech/udp-custom/raw/main/bin/dnstt-server")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("échec téléchargement dnstt: %v", err)
		}
		os.Chmod(s.BinPath, 0755)
	}

	// 2. Générer les clés (Privée et Publique) si absentes
	if _, err := os.Stat(s.PrivKeyPath); os.IsNotExist(err) {
		fmt.Println("Génération des clés SlowDNS...")
		// dnstt-server -gen-key retourne les clés dans la console, il faut les capturer
		out, err := exec.Command(s.BinPath, "-gen-key").Output()
		if err != nil {
			return fmt.Errorf("échec génération clés: %v", err)
		}
		
		// L'output ressemble à :
		// privkey: xxxx...
		// pubkey: yyyy...
		lines := strings.Split(string(out), "\n")
		var priv, pub string
		for _, line := range lines {
			if strings.HasPrefix(line, "privkey:") {
				priv = strings.TrimSpace(strings.TrimPrefix(line, "privkey:"))
			}
			if strings.HasPrefix(line, "pubkey:") {
				pub = strings.TrimSpace(strings.TrimPrefix(line, "pubkey:"))
			}
		}

		// Sauvegarde des clés dans des fichiers
		os.WriteFile(s.PrivKeyPath, []byte(priv), 0600)
		os.WriteFile(s.PubKeyPath, []byte(pub), 0644)
	}

	// 3. Créer le service Systemd
	// Le port UDP 5300 est standard pour DNSTT interne.
	// 127.0.0.1:22 signifie qu'on redirige le trafic vers le SSH local.
	serviceContent := fmt.Sprintf(`[Unit]
Description=SlowDNS DNSTT Service
After=network.target

[Service]
User=root
WorkingDirectory=%s
ExecStart=%s -udp :5300 -privkey-file %s %s 127.0.0.1:22
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target`, filepath.Dir(s.BinPath), s.BinPath, s.PrivKeyPath, nameserver)

	os.WriteFile("/etc/systemd/system/slowdns.service", []byte(serviceContent), 0644)

	// 4. Démarrer le service
	exec.Command("systemctl", "daemon-reload").Run()
	exec.Command("systemctl", "enable", "slowdns").Run()
	err := exec.Command("systemctl", "restart", "slowdns").Run()

	return err
}

// GetPublicKey : Récupère la clé publique pour l'afficher dans le Panel (Clients en ont besoin)
func (s *SlowDNSService) GetPublicKey() string {
	content, err := os.ReadFile(s.PubKeyPath)
	if err != nil {
		return "Clé non générée (Lancez le service d'abord)"
	}
	return string(content)
}

// Stop : Arrête SlowDNS
func (s *SlowDNSService) Stop() error {
	return exec.Command("systemctl", "stop", "slowdns").Run()
}
