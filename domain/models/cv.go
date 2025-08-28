package models

import "time"



type CV struct {
    ID                  string                 
    UserID              string                 
    OriginalText        string                 
    ExtractedSkills     map[string]interface{} 
    ExtractedExperience map[string]interface{} 
    ExtractedEducation  map[string]interface{} 
    Summary             string
    Language            string                 
    IsActive            bool                  
    CreatedAt           time.Time              
    UpdatedAt           time.Time              

}

type CVFeedback struct {
    ID                    string    
    SessionID             string    
    UserID                string    
    CVID                  string    
    Strengths             string    
    Weaknesses            string    
    ImprovementSuggestions string   
    GeneratedAt           time.Time
}
