package dto

type JobChatMessageDTO struct {
	Role    string `json:"role"`
	Message string `json:"message"`
}

type JobSuggestionRequest struct {
	UserID      string              `json:"user_id"`
	LookingFor  string              `json:"looking_for"` // "local", "remote", "freelance"
	Field       string              `json:"field"`
	Skills      []string            `json:"skills"`
	Experience  string              `json:"experience"`
	Language    string              `json:"language"` // "en", "am"
	ChatHistory []JobChatMessageDTO `json:"chat_history"`
}

type JobSuggestionResponse struct {
	Jobs    []JobDTO `json:"jobs"`
	Message string   `json:"message"`
}

type JobDTO struct {
	Title        string   `json:"title"`
	Company      string   `json:"company"`
	Location     string   `json:"location"`
	Requirements []string `json:"requirements"`
	Type         string   `json:"type"`
	Source       string   `json:"source"`
	Link         string   `json:"link"`
	Language     string   `json:"language"`
}
