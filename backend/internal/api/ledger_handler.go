package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"taskmaster-license/internal/repository"
)

type LedgerHandler struct {
	repo *repository.Repository
}

func NewLedgerHandler(repo *repository.Repository) *LedgerHandler {
	return &LedgerHandler{repo: repo}
}

func (h *LedgerHandler) GetLedger(c *gin.Context) {
	orgID := c.Query("org_id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "org_id required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "100")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	entries, total, err := h.repo.GetLedgerEntries(orgID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"entries": entries,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}
