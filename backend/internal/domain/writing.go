package domain

import "time"

type WritingVariant struct {
	ID      uint          `json:"id" gorm:"primaryKey"`
	Name    string        `json:"name"`
	Tasks   []WritingTask `json:"tasks" gorm:"foreignKey:VariantID"`
	Version int32         `json:"version"`
}

type WritingTask struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	VariantID  uint   `json:"variant_id"`
	TaskNumber int    `json:"task_number"` // 1 or 2
	Prompt     string `json:"prompt"`
	ImageURL   string `json:"image_url,omitempty"`
	Version    int32  `json:"version"`
}

type WritingSubmission struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	TaskID      uint      `json:"task_id"`
	UserID      uint      `json:"user_id"`
	Answer      string    `json:"answer"`
	SubmittedAt time.Time `json:"submitted_at"`
	Version     int32     `json:"version"`
}
