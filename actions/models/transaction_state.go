package models

import "time"

type TransactionState struct {
	TransactionID uint64 `json:"transaction_id"`

	EntityID   uint64 `json:"entity_id"`
	EntityName string `json:"entity_name"`
	EntityType string `json:"entity_type"`

	State        string    `json:"state"`
	LastUpdateAt time.Time `json:"last_update_at"`

	CreatedByID uint64 `json:"created_by_id"`
	UpdatedByID uint64 `json:"updated_by_id"`
}
