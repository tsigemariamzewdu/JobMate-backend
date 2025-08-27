package dto



type OTPRequestDTO struct {
    Phone string `json:"phone" binding:"required"`
}

type OTPResponseDTO struct {
    Message string `json:"message"`
}
