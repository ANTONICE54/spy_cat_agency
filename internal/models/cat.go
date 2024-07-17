package models

import "time"

type Cat struct {
	ID                uint      `json:"id"`
	Name              string    `json:"name" binding:"required,alpha"`
	YearsOfExperience uint      `json:"years_of_experience" binding:"required,numeric"`
	Breed             string    `json:"breed" binding:"required,alpha,breed"`
	Salary            float64   `json:"salary" binding:"required,numeric"`
	CreatedAt         time.Time `json:"created_at"`
}
