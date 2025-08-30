package models

type Job struct {
	Title        string
	Company      string
	Location     string
	Requirements []string
	Type         string      // "local", "remote", "freelance"
	Source       string      
	Link         string
	Language     string     
}
