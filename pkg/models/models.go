// Package models contains the data models used by the application.
package models

import "time"

// Trade represents a trade.
type Trade struct {
	ID           int64     `validate:"required" json:"id,omitempty"`
	CreatedAt    string    `validate:"required" json:"created_at,omitempty"`
	InstrumentID int64     `validate:"required" json:"instrument_id,omitempty"`
	Size         int64     `validate:"required" json:"size,omitempty"`
	Price        float64   `validate:"required" json:"price,omitempty"` // not a suitable money type, but ok for demo purposes
	Timestamp    time.Time `validate:"required" json:"timestamp,omitempty"`
}

// Position represents a position.
type Position struct {
	ID           int64     `validate:"required" json:"id,omitempty"`
	CreatedAt    string    `validate:"required" json:"created_at,omitempty"`
	InstrumentID int64     `validate:"required" json:"instrument_id,omitempty"`
	Size         int64     `validate:"required" json:"size,omitempty"`
	Timestamp    time.Time `validate:"required" json:"timestamp,omitempty"`
}
