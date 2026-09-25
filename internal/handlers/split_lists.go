package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sokosplit/sokosplit-core-api/internal/events"
	"github.com/sokosplit/sokosplit-core-api/internal/models"
	"gorm.io/gorm"
)

type SplitListHandler struct {
	DB    *gorm.DB
	Bus   *events.Bus
}

func NewSplitListHandler(db *gorm.DB, bus *events.Bus) *SplitListHandler {
	return &SplitListHandler{DB: db, Bus: bus}
}

type createSplitListRequest struct {
	SplitID    string `json:"split_id" binding:"required"`
	Token      string `json:"token" binding:"required"`
	Amount     int64  `json:"amount" binding:"required"`
	Recipients []struct {
		Address string `json:"address" binding:"required"`
		Bps     uint32 `json:"bps" binding:"required"`
	} `json:"recipients" binding:"required,min=1"`
}

// Create handles POST /split-lists — creates a split list definition.
// This records intent only; funding happens on-chain via
// sokosplit-wallet-service, which then calls back to confirm escrow.
func (h *SplitListHandler) Create(c *gin.Context) {
	var req createSplitListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var totalBps uint32
	for _, r := range req.Recipients {
		totalBps += r.Bps
	}
	if totalBps != 10_000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "recipient bps must sum to exactly 10000"})
		return
	}

	ownerID, _ := c.Get("user_id") // set by auth middleware

	split := models.SplitList{
		SplitID: req.SplitID,
		OwnerID: ownerID.(uuid.UUID),
		Token:   req.Token,
		Amount:  req.Amount,
		Status:  models.SplitStatusPending,
	}
	for _, r := range req.Recipients {
		split.Recipients = append(split.Recipients, models.Recipient{
			Address: r.Address,
			Bps:     r.Bps,
		})
	}

	if err := h.DB.Create(&split).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create split list"})
		return
	}

	c.JSON(http.StatusCreated, split)
}

// Get handles GET /split-lists/:id
func (h *SplitListHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var split models.SplitList
	if err := h.DB.Preload("Recipients").First(&split, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "split list not found"})
		return
	}
	c.JSON(http.StatusOK, split)
}

// List handles GET /split-lists — history for the authenticated owner.
func (h *SplitListHandler) List(c *gin.Context) {
	ownerID, _ := c.Get("user_id")
	var splits []models.SplitList
	if err := h.DB.Preload("Recipients").Where("owner_id = ?", ownerID).
		Order("created_at desc").Find(&splits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not list split lists"})
		return
	}
	c.JSON(http.StatusOK, splits)
}

// Release handles POST /split-lists/:id/release — marks a split ready for
// release and emits a release event for the wallet service to act on.
func (h *SplitListHandler) Release(c *gin.Context) {
	id := c.Param("id")
	var split models.SplitList
	if err := h.DB.First(&split, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "split list not found"})
		return
	}
	if split.Status == models.SplitStatusReleased {
		c.JSON(http.StatusConflict, gin.H{"error": "split already released"})
		return
	}

	if err := h.Bus.PublishReleaseFunds(events.ReleaseFundsEvent{SplitID: split.SplitID}); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "could not notify wallet service: " + err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "release requested"})
}

// WalletWebhook handles POST /webhooks/wallet-service — the wallet service
// calls back here once a deposit is confirmed on-chain or a release
// completes, so core-api can update its own record of status.
func WalletWebhook(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			SplitID string             `json:"split_id" binding:"required"`
			Status  models.SplitStatus `json:"status" binding:"required"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Model(&models.SplitList{}).
			Where("split_id = ?", payload.SplitID).
			Update("status", payload.Status).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update status"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "updated"})
	}
}
