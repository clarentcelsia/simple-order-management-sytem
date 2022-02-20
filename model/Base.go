package model

import (
	"time"
)

type Base struct {
	CreatedAt string `json:"crt_at"`
	UpdatedAt string `json:"upd_at"`
}

func (base *Base) PrePersist() error {
	if base.CreatedAt == "" {
		base.CreatedAt = time.Now().Format(time.RFC3339Nano)
	}
	if base.UpdatedAt == "" {
		base.UpdatedAt = time.Now().Format(time.RFC3339Nano)
	}
	return nil
}

func (base *Base) PreUpdate() error {
	base.UpdatedAt = time.Now().Format(time.RFC3339Nano)
	return nil
}
