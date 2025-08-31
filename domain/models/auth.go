package models

import (
	"time"
)

type RefreshToken struct {
	ID        	string    
	UserID    	string  
	TokenHash 	string    
	IsRevoked 	bool      
	ExpiresAt 	time.Time 
	CreatedAt 	time.Time 
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    time.Duration
	User         *User
}
