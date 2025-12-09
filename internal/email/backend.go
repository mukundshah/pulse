package email

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	mail "github.com/wneessen/go-mail"
)

// Email represents an email message
type Email struct {
	To       string
	Subject  string
	HTMLBody string
	TextBody string
}

// Backend defines the interface for email backends
type Backend interface {
	SendEmail(ctx context.Context, email *Email) error
}

// SMTPBackend implements email sending via SMTP using go-mail
type SMTPBackend struct {
	client *mail.Client
	from   string
}

// NewSMTPBackend creates a new SMTP backend from URL
// Format: smtp://user:pass@host:port or smtp+tls:// or smtp+ssl://
func NewSMTPBackend(backendURL string, from string) (*SMTPBackend, error) {
	u, err := url.Parse(backendURL)
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP URL: %w", err)
	}

	port := 25
	if u.Port() != "" {
		p, err := strconv.Atoi(u.Port())
		if err != nil {
			return nil, fmt.Errorf("invalid port: %w", err)
		}
		port = p
	} else {
		// Default ports based on scheme
		switch u.Scheme {
		case "smtp+ssl":
			port = 465
		case "smtp+tls", "smtp":
			port = 587
		}
	}

	// Build client options
	opts := []mail.Option{
		mail.WithPort(port),
	}

	// Configure TLS/SSL
	switch u.Scheme {
	case "smtp+ssl":
		opts = append(opts, mail.WithSSLPort(false))
	case "smtp+tls":
		opts = append(opts, mail.WithTLSPolicy(mail.TLSMandatory))
	default:
		opts = append(opts, mail.WithTLSPolicy(mail.TLSOpportunistic))
	}

	// Configure authentication
	if u.User != nil {
		username := u.User.Username()
		password, _ := u.User.Password()
		if username != "" && password != "" {
			opts = append(opts, mail.WithSMTPAuth(mail.SMTPAuthPlain))
			opts = append(opts, mail.WithUsername(username))
			opts = append(opts, mail.WithPassword(password))
		}
	}

	client, err := mail.NewClient(u.Hostname(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create mail client: %w", err)
	}

	return &SMTPBackend{
		client: client,
		from:   from,
	}, nil
}

