package controller

import (
	"net/http"
	"strconv"
	"time"
	"utm-panel/config"   // <--- C'EST CETTE LIGNE QUI MANQUAIT
	"utm-panel/database"
	"utm-panel/service"

	"github.com/gin-gonic/gin"
)

type ClientController struct {
	clientService *service.ClientService
}

func NewClientController() *ClientController {
	return &ClientController{
		clientService: service.NewClientService(),
	}
}

// --- LOGIN ---
func (cc *ClientController) Login(c *gin.Context) {
	var form struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "msg": "Données invalides"})
		return
	}
	
	// Utilise la config globale pour vérifier le mot de passe
	if form.Username == config.GlobalConfig.AdminUser && form.Password == config.GlobalConfig.AdminPass {
		c.SetCookie("session", "logged_in", 3600*24, "/", "", false, false)
		c.JSON(http.StatusOK, gin.H{"success": true})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "msg": "Identifiants incorrects"})
	}
}

// --- AJOUT CLIENT ---
type AddRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	TotalGB    int    `json:"total_gb"`
	ExpiryDays int    `json:"expiry_days"`
	UDP        bool   `json:"udp"`
	ZiVPN      bool   `json:"zivpn"`
	SlowDNS    bool   `json:"slowdns"`
}

func (cc *ClientController) AddClient(c *gin.Context) {
	var req AddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "msg": err.Error()})
		return
	}

	totalBytes := int64(req.TotalGB) * 1024 * 1024 * 1024
	expiryTime := time.Now().AddDate(0, 0, req.ExpiryDays).Unix()

	newClient := database.Client{
		Username:     req.Username,
		Password:     req.Password,
		Total:        totalBytes,
		ExpiryTime:   expiryTime,
		Enable:       true,
		UDPCustom:    req.UDP,
		ZiVPN:        req.ZiVPN,
		SlowDNS:      req.SlowDNS,
		CreatedTime:  time.Now().Unix(),
	}

	err := cc.clientService.AddClient(&newClient)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "msg": "Erreur: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}

// --- LISTE ---
func (cc *ClientController) ListClients(c *gin.Context) {
	clients, err := cc.clientService.GetAllClients()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"obj": clients})
}

// --- SUPPRESSION ---
func (cc *ClientController) DeleteClient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "msg": "ID invalide"})
		return
	}

	err = cc.clientService.DeleteClient(id)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}

// --- RESET TRAFIC ---
func (cc *ClientController) ResetClient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "msg": "ID invalide"})
		return
	}
	cc.clientService.ResetTraffic(id)
	c.JSON(200, gin.H{"success": true})
}

// --- UPDATE ---
func (cc *ClientController) UpdateClient(c *gin.Context) {
	c.JSON(200, gin.H{"success": true})
}
