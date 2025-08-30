package models

import (
	"time"
)

type JobChatMessage struct {
	Role      string    `bson:"role" json:"role"` 
	Message   string    `bson:"message" json:"message"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}

type JobChat struct {
	ID             string           `bson:"_id,omitempty" json:"id"`
	UserID         string           `bson:"user_id" json:"user_id"`
	Messages       []JobChatMessage `bson:"messages" json:"messages"`
	JobSearchQuery map[string]any   `bson:"job_search_query" json:"job_search_query"`
	JobResults     []Job            `bson:"job_results" json:"job_results"`
	CreatedAt      time.Time        `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time        `bson:"updated_at" json:"updated_at"`
}
