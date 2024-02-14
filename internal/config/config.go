package config

import (
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

type SslMode string

const (
	SslModeDisable SslMode = "disable"
	SslModeAllow   SslMode = "allow"
	SslModePrefer  SslMode = "prefer"
	SslModeRequire SslMode = "require"
)

// Config represents the top-level configuration structure.
type Config struct {
	Server   ServerConfig   `toml:"server"`
	Database DatabaseConfig `toml:"database"`
	Security SecurityConfig `toml:"security"`
	Session  SessionConfig  `toml:"session"`
}

// DefaultConfig creates a Config instance populated with default values.
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: 80,
			Mode: gin.ReleaseMode,
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "password",
			Database: "db",
			Sslmode:  SslModeDisable,
		},
		Security: SecurityConfig{
			AllowedHosts:          []string{},
			StsSeconds:            86400 * 365,
			StsIncludeSubdomains:  true,
			FrameDeny:             true,
			ContentTypeNosniff:    true,
			BrowserXSSFilter:      true,
			ContentSecurityPolicy: "default-src 'self'",
			CsrfSecret:            "some_secret_key",
		},
		Session: SessionConfig{
			Key:      "13d45bf0a822b832cc8886fa41ce4ced30584189bad02ec8ce552ace0d1ae8b1",
			EncKey:   "2bb61a68ac3dec4f7c25efb062f4ae3b",
			Path:     "/",
			Domain:   "",
			MaxAge:   86400 * 30,
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		},
	}
}

// NewConfigFromPath creates and validates a new Config from a .toml file.
// Uses a set of default values from DefaultConfig.
func NewConfigFromPath(path string) (*Config, error) {
	config := DefaultConfig()
	if _, err := toml.DecodeFile(path, config); err != nil {
		return nil, err
	}
	if err := config.validate(); err != nil {
		return nil, err
	}
	return config, nil
}

// ServerConfig represents the server configuration.
type ServerConfig struct {
	Port uint16 `toml:"port"`
	Mode string `toml:"mode" validate:"mode"`
}

// DatabaseConfig represents the database configuration.
type DatabaseConfig struct {
	Host     string  `toml:"host"`
	Port     uint16  `toml:"port"`
	User     string  `toml:"user"`
	Password string  `toml:"password"`
	Database string  `toml:"database"`
	Sslmode  SslMode `toml:"sslmode" validate:"sslmode"`
}

// SecurityConfig represents the security configuration.
type SecurityConfig struct {
	AllowedHosts          []string `toml:"allowed_hosts"`
	StsSeconds            int64    `toml:"sts_seconds"`
	StsIncludeSubdomains  bool     `toml:"sts_include_subdomains"`
	FrameDeny             bool     `toml:"frame_deny"`
	ContentTypeNosniff    bool     `toml:"content_type_nosniff"`
	BrowserXSSFilter      bool     `toml:"browser_xss_filter"`
	ContentSecurityPolicy string   `toml:"content_security_policy"`
	CsrfSecret            string   `toml:"csrf_secret"`
}

// SessionConfig represents the session configuration.
type SessionConfig struct {
	Key      string        `toml:"key"`
	EncKey   string        `toml:"enc_key"`
	Path     string        `toml:"path"`
	Domain   string        `toml:"domain"`
	MaxAge   int           `toml:"max_age"`
	Secure   bool          `toml:"secure"`
	HttpOnly bool          `toml:"http_only"`
	SameSite http.SameSite `toml:"same_site" validate:"samesite"`
}

// validate the config
func (s *Config) validate() error {
	validate := validator.New()
	if err := validate.RegisterValidation("mode", validateMode); err != nil {
		log.Fatal(err)
	}
	if err := validate.RegisterValidation("samesite", validateSameSite); err != nil {
		log.Fatal(err)
	}
	if err := validate.RegisterValidation("sslmode", validateSslMode); err != nil {
		log.Fatal(err)
	}
	return validate.Struct(s)
}

// validateMode ensures a valid gin mode.
func validateMode(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	switch val {
	case gin.DebugMode, gin.ReleaseMode, gin.TestMode:
		return true
	default:
		return false
	}
}

// validateSameSite ensures a valid http.SameSite.
func validateSameSite(fl validator.FieldLevel) bool {
	val := fl.Field().Int()
	switch http.SameSite(val) {
	case http.SameSiteDefaultMode, http.SameSiteLaxMode, http.SameSiteStrictMode, http.SameSiteNoneMode:
		return true
	default:
		return false
	}
}

// validateSslMode ensures a valid SslMode.
func validateSslMode(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	switch SslMode(val) {
	case SslModeDisable, SslModeAllow, SslModePrefer, SslModeRequire:
		return true
	default:
		return false
	}
}
