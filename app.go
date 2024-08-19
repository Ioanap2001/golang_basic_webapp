package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/alexedwards/scs/v2"
)

type App struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	sessionManager *scs.SessionManager
	serverCfg      ServerConfig
	templateCache  map[string]*template.Template
}

type ServerConfig struct {
	Address string `json:"address"`
}

type AppConfig struct {
	Server ServerConfig
}

func loadAppConfig(filename string) AppConfig {
	configFile := filename
	var config AppConfig

	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error reading the configuration file %q: %v", filename, err)
	}

	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("Error parsing the configuration file %q: %v", filename, err)
	}

	return config
}

func initializeApp() (*App, error) {
	appCfg := loadAppConfig("config.json")
	sessionManager := scs.New()
	//sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true
	infoLog := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)

	templateCache, err := newTemplateCache()
	if err != nil {
		return nil, err
	}

	return &App{
		errorLog:       errorLog,
		infoLog:        infoLog,
		sessionManager: sessionManager,
		serverCfg:      appCfg.Server,
		templateCache:  templateCache,
	}, nil
}

func (app *App) getServer() *http.Server {
	srv := &http.Server{
		Addr:         app.serverCfg.Address,
		ErrorLog:     app.errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return srv
}
