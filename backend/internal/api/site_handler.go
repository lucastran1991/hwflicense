package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"taskmaster-license/internal/service"
)

type SiteHandler struct {
	siteService *service.SiteService
}

func NewSiteHandler(siteService *service.SiteService) *SiteHandler {
	return &SiteHandler{siteService: siteService}
}

type CreateSiteRequest struct {
	SiteID      string                 `json:"site_id" binding:"required"`
	Fingerprint map[string]interface{} `json:"fingerprint"`
}

func (h *SiteHandler) CreateSite(c *gin.Context) {
	var req CreateSiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID := c.GetString("org_id")
	if orgID == "" {
		orgID = c.Query("org_id")
	}

	site, err := h.siteService.CreateSiteLicense(orgID, req.SiteID, req.Fingerprint)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"license":   site,
		"site.lic":  string(site.LicenseData),
	})
}

func (h *SiteHandler) ListSites(c *gin.Context) {
	orgID := c.Query("org_id")
	status := c.Query("status")
	
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	sites, total, err := h.siteService.ListSiteLicenses(orgID, status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sites": sites,
		"total": total,
		"limit": limit,
		"offset": offset,
	})
}

func (h *SiteHandler) GetSite(c *gin.Context) {
	siteID := c.Param("site_id")

	site, err := h.siteService.GetSiteLicense(siteID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"license": site})
}

func (h *SiteHandler) DeleteSite(c *gin.Context) {
	siteID := c.Param("site_id")

	if err := h.siteService.RevokeSiteLicense(siteID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "revoked"})
}

func (h *SiteHandler) Heartbeat(c *gin.Context) {
	siteID := c.Param("site_id")

	if err := h.siteService.UpdateHeartbeat(siteID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

type ValidateLicenseRequest struct {
	License    map[string]interface{} `json:"license" binding:"required"`
	Fingerprint map[string]interface{} `json:"fingerprint"`
}

func (h *SiteHandler) Validate(c *gin.Context) {
	var req ValidateLicenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert license map to JSON string
	licenseData, err := c.MustGet("jsonEncoder").(func(interface{}) ([]byte, error))(req.License)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to encode license data"})
		return
	}

	result, err := h.siteService.ValidateLicense(string(licenseData), req.Fingerprint)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !result.Valid {
		c.JSON(http.StatusOK, gin.H{
			"valid":   false,
			"message": result.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":      true,
		"message":    "License valid",
		"features":   result.Features,
		"expires_at": result.ExpiresAt,
	})
}
