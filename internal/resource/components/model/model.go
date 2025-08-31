package model

import "time"

type Component struct {
	Id           uint64     `json:"id"`
	AddedBy      uint64     `json:"added_by"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Footprint    string     `json:"footprint"`
	Manufacturer string     `json:"manufacturer"`
	Supplier     string     `json:"supplier"`
	Amount       uint64     `json:"amount"`
	CreatedAt    *time.Time `json:"created_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}
