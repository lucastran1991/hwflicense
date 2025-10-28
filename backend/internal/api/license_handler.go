package api

import (
	"net/http"
	"time"

	"taskmaster-license/internal/models"
	"taskmaster-license/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LicenseHandler struct {
	siteService  *service.LicenseSiteService
	statsService *service.LicenseStatsService
	alertService *service.LicenseAlertService
	jwtSecret    string
}

func NewLicenseHandler(siteService *service.LicenseSiteService, statsService *service.LicenseStatsService, 
	alertService *service.LicenseAlertService, jwtSecret string) *LicenseHandler {
	return &LicenseHandler{
		siteService:  siteService,
		statsService: statsService,
		alertService: alertService,
		jwtSecret:    jwtSecret,
	}
}

// CreateSiteKey handles POST /api/keys/create
func (h *LicenseHandler) CreateSiteKey(c *gin.Context) {
	var req struct {
		SiteID       string `json:"site_id" binding:"required"`
		EnterpriseID string `json:"enterprise_id" binding:"required"`
		Mode         string `json:"mode" binding:"required"`
		OrgID        string `json:"org_id" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Validate mode
	if req.Mode != "production" && req.Mode != "dev" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "mode must be 'production' or 'dev'"})
		return
	}
	
	siteKey, err := h.siteService.CreateSiteKey(req.SiteID, req.EnterpriseID, req.Mode, req.OrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, siteKey)
}

// ListSiteKeys handles GET /api/keys
func (h *LicenseHandler) ListSiteKeys(c *gin.Context) {
	enterpriseID := c.Query("enterprise_id")
	
	keys, err := h.siteService.ListSiteKeys(enterpriseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"keys":  keys,
		"total": len(keys),
	})
}

// UpdateSiteKey handles PUT /api/keys/:id
func (h *LicenseHandler) UpdateSiteKey(c *gin.Context) {
	siteID := c.Param("id")
	
	var req struct {
		Status string `json:"status"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if req.Status != "active" && req.Status != "revoked" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status must be 'active' or 'revoked'"})
		return
	}
	
	if err := h.siteService.UpdateSiteKeyStatus(siteID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"id":          siteID,
		"status":      req.Status,
		"updated_at":  time.Now(),
	})
}

// RefreshKey handles POST /api/keys/refresh
func (h *LicenseHandler) RefreshKey(c *gin.Context) {
	var req struct {
		SiteID string `json:"site_id" binding:"required"`
		OldKey string `json:"old_key" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	siteKey, err := h.siteService.RefreshSiteKey(req.SiteID, req.OldKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, siteKey)
}

// AggregateStats handles POST /api/stats/aggregate
func (h *LicenseHandler) AggregateStats(c *gin.Context) {
	var stats models.QuarterlyStats
	
	if err := c.ShouldBindJSON(&stats); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.statsService.SaveQuarterlyStats(&stats); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":    "saved",
		"period":    stats.Period,
		"timestamp": time.Now(),
	})
}

// ValidateKey handles POST /api/keys/validate
func (h *LicenseHandler) ValidateKey(c *gin.Context) {
	var req struct {
		SiteID string `json:"site_id" binding:"required"`
		Key    string `json:"key" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	siteKey, err := h.siteService.ValidateSiteKey(req.SiteID, req.Key)
	if err != nil {
		c.JSON(http.StatusOK, models.ValidationResponse{
			Valid:   false,
			Message: err.Error(),
		})
		return
	}
	
	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"site_id": siteKey.SiteID,
		"enterprise_id": siteKey.EnterpriseID,
		"key_type": siteKey.KeyType,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})
	
	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	
	c.JSON(http.StatusOK, models.ValidationResponse{
		Valid:     true,
		Token:     tokenString,
		ExpiresIn: 3600,
		Message:   "License valid",
	})
}

// SendAlert handles POST /api/alerts
func (h *LicenseHandler) SendAlert(c *gin.Context) {
	var alert models.Alert
	
	if err := c.ShouldBindJSON(&alert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.alertService.SaveAlert(&alert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":     "received",
		"site_id":    alert.SiteID,
		"alert_type": alert.AlertType,
		"timestamp":  alert.Timestamp,
	})
}

