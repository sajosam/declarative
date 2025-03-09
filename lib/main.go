package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

// Config struct to map YAML configuration
type Config struct {
	Server struct {
		Port           int      `yaml:"port"`
		Timeout        string   `yaml:"timeout"`
		AllowCORS      bool     `yaml:"allow_cors"`
		AllowedOrigins []string `yaml:"allowed_origins"`
		Cookie         struct {
			Secure   bool   `yaml:"secure"`
			HTTPOnly bool   `yaml:"http_only"`
			SameSite string `yaml:"same_site"`
		} `yaml:"cookie"`
	} `yaml:"server"`

	Framework string `yaml:"framework"`

	Middlewares struct {
		Global []string `yaml:"global"`
	} `yaml:"middlewares"`

	Routes struct {
		Groups []struct {
			Base   string `yaml:"base"`
			Routes []struct {
				Path       string `yaml:"path"`
				Handler    string `yaml:"handler"`
				Method     string `yaml:"method"`
				BodyParams []struct {
					Name string `yaml:"name"`
					Type string `yaml:"type"`
				} `yaml:"body_params"`
			} `yaml:"routes"`
		} `yaml:"groups"`
	} `yaml:"routes"`
}

// Load YAML config file
func loadConfig(configPath string) (*Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// StartServer initializes the Gin server using the YAML config
func StartServer(configPath string, handlerMap map[string]gin.HandlerFunc) {
	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize Gin engine
	r := gin.Default()

	// Apply middlewares
	if config.Server.AllowCORS {
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowOrigins = config.Server.AllowedOrigins
		r.Use(cors.New(corsConfig))
	}

	// Register routes from config
	for _, group := range config.Routes.Groups {
		api := r.Group(group.Base)
		for _, route := range group.Routes {
			handler, exists := handlerMap[route.Handler]
			if !exists {
				log.Fatalf("Handler not found for route: %s", route.Handler)
			}

			switch route.Method {
			case "GET":
				api.GET(route.Path, handler)
			case "POST":
				api.POST(route.Path, handler)
			case "PUT":
				api.PUT(route.Path, handler)
			case "DELETE":
				api.DELETE(route.Path, handler)
			default:
				log.Fatalf("Unsupported HTTP method: %s", route.Method)
			}
		}
	}

	// Start server
	port := fmt.Sprintf(":%d", config.Server.Port)
	log.Printf("Server running on %s", port)
	r.Run(port)
}
