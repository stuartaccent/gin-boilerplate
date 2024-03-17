package config

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type (
	SslMode    string
	ServerMode string
)

//goland:noinspection GoUnusedConst
const (
	SslModeDisable SslMode = "disable"
	SslModeAllow   SslMode = "allow"
	SslModePrefer  SslMode = "prefer"
	SslModeRequire SslMode = "require"

	ServerModeDebug   ServerMode = "debug"
	ServerModeRelease ServerMode = "release"
	ServerModeTest    ServerMode = "test"
)

// ToGinMode convert string to gin mode
func (m ServerMode) ToGinMode() string {
	switch m {
	case ServerModeDebug:
		return gin.DebugMode
	case ServerModeRelease:
		return gin.ReleaseMode
	case ServerModeTest:
		return gin.TestMode
	default:
		log.Printf("Invalid server mode '%s', falling back to '%s'", m, gin.ReleaseMode)
		return gin.ReleaseMode
	}
}

// Config represents the top-level configuration structure.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Security SecurityConfig `mapstructure:"security"`
	Session  SessionConfig  `mapstructure:"session"`
}

// FromPath creates and validates a new Config from a .toml file.
// If environment variables are set, they will override the values from the .toml file.
// The environment variables must be prefixed with the
// name of the configuration structure in uppercase, and the keys must be separated
// by underscores. For example, to override the `port` value in the `server` structure,
// the environment variable must be `SERVER_PORT`.
func FromPath(path string) (*Config, error) {
	config := &Config{}

	viper.SetConfigFile(path)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return config, nil
}

// ServerConfig represents the server configuration.
type ServerConfig struct {
	Port uint16     `mapstructure:"port"`
	Mode ServerMode `mapstructure:"mode"`
}

// DatabaseConfig represents the database configuration.
type DatabaseConfig struct {
	Host     string  `mapstructure:"host"`
	Port     uint16  `mapstructure:"port"`
	User     string  `mapstructure:"user"`
	Password string  `mapstructure:"password"`
	Db       string  `mapstructure:"db"`
	SslMode  SslMode `mapstructure:"ssl_mode"`
}

// SecurityConfig represents the security configuration.
type SecurityConfig struct {
	AllowedHosts          []string `mapstructure:"allowed_hosts"`
	StsSeconds            int64    `mapstructure:"sts_seconds"`
	StsIncludeSubdomains  bool     `mapstructure:"sts_include_subdomains"`
	FrameDeny             bool     `mapstructure:"frame_deny"`
	ContentTypeNosniff    bool     `mapstructure:"content_type_nosniff"`
	BrowserXSSFilter      bool     `mapstructure:"browser_xss_filter"`
	ContentSecurityPolicy string   `mapstructure:"content_security_policy"`
	CsrfSecret            string   `mapstructure:"csrf_secret"`
}

// SessionConfig represents the session configuration.
type SessionConfig struct {
	Key      string        `mapstructure:"key"`
	EncKey   string        `mapstructure:"enc_key"`
	Path     string        `mapstructure:"path"`
	Domain   string        `mapstructure:"domain"`
	MaxAge   int           `mapstructure:"max_age"`
	Secure   bool          `mapstructure:"secure"`
	HttpOnly bool          `mapstructure:"http_only"`
	SameSite http.SameSite `mapstructure:"same_site"`
}
