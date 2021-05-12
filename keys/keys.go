// Package keys provides key management functions.
package keys

import (
	"context"
	"time"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"gorm.io/gorm"
)

const (
	ACCOUNT_KEY_TYPE_LOCAL      = "local"
	ACCOUNT_KEY_TYPE_GOOGLE_KMS = "google_kms"
)

// Manager provides the functions needed for key management.
type Manager interface {
	// Generate generates a new Key using provided key index and weight.
	Generate(ctx context.Context, keyIndex int, weight int) (Wrapped, error)
	// GenerateDefault generates a new Key using application defaults.
	GenerateDefault(context.Context) (Wrapped, error)
	// Save is responsible for converting an "in flight" key to a storable key.
	Save(Key) (StorableKey, error)
	// Load is responsible for converting a storable key to an "in flight" key.
	Load(StorableKey) (Key, error)
	// AdminAuthorizer returns an Authorizer for the applications admin account.
	AdminAuthorizer(context.Context) (Authorizer, error)
	// UserAuthorizer returns an Authorizer for the given address.
	UserAuthorizer(ctx context.Context, address string) (Authorizer, error)
}

// StorableKey struct represents a storable account key.
// StorableKey.Value is an encrypted byte representation of
// the actual private key when using local key management
// or resource id when using a remote key management system (e.g. Google KMS).
type StorableKey struct {
	ID             int            `json:"-" gorm:"primaryKey"`
	AccountAddress string         `json:"-" gorm:"index"`
	Index          int            `json:"index" gorm:"index"`
	Type           string         `json:"type"`
	Value          []byte         `json:"-"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// Key is an "in flight" account key meaning its Value should be the actual
// private key or resource id (unencrypted).
type Key struct {
	Index int    `json:"index"`
	Type  string `json:"type"`
	Value string `json:"-"`
}

// Authorizer groups the necessary items for transaction signing.
type Authorizer struct {
	Address flow.Address
	Key     *flow.AccountKey
	Signer  crypto.Signer
}

// Wrapped simply provides a way to pass a flow.AccountKey and the corresponding Key together.
type Wrapped struct {
	FlowKey    *flow.AccountKey
	AccountKey Key
}
