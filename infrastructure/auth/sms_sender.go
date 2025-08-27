package auth

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	config "github.com/tsigemariamzewdu/JobMate-backend/infrastructure/config"
	svc "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/services"
)

// TwilioOTPSender sends OTP via Twilio REST API (no SDK required)
type TwilioOTPSender struct {
	AccountSID string
	AuthToken  string
	From       string
	Client     *http.Client
}

var _ svc.IOTPSender = (*TwilioOTPSender)(nil)

// NewTwilioOTPSender returns a Twilio sender configured from cfg
func NewTwilioOTPSender(cfg *config.Config) (*TwilioOTPSender, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	if cfg.TwilioAccountSID == "" || cfg.TwilioAuthToken == "" || cfg.TwilioFromNumber == "" {
		return nil, fmt.Errorf("twilio configuration missing (TWILIO_ACCOUNT_SID/TWILIO_AUTH_TOKEN/TWILIO_FROM_NUMBER)")
	}
	return &TwilioOTPSender{
		AccountSID: cfg.TwilioAccountSID,
		AuthToken:  cfg.TwilioAuthToken,
		From:       cfg.TwilioFromNumber,
		Client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}, nil
}

// SendOTP sends the OTP message to the given phone (E.164 format, e.g. +2519xxxxxxx)
func (s *TwilioOTPSender) SendOTP(phone string, code string) error {
	if phone == "" {
		return fmt.Errorf("phone is empty")
	}
	message := fmt.Sprintf("JobMate: your verification code is %s", code)

	form := url.Values{}
	form.Set("From", s.From)
	form.Set("To", phone)
	form.Set("Body", message)

	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", s.AccountSID)
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	// Twilio uses HTTP Basic Auth (Account SID : Auth Token)
	req.SetBasicAuth(s.AccountSID, s.AuthToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("twilio request error: %w", err)
	}
	defer resp.Body.Close()

	// Twilio returns 201 Created on success; accept any 2xx.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("twilio send failed: status=%s body=%s", resp.Status, strings.TrimSpace(string(body)))
	}

	return nil
}

// StubOTPSender prints OTP to stdout — use this in development
type StubOTPSender struct{}

var _ svc.IOTPSender = (*StubOTPSender)(nil)

func (s *StubOTPSender) SendOTP(phone string, code string) error {
	fmt.Printf("[DEV SMS] OTP for %s: %s\n", phone, code)
	return nil
}

// NewOTPSenderFromEnv picks Twilio in production (with valid creds) else returns the dev stub.
// This avoids accidentally sending real SMS from development.
func NewOTPSenderFromEnv(cfg *config.Config) (svc.IOTPSender, error) {
	// Use Twilio only if explicitly in production and creds are present
	if cfg != nil && strings.ToLower(cfg.AppEnv) == "production" &&
		cfg.TwilioAccountSID != "" && cfg.TwilioAuthToken != "" && cfg.TwilioFromNumber != "" {

		if tw, err := NewTwilioOTPSender(cfg); err == nil {
			return tw, nil
		} else {
			// If Twilio init fails for some reason, return error so caller can handle or fallback
			return nil, fmt.Errorf("failed to initialize Twilio sender: %w", err)
		}
	}

	// Default: dev stub
	return &StubOTPSender{}, nil
}
