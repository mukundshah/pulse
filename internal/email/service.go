package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"time"

	"pulse/internal/config"
)

const (
	// EmailSendTimeout is the maximum time to wait for an email to send
	EmailSendTimeout = 30 * time.Second
)

// Service provides email sending functionality
type Service struct {
	backend Backend
	tmpl    *template.Template
	config  *config.Config
}

// NewService creates a new email service
func NewService(cfg *config.Config) (*Service, error) {

	backend, err := NewBackendFromURL(cfg.EmailURL, cfg.EmailFrom)
	if err != nil {
		return nil, fmt.Errorf("failed to create email backend: %w", err)
	}

	// Load email templates
	tmpl, err := loadTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to load email templates: %w", err)
	}

	return &Service{
		backend: backend,
		tmpl:    tmpl,
		config:  cfg,
	}, nil
}

// loadTemplates loads all email templates from the templates directory
func loadTemplates() (*template.Template, error) {
	templatesDir := "templates/email"

	// Load all .tmpl files in the templates directory
	tmpl := template.New("email").Funcs(template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"now": func() time.Time {
			return time.Now()
		},
	})

	// Find all template files
	pattern := filepath.Join(templatesDir, "*.tmpl")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to find template files: %w", err)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no template files found in %s", templatesDir)
	}

	// Parse all template files
	for _, match := range matches {
		_, err := tmpl.ParseFiles(match)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", match, err)
		}
	}

	return tmpl, nil
}

// SendPasswordResetEmail sends a password reset email synchronously
func (s *Service) SendPasswordResetEmail(ctx context.Context, to, resetToken string) error {
	// Build reset URL
	resetURL := fmt.Sprintf("%s/auth/password/reset/%s", s.config.FrontendURL, resetToken)

	// Prepare template data
	data := map[string]interface{}{
		"ResetURL": resetURL,
		"Email":    to,
		"Token":    resetToken,
	}

	// Render HTML template
	var htmlBuf bytes.Buffer
	if err := s.tmpl.ExecuteTemplate(&htmlBuf, "password_reset.html", data); err != nil {
		return fmt.Errorf("failed to render HTML template: %w", err)
	}

	// Render text template
	var textBuf bytes.Buffer
	if err := s.tmpl.ExecuteTemplate(&textBuf, "password_reset.txt", data); err != nil {
		return fmt.Errorf("failed to render text template: %w", err)
	}

	subject := "Reset Your Password"

	emailMsg := &Email{
		To:       to,
		Subject:  subject,
		HTMLBody: htmlBuf.String(),
		TextBody: textBuf.String(),
	}

	return s.backend.SendEmail(ctx, emailMsg)
}

// SendPasswordResetEmailAsync sends a password reset email asynchronously in a goroutine
// Errors are logged but not returned to the caller
func (s *Service) SendPasswordResetEmailAsync(ctx context.Context, to, resetToken string) {
	go func() {
		// Create a new context with timeout for the background operation
		// Don't use the parent context as it might be cancelled before email sends
		bgCtx, cancel := context.WithTimeout(context.Background(), EmailSendTimeout)
		defer cancel()

		if err := s.SendPasswordResetEmail(bgCtx, to, resetToken); err != nil {
			log.Printf("Failed to send password reset email to %s: %v", to, err)
		}
	}()
}

// SendEmailVerification sends an email verification email synchronously
func (s *Service) SendEmailVerification(ctx context.Context, to, email, verificationToken string) error {
	// Build verification URL - email is included in the token, so only token is needed
	// Format: /auth/verify-email?token={token}
	verificationURL := fmt.Sprintf("%s/auth/verify-email?token=%s", s.config.FrontendURL, verificationToken)

	// Prepare template data
	data := map[string]interface{}{
		"VerificationURL": verificationURL,
		"Email":           email,
		"Token":           verificationToken,
	}

	// Render HTML template
	var htmlBuf bytes.Buffer
	if err := s.tmpl.ExecuteTemplate(&htmlBuf, "email_verification.html", data); err != nil {
		return fmt.Errorf("failed to render HTML template: %w", err)
	}

	// Render text template
	var textBuf bytes.Buffer
	if err := s.tmpl.ExecuteTemplate(&textBuf, "email_verification.txt", data); err != nil {
		return fmt.Errorf("failed to render text template: %w", err)
	}

	subject := "Verify Your Email Address"

	emailMsg := &Email{
		To:       to,
		Subject:  subject,
		HTMLBody: htmlBuf.String(),
		TextBody: textBuf.String(),
	}

	return s.backend.SendEmail(ctx, emailMsg)
}

// SendEmailVerificationAsync sends an email verification email asynchronously in a goroutine
// Errors are logged but not returned to the caller
func (s *Service) SendEmailVerificationAsync(to, verificationToken string) {
	go func() {
		if err := s.SendEmailVerification(to, to, verificationToken); err != nil {
			log.Printf("Failed to send email verification to %s: %v", to, err)
		}
	}()
}
