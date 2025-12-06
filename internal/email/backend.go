package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Backend defines the interface for email backends
type Backend interface {
	SendEmail(to, subject, htmlBody, textBody string) error
}

// SMTPBackend implements email sending via SMTP
type SMTPBackend struct {
	host     string
	port     int
	username string
	password string
	from     string
	useTLS   bool
	useSSL   bool
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

	username := ""
	password := ""
	if u.User != nil {
		username = u.User.Username()
		if p, ok := u.User.Password(); ok {
			password = p
		}
	}

	useTLS := u.Scheme == "smtp+tls"
	useSSL := u.Scheme == "smtp+ssl"

	return &SMTPBackend{
		host:     u.Hostname(),
		port:     port,
		username: username,
		password: password,
		from:     from,
		useTLS:   useTLS,
		useSSL:   useSSL,
	}, nil
}

// SendEmail sends an email via SMTP
func (b *SMTPBackend) SendEmail(to, subject, htmlBody, textBody string) error {
	addr := fmt.Sprintf("%s:%d", b.host, b.port)

	var auth smtp.Auth
	if b.username != "" && b.password != "" {
		auth = smtp.PlainAuth("", b.username, b.password, b.host)
	}

	// Create multipart message
	msg := []byte(fmt.Sprintf("From: %s\r\n", b.from) +
		fmt.Sprintf("To: %s\r\n", to) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: multipart/alternative; boundary=\"boundary123\"\r\n" +
		"\r\n" +
		"--boundary123\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		textBody +
		"\r\n" +
		"--boundary123\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		htmlBody +
		"\r\n" +
		"--boundary123--\r\n")

	if b.useSSL {
		// SSL connection (port 465)
		conn, err := tls.Dial("tcp", addr, &tls.Config{
			ServerName: b.host,
		})
		if err != nil {
			return fmt.Errorf("failed to connect via SSL: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, b.host)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
		defer client.Close()

		if auth != nil {
			if err := client.Auth(auth); err != nil {
				return fmt.Errorf("failed to authenticate: %w", err)
			}
		}

		if err := client.Mail(b.from); err != nil {
			return fmt.Errorf("failed to set sender: %w", err)
		}
		if err := client.Rcpt(to); err != nil {
			return fmt.Errorf("failed to set recipient: %w", err)
		}

		writer, err := client.Data()
		if err != nil {
			return fmt.Errorf("failed to get data writer: %w", err)
		}
		_, err = writer.Write(msg)
		if err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}
		err = writer.Close()
		if err != nil {
			return fmt.Errorf("failed to close writer: %w", err)
		}

		return client.Quit()
	}

	// Standard SMTP with optional STARTTLS
	if b.useTLS {
		return smtp.SendMail(addr, auth, b.from, []string{to}, msg)
	}

	// Try STARTTLS if available
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to dial SMTP server: %w", err)
	}
	defer client.Close()

	// Try STARTTLS
	if ok, _ := client.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: b.host}
		if err := client.StartTLS(config); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	}

	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
	}

	if err := client.Mail(b.from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	_, err = writer.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return client.Quit()
}

// ConsoleBackend implements email sending by logging to console
type ConsoleBackend struct{}

// NewConsoleBackend creates a new console backend
func NewConsoleBackend() *ConsoleBackend {
	return &ConsoleBackend{}
}

// SendEmail logs the email to console
func (b *ConsoleBackend) SendEmail(to, subject, htmlBody, textBody string) error {
	log.Println("=" + strings.Repeat("=", 78))
	log.Printf("EMAIL (Console Backend)")
	log.Println("=" + strings.Repeat("=", 78))
	log.Printf("To: %s", to)
	log.Printf("Subject: %s", subject)
	log.Println("-" + strings.Repeat("-", 78))
	log.Println("TEXT VERSION:")
	log.Println(textBody)
	log.Println("-" + strings.Repeat("-", 78))
	log.Println("HTML VERSION:")
	log.Println(htmlBody)
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
func (b *FileBackend) SendEmail(to, subject, htmlBody, textBody string) error {
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("email-%s-%s.txt", timestamp, strings.ReplaceAll(to, "@", "_at_"))
	filepath := filepath.Join(b.outputDir, filename)

	content := fmt.Sprintf("To: %s\n", to) +
		fmt.Sprintf("Subject: %s\n", subject) +
		fmt.Sprintf("Time: %s\n", time.Now().Format(time.RFC3339)) +
		strings.Repeat("=", 80) + "\n" +
		"TEXT VERSION:\n" +
		strings.Repeat("-", 80) + "\n" +
		textBody + "\n" +
		strings.Repeat("-", 80) + "\n" +
		"HTML VERSION:\n" +
		strings.Repeat("-", 80) + "\n" +
		htmlBody + "\n" +
		strings.Repeat("=", 80) + "\n"

	if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write email file: %w", err)
	}

	log.Printf("Email written to file: %s", filepath)
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
func (b *LocMemBackend) SendEmail(to, subject, htmlBody, textBody string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.emails = append(b.emails, EmailMessage{
		To:       to,
		Subject:  subject,
		HTMLBody: htmlBody,
		TextBody: textBody,
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
func (b *DummyBackend) SendEmail(to, subject, htmlBody, textBody string) error {
	// Do nothing
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
