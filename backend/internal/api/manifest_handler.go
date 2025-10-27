package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"taskmaster-license/internal/service"
)

type ManifestHandler struct {
	manifestService *service.ManifestService
}

func NewManifestHandler(manifestService *service.ManifestService) *ManifestHandler {
	return &ManifestHandler{manifestService: manifestService}
}

type GenerateManifestRequest struct {
	Period string `json:"period" binding:"required"`
}

func (h *ManifestHandler) GenerateManifest(c *gin.Context) {
	var req GenerateManifestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	manifest, err := h.manifestService.GenerateManifest(req.Period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"manifest":  string(manifest.ManifestData),
		"signature": manifest.Signature,
	})
}

func (h *ManifestHandler) ListManifests(c *gin.Context) {
	period := c.Query("period")

	manifests, err := h.manifestService.ListManifests(period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"manifests": manifests})
}

func (h *ManifestHandler) GetManifest(c *gin.Context) {
	id := c.Param("manifest_id")

	manifest, err := h.manifestService.GetManifest(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"manifest": manifest})
}

func (h *ManifestHandler) DownloadManifest(c *gin.Context) {
	id := c.Param("manifest_id")

	manifest, err := h.manifestService.GetManifest(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=manifest_%s.json", id))
	c.String(http.StatusOK, string(manifest.ManifestData))
}

type SendManifestRequest struct {
	ManifestID  string `json:"manifest_id" binding:"required"`
	AStackEndpoint string `json:"astack_endpoint" binding:"required"`
}

func (h *ManifestHandler) SendManifest(c *gin.Context) {
	var req SendManifestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get manifest
	manifest, err := h.manifestService.GetManifest(req.ManifestID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// TODO: Send to A-Stack endpoint
	// For now, just mark as sent
	if err := h.manifestService.MarkManifestSent(req.ManifestID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "sent",
		"timestamp": manifest.SentAt,
	})
}