// SendEmail sends an email via SMTP
func (b *SMTPBackend) SendEmail(ctx context.Context, email *Email) error {
	msg := mail.NewMsg()

	if err := msg.From(b.from); err != nil {
		return fmt.Errorf("failed to set from address: %w", err)
	}

	if err := msg.To(email.To); err != nil {
		return fmt.Errorf("failed to set to address: %w", err)
	}

	msg.Subject(email.Subject)

	// Set both plain text and HTML bodies
	if email.TextBody != "" {
		msg.SetBodyString(mail.TypeTextPlain, email.TextBody)
	}
	if email.HTMLBody != "" {
		msg.AddAlternativeString(mail.TypeTextHTML, email.HTMLBody)
	}

	if err := b.client.DialAndSendWithContext(ctx, msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// Close closes the SMTP client connection
func (b *SMTPBackend) Close() error {
	return b.client.Close()
}

// ConsoleBackend implements email sending by logging to console
type ConsoleBackend struct{}

// NewConsoleBackend creates a new console backend
func NewConsoleBackend() *ConsoleBackend {
	return &ConsoleBackend{}
}

// SendEmail logs the email to console
func (b *ConsoleBackend) SendEmail(ctx context.Context, email *Email) error {
	log.Println("=" + strings.Repeat("=", 78))
	log.Printf("EMAIL (Console Backend)")
	log.Println("=" + strings.Repeat("=", 78))
	log.Printf("To: %s", email.To)
	log.Printf("Subject: %s", email.Subject)
	log.Println("-" + strings.Repeat("-", 78))
	log.Println("TEXT VERSION:")
	log.Println(email.TextBody)
	log.Println("-" + strings.Repeat("-", 78))
	log.Println("HTML VERSION:")
	log.Println(email.HTMLBody)
	log.Println("=" + strings.Repeat("=", 78))
	return nil
}

// FileBackend implements email sending by writing to files
type FileBackend struct {
	outputDir string
}

// NewFileBackend creates a new file backend
func NewFileBackend(backendURL string) (*FileBackend, error) {
	u, err := url.Parse(backendURL)
	if err != nil {
		return nil, fmt.Errorf("invalid filemail URL: %w", err)
	}

	outputDir := u.Path
	if outputDir == "" {
		outputDir = "./tmp/emails"
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	return &FileBackend{
		outputDir: outputDir,
	}, nil
}

// SendEmail writes the email to a file
func (b *FileBackend) SendEmail(ctx context.Context, email *Email) error {
	timestamp := time.Now().Format("20060102-150405")
	safeEmail := strings.NewReplacer(
		"@", "_at_",
		":", "_",
		"/", "_",
	).Replace(email.To)

	filename := fmt.Sprintf("email-%s-%s.txt", timestamp, safeEmail)
	filePath := filepath.Join(b.outputDir, filename)

	content := fmt.Sprintf("To: %s\n", email.To) +
		fmt.Sprintf("Subject: %s\n", email.Subject) +
		fmt.Sprintf("Time: %s\n", time.Now().Format(time.RFC3339)) +
		strings.Repeat("=", 80) + "\n" +
		"TEXT VERSION:\n" +
		strings.Repeat("-", 80) + "\n" +
		email.TextBody + "\n" +
		strings.Repeat("-", 80) + "\n" +
		"HTML VERSION:\n" +
		strings.Repeat("-", 80) + "\n" +
		email.HTMLBody + "\n" +
		strings.Repeat("=", 80) + "\n"

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write email file: %w", err)
	}

	log.Printf("Email written to file: %s", filePath)
	return nil
}

// LocMemBackend stores emails in memory (useful for testing)
type LocMemBackend struct {
	emails []EmailMessage
	mu     sync.RWMutex
}

// EmailMessage represents a stored email
type EmailMessage struct {
	To       string
	Subject  string
	HTMLBody string
	TextBody string
	SentAt   time.Time
}

// NewLocMemBackend creates a new in-memory backend
func NewLocMemBackend() *LocMemBackend {
	return &LocMemBackend{
		emails: make([]EmailMessage, 0),
	}
}

// SendEmail stores the email in memory
func (b *LocMemBackend) SendEmail(ctx context.Context, email *Email) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.emails = append(b.emails, EmailMessage{
		To:       email.To,
		Subject:  email.Subject,
		HTMLBody: email.HTMLBody,
		TextBody: email.TextBody,
		SentAt:   time.Now(),
	})

	log.Printf("Email stored in memory (total: %d)", len(b.emails))
	return nil
}

// GetEmails returns all stored emails
func (b *LocMemBackend) GetEmails() []EmailMessage {
	b.mu.RLock()
	defer b.mu.RUnlock()

	emails := make([]EmailMessage, len(b.emails))
	copy(emails, b.emails)
	return emails
}

// ClearEmails clears all stored emails
func (b *LocMemBackend) ClearEmails() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.emails = make([]EmailMessage, 0)
}

// DummyBackend does nothing (useful for testing)
type DummyBackend struct{}

// NewDummyBackend creates a new dummy backend
func NewDummyBackend() *DummyBackend {
	return &DummyBackend{}
}

// SendEmail does nothing
func (b *DummyBackend) SendEmail(ctx context.Context, email *Email) error {
	return nil
}

// NewBackendFromURL creates a backend based on the URL scheme
func NewBackendFromURL(backendURL string, from string) (Backend, error) {
	if backendURL == "" {
		return NewConsoleBackend(), nil
	}

	u, err := url.Parse(backendURL)
	if err != nil {
		return nil, fmt.Errorf("invalid email backend URL: %w", err)
	}

	switch u.Scheme {
	case "smtp", "smtp+tls", "smtp+ssl":
		return NewSMTPBackend(backendURL, from)
	case "consolemail":
		return NewConsoleBackend(), nil
	case "filemail":
		return NewFileBackend(backendURL)
	case "memorymail", "locmemmail":
		return NewLocMemBackend(), nil
	case "dummymail":
		return NewDummyBackend(), nil
	default:
		return nil, fmt.Errorf("unsupported email backend scheme: %s", u.Scheme)
	}
}
