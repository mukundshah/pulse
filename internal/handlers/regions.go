package handlers

import (
	"net/http"

	"pulse/internal/store"

	"github.com/gin-gonic/gin"
)

type RegionHandler struct {
	store *store.Store
}

func NewRegionHandler(s *store.Store) *RegionHandler {
	return &RegionHandler{store: s}
}

// ListRegions handles GET /regions
func (h *RegionHandler) ListRegions(c *gin.Context) {
	regions, err := h.store.ListRegions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list regions"})
		return
	}

	c.JSON(http.StatusOK, regions)
}
