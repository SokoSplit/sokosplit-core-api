package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User is a SokoSplit account — a freelancer, SME, or marketplace admin.
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	GitHubID  string    `gorm:"uniqueIndex" json:"github_id,omitempty"`
	WalletKey string    `json:"wallet_key,omitempty"` // Stellar public key, once linked
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// SplitStatus tracks a split list's lifecycle.
type SplitStatus string

const (
	SplitStatusPending  SplitStatus = "pending"
	SplitStatusEscrowed SplitStatus = "escrowed"
	SplitStatusReleased SplitStatus = "released"
)

// Recipient is one payee in a split list, matching the on-chain
// SplitEscrow.Recipient shape (address + basis points).
type Recipient struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SplitListID uuid.UUID `gorm:"type:uuid;index" json:"split_list_id"`
	Address     string    `gorm:"not null" json:"address"` // Stellar public key
	Bps         uint32    `gorm:"not null" json:"bps"`      // basis points, must sum to 10000 per split
}

func (r *Recipient) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// SplitList is the system-of-record definition of who gets paid what,
// and its current on-chain status. The `SplitID` field is the symbol
// passed to the SplitEscrow contract's deposit()/release()/get_split().
type SplitList struct {
	ID          uuid.UUID   `gorm:"type:uuid;primaryKey" json:"id"`
	SplitID     string      `gorm:"uniqueIndex;not null" json:"split_id"` // on-chain symbol
	OwnerID     uuid.UUID   `gorm:"type:uuid;index;not null" json:"owner_id"`
	Token       string      `gorm:"not null" json:"token"` // Stellar asset contract address
	Amount      int64       `gorm:"not null" json:"amount"`
	Recipients  []Recipient `gorm:"foreignKey:SplitListID" json:"recipients"`
	Status      SplitStatus `gorm:"not null;default:pending" json:"status"`
	ReleaseAt   *time.Time  `json:"release_at,omitempty"` // set if release condition is time-based
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

func (s *SplitList) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
