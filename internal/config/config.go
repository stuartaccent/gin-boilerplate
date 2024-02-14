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

// NewConfigFromPath creates and validates a new Config from a .toml file.
func NewConfigFromPath(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	if err := config.validate(); err != nil {
		return nil, err
	}
	return &config, nil
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
