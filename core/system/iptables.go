package system

import (
	"os/exec"
	"strconv"
	"strings"
)

// TrafficMonitor : S'occupe de lire la consommation des données
type TrafficMonitor struct{}

// GetUserUsage : Récupère le nombre d'octets (Bytes) consommés par un utilisateur
// Cela correspond au téléchargement du client (Trafic SORTANT du serveur vers le client)
func (t *TrafficMonitor) GetUserUsage(username string) (int64, error) {
	// Exécute la commande : iptables -nvx -L OUTPUT
	// -n : Affiche les IP (pas de résolution DNS, plus rapide)
	// -v : Affiche les détails (packets, bytes)
	// -x : Affiche les chiffres exacts (pas de "K" ou "M", on veut les octets précis)
	cmd := exec.Command("iptables", "-nvx", "-L", "OUTPUT")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	// Analyse ligne par ligne
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		// On cherche la ligne qui contient "owner UID match" ET le nom de l'utilisateur
		// Exemple de ligne Linux : 
		// 500  45000  ACCEPT  all  --  * * 0.0.0.0/0  0.0.0.0/0  owner UID match monuser
		if strings.Contains(line, "owner UID match") && strings.Contains(line, username) {
			
			// On découpe la ligne par espaces
			fields := strings.Fields(line)
			
			// Le 2ème champ correspond aux BYTES (Octets)
			if len(fields) >= 2 {
				bytesStr := fields[1]
				bytes, err := strconv.ParseInt(bytesStr, 10, 64)
				if err == nil {
					return bytes, nil
				}
			}
		}
	}

	return 0, nil
}

// ResetUserUsage : Remet le compteur à zéro (utile lors du renouvellement)
func (t *TrafficMonitor) ResetUserUsage(username string) error {
	// Pour "reset", on supprime la règle et on la remet, c'est le plus simple avec iptables
	// 1. Supprimer (-D)
	exec.Command("iptables", "-D", "OUTPUT", "-m", "owner", "--uid-owner", username, "-j", "ACCEPT").Run()
	
	// 2. Remettre (-I) -> Le compteur repartira à 0
	err := exec.Command("iptables", "-I", "OUTPUT", "-m", "owner", "--uid-owner", username, "-j", "ACCEPT").Run()
	return err
}
