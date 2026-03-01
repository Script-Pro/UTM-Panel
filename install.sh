#!/bin/bash

# Couleurs pour l'affichage
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

echo -e "${GREEN}>>> Démarrage de l'installation de UTM-PANEL...${NC}"

# 1. Vérification Root
if [[ $EUID -ne 0 ]]; then
   echo -e "${RED}Erreur : Ce script doit être lancé en root ! (utilisez sudo -i)${NC}" 
   exit 1
fi

# 2. Installation des dépendances système
echo -e "\n--- 1/6 Installation des dépendances ---"
apt update && apt upgrade -y
# On installe git, wget, curl, et de quoi compiler
apt install -y git wget curl unzip tar socat build-essential

# 3. Installation de Go (Golang) - Version récente requise
echo -e "\n--- 2/6 Installation de Golang ---"
# On supprime toute vieille version
rm -rf /usr/local/go
# On télécharge Go 1.22
wget -q -O go.tar.gz https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
tar -C /usr/local -xzf go.tar.gz
rm go.tar.gz
# On ajoute Go au PATH système
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# 4. Clonage du projet depuis GitHub
echo -e "\n--- 3/6 Téléchargement du Panel ---"
APP_DIR="/etc/utm-panel"

# Si le dossier existe déjà, on le nettoie pour une install propre
if [ -d "$APP_DIR" ]; then
    rm -rf "$APP_DIR"
fi
mkdir -p "$APP_DIR"
cd "$APP_DIR"

# ---  LIEN GITHUB ---
git clone https://github.com/script-pro/utm-panel.git .
# ---------------------------------------

if [ ! -f "go.mod" ]; then
    echo -e "${RED}Erreur : Le clonage a échoué ou le dossier est vide.${NC}"
    exit 1
fi

# 5. Compilation du projet
echo -e "\n--- 4/6 Compilation du Code Source ---"
# Télécharge les librairies Go nécessaires
/usr/local/go/bin/go mod tidy
# Compile le fichier main.go en un exécutable nommé 'utm-panel'
/usr/local/go/bin/go build -o utm-panel cmd/server/main.go
chmod +x utm-panel

# 6. Création du Service Systemd (Démarrage auto)
echo -e "\n--- 5/6 Configuration du Service ---"
cat <<EOF > /etc/systemd/system/utm-panel.service
[Unit]
Description=UTM Panel Service (UDP/SSH/ZiVPN)
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=$APP_DIR
ExecStart=$APP_DIR/utm-panel
Restart=always
RestartSec=3
# Configuration par défaut
Environment="PANEL_PORT=8081"
Environment="ADMIN_USER=admin"
Environment="ADMIN_PASS=admin"
Environment="DB_PATH=$APP_DIR/utm.db"

[Install]
WantedBy=multi-user.target
EOF

# 7. Démarrage final
echo -e "\n--- 6/6 Démarrage ---"
systemctl daemon-reload
systemctl enable utm-panel
systemctl restart utm-panel

# Récupération de l'IP pour l'affichage
IP=$(curl -s ipv4.icanhazip.com)

echo -e "${GREEN}==============================================${NC}"
echo -e "${GREEN}      INSTALLATION TERMINÉE AVEC SUCCÈS !     ${NC}"
echo -e "${GREEN}==============================================${NC}"
echo -e ""
echo -e " 🔗 URL Panel : http://$IP:8081/login"
echo -e " 👤 User      : admin"
echo -e " 🔑 Password  : admin"
echo -e ""
echo -e "${GREEN}==============================================${NC}"
echo -e "Commande pour redémarrer le panel : systemctl restart utm-panel"
