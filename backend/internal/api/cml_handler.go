package api

import (
	"encoding/csv"
	"encoding/json"
	"mime/multipart"
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
	CMLData   map[string]interface{} `json:"cml_data" binding:"required"`
	Signature string                 `json:"signature" binding:"required"`
	PublicKey string                 `json:"public_key" binding:"required"`
}

type RefreshCMLRequest struct {
	CMLData   string `json:"cml_data" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

func (h *CMLHandler) UploadCML(c *gin.Context) {
	// Check if CSV file is being uploaded
	file, err := c.FormFile("csv_file")
	if err == nil && file != nil {
		// Handle CSV upload
		h.handleCSVUpload(c, file)
		return
	}

	// Handle JSON upload
	var req UploadCMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Marshal CMLData to JSON string
	cmlDataJSON, err := json.Marshal(req.CMLData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to marshal cml_data: " + err.Error()})
		return
	}

	cml, err := h.cmlService.UploadCML(string(cmlDataJSON), req.Signature, req.PublicKey)
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

func (h *CMLHandler) handleCSVUpload(c *gin.Context, fileHeader *multipart.FileHeader) {
	f, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to open uploaded file"})
		return
	}
	defer f.Close()

	// Parse CSV
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse CSV: " + err.Error()})
		return
	}

	if len(records) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CSV must have at least header and one data row"})
		return
	}

	// Process each row (skip header)
	uploadedCount := 0
	errors := []string{}
	
	for i := 1; i < len(records); i++ {
		row := records[i]
		if len(row) < 6 {
			errors = append(errors, "Row "+string(rune(i))+": insufficient columns")
			continue
		}

		// CSV format: org_id, max_sites, validity, feature_packs, dev_key_public, prod_key_public
		// Convert to JSON string
		cmlJSON := `{"org_id":"` + row[0] + `","max_sites":` + row[1] + `,"validity":"` + row[2] + `","feature_packs":"` + row[3] + `","key_type":"dev","issued_by":"system"}`
		
		// Upload with placeholder signature
		_, err := h.cmlService.UploadCML(cmlJSON, "default_signature", "default_public_key")
		if err != nil {
			errors = append(errors, "Row "+string(rune(i))+": "+err.Error())
			continue
		}
		
		uploadedCount++
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         "uploaded",
		"uploaded_count": uploadedCount,
		"total_rows":     len(records) - 1,
		"errors":         errors,
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
