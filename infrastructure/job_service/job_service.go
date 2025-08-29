package job_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
)

type JobService struct {
	JobDataAPIKey string
}

func NewJobService(apiKey string) *JobService {
	return &JobService{JobDataAPIKey: apiKey}
}

func (s *JobService) GetCuratedJobs(field, lookingFor, experience string, skills []string, language string) ([]models.Job, string, error) {
	var jobs []models.Job
	// fetch from JobDataAPI for local jobs
	if lookingFor == "local" {
		localJobs, err := fetchJobsFromJobDataAPI(s.JobDataAPIKey, field)
		if err == nil && len(localJobs) > 0 {
			jobs = append(jobs, localJobs...)
		}
	}
	// fetch from Upwork for remote/freelance jobs
	if lookingFor == "remote" || lookingFor == "freelance" {
		upworkJobs, err := fetchUpworkJobs(field)
		if err == nil {
			jobs = append(jobs, upworkJobs...)
		}
	}
	if len(jobs) == 0 {
		userMsg := "No jobs found for your criteria. Please check your spelling, try related keywords, or broaden your search."
		if language == "am" {
			userMsg = "ምንም ስራ አልተገኘም። እባክዎ ፊደላትን ያረጋግጡ፣ ተዛማጅ ቃላት ይሞክሩ፣ ወይም ፍለጋዎን ያሰፋፉ።"
		}
		return nil, userMsg, errors.New("no jobs found for your criteria")
	}
	msg := "Here are some opportunities for you:"
	if language == "am" {
		msg = "እነዚህ ስራዎች ለ" + field + " የሚስማሙ ናቸው።"
	}
	return jobs, msg, nil
}

// fetch jobs from JobDataAPI for Ethiopia
func fetchJobsFromJobDataAPI(apiKey, titleFilter string) ([]models.Job, error) {
	resp, err := http.Get("https://jobdataapi.com/api/jobcountries/")
	if err != nil {
		return nil, fmt.Errorf("error fetching countries: %w", err)
	}
	defer resp.Body.Close()

	var countries []struct {
		Name string `json:"name"`
		Code string `json:"code"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&countries); err != nil {
		return nil, fmt.Errorf("error decoding countries: %w", err)
	}

	var countryCode string
	for _, c := range countries {
		if c.Name == "Ethiopia" {
			countryCode = c.Code
			break
		}
	}
	if countryCode == "" {
		return nil, fmt.Errorf("Ethiopia not found in countries list")
	}

	req, err := http.NewRequest("GET", "https://jobdataapi.com/api/jobs/", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Api-Key "+apiKey)

	q := req.URL.Query()
	q.Add("country_code", countryCode)
	if titleFilter != "" {
		q.Add("title", titleFilter)
	}
	req.URL.RawQuery = q.Encode()

	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching jobs: %w", err)
	}
	defer resp2.Body.Close()

	var data struct {
		Count   int `json:"count"`
		Results []struct {
			ID      int `json:"id"`
			Company struct {
				Name string `json:"name"`
				Logo string `json:"logo"`
			} `json:"company"`
			Title          string `json:"title"`
			Location       string `json:"location"`
			Description    string `json:"description"`
			ApplicationURL string `json:"application_url"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp2.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error decoding jobs: %w", err)
	}

	var jobs []models.Job
	for _, j := range data.Results {
		jobs = append(jobs, models.Job{
			Title:        j.Title,
			Company:      j.Company.Name,
			Location:     j.Location,
			Link:         j.ApplicationURL,
			Source:       "JobDataAPI",
			Requirements: []string{},
		})
	}

	return jobs, nil
}

// fetch jobs from Upwork for remote/freelance positions
func fetchUpworkJobs(field string) ([]models.Job, error) {
	searchURL := fmt.Sprintf("https://www.upwork.com/ab/jobs/search/?q=%s", url.QueryEscape(field))
	job := models.Job{
		Title:        fmt.Sprintf("Freelance %s Jobs", field),
		Company:      "Upwork",
		Location:     "Remote",
		Requirements: []string{},
		Type:         "freelance",
		Source:       "Upwork",
		Link:         searchURL,
		Language:     "en",
	}
	return []models.Job{job}, nil
}
