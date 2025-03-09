package loader

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port           int      `yaml:"port"`
		Timeout        string   `yaml:"timeout"`
		AllowCORS      bool     `yaml:"allow_cors"`
		AllowedOrigins []string `yaml:"allowed_origins"`
	} `yaml:"server"`

	Framework  string `yaml:"framework"`
	Middleware struct {
		Global []string `yaml:"global"`
	} `yaml:"middlewares"`

	Routes struct {
		Groups []RouteGroup `yaml:"groups"`
	} `yaml:"routes"`
}

type RouteGroup struct {
	Base   string  `yaml:"base"`
	Routes []Route `yaml:"routes"`
}

type Route struct {
	Path       string      `yaml:"path"`
	Handler    string      `yaml:"handler"`
	Method     string      `yaml:"method"`
	BodyParams *BodyParams `yaml:"body_params,omitempty"`
}

type BodyParams struct {
	Required []Param `yaml:"required,omitempty"`
}

type Param struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

var AppConfig Config

func LoadConfig() {
	file, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err = yaml.Unmarshal(file, &AppConfig)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
}
