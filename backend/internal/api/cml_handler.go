package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"taskmaster-license/internal/service"
)

type CMLHandler struct {
	cmlService *service.CMLService
}

func NewCMLHandler(cmlService *service.CMLService) *CMLHandler {
	return &CMLHandler{cmlService: cmlService}
}

type UploadCMLRequest struct {
	CMLData   string `json:"cml_data" binding:"required"`
	Signature string `json:"signature" binding:"required"`
	PublicKey string `json:"public_key" binding:"required"`
}

type RefreshCMLRequest struct {
	CMLData   string `json:"cml_data" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

func (h *CMLHandler) UploadCML(c *gin.Context) {
	var req UploadCMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cml, err := h.cmlService.UploadCML(req.CMLData, req.Signature, req.PublicKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "valid",
		"org_id": cml.OrgID,
		"cml":    cml,
	})
}

func (h *CMLHandler) GetCML(c *gin.Context) {
	orgID := c.Query("org_id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "org_id required"})
		return
	}

	cml, err := h.cmlService.GetCML(orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cml":    cml,
		"status": "active",
	})
}

func (h *CMLHandler) RefreshCML(c *gin.Context) {
	orgID := c.Query("org_id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "org_id required"})
		return
	}

	var req RefreshCMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.cmlService.RefreshCML(orgID, req.CMLData, req.Signature); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "refreshed"})
}
