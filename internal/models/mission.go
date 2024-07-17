package models

import "time"

type Mission struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name" binding:"required,alpha"`
	CatId       *uint     `json:"cat_id"`
	TargetList  []Target  `json:"target_list" binding:"required"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
}

type Target struct {
	ID          uint      `json:"id"`
	MissionID   uint      `json:"mission_id" `
	Name        string    `json:"name" binding:"required,alpha"`
	Country     string    `json:"country" binding:"required,alpha"`
	Notes       string    `json:"notes"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
}
