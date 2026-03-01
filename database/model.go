package database

// Client : Structure utilisateur
type Client struct {
	ID           int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string `gorm:"unique;not null" json:"username"`
	Password     string `json:"password"`
	
	// Gestion de l'expiration et du quota
	ExpiryTime   int64  `json:"expiry_time"` 
	Total        int64  `json:"total"`       
	Up           int64  `json:"up"`          
	Down         int64  `json:"down"`        
	
	// État
	Enable       bool   `json:"enable"`      
	
	// Permissions
	UDPCustom    bool   `json:"udp_custom"`
	ZiVPN        bool   `json:"zivpn"`
	SlowDNS      bool   `json:"slowdns"`

	// Le champ qui manquait :
	CreatedTime  int64  `json:"created_time"`
}

// Setting : Configuration
type Setting struct {
	ID    int    `gorm:"primaryKey" json:"id"`
	Key   string `gorm:"unique" json:"key"`
	Value string `json:"value"`
}
