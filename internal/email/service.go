package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"path/filepath"

	"pulse/internal/config"
)

// Service handles email sending with template rendering
type Service struct {
	backend Backend
	tmpl    *template.Template
	config  *config.Config
}

// NewService creates a new email service
// Uses EMAIL_BACKEND URL to determine which backend to use
func NewService(cfg *config.Config) (*Service, error) {
	backendURL := cfg.EmailURL
	if backendURL == "" {
		backendURL = "consolemail://"
	}

	backend, err := NewBackendFromURL(backendURL, cfg.EmailFrom)
	if err != nil {
		return nil, fmt.Errorf("failed to create email backend: %w", err)
	}

	log.Printf("Email service: Using backend %s", backendURL)

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
func (s *Service) SendPasswordResetEmail(to, resetToken string) error {
	// Build reset URL
	resetURL := fmt.Sprintf("%s/auth/password/reset/%s", s.config.FrontendURL, resetToken)

	// Prepare template data
	data := map[string]interface{}{
		"ResetURL": resetURL,
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

	subject := "Password Reset Request"

	return s.backend.SendEmail(to, subject, htmlBuf.String(), textBuf.String())
}

// SendPasswordResetEmailAsync sends a password reset email asynchronously in a goroutine
// Errors are logged but not returned to the caller
func (s *Service) SendPasswordResetEmailAsync(to, resetToken string) {
	go func() {
		if err := s.SendPasswordResetEmail(to, resetToken); err != nil {
			log.Printf("Failed to send password reset email to %s: %v", to, err)
		}
	}()
}
