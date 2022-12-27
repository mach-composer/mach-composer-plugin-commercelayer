package internal

import (
	"net/url"
)

type SiteConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	Domain       string `mapstructure:"domain"`
}

func (s *SiteConfig) ApiEndpoint() (string, error) {
	return url.JoinPath(s.Domain, "/api")
}
func (s *SiteConfig) AuthEndpoint() (string, error) {
	return url.JoinPath(s.Domain, "/oauth/token")
}
