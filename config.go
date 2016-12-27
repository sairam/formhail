package main

import (
	"os"
	"sync"
)

// AppConfig is
type AppConfig struct {
	WebsiteURL          string // https://website.com
	LocalServer         string // host:port combination used for starting the server
	TemplateDir         string // tmpl/
	TemplatePartialsDir string // tmpl/partials/

	once sync.Once
}

var config = new(AppConfig)

func init() {
	config.once.Do(func() { config.init() })
}

func (config *AppConfig) init() {
	config.WebsiteURL = os.Getenv("WEBSITE_URL")
	if config.WebsiteURL == "" {
		config.WebsiteURL = os.Getenv("LOCAL_SERVER")
	}
	config.LocalServer = os.Getenv("LOCAL_SERVER")

	config.TemplateDir = os.Getenv("TEMPLATE_DIR")
	config.TemplatePartialsDir = os.Getenv("TEMPLATE_PARTIALS_DIR")

}
