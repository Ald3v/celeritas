package celeritas

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ald3v/celeritas/logger"
	"github.com/ald3v/celeritas/session"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

const version = "1.0.0"

type Celeritas struct {
	AppName string
	Debug   bool
	Version string
	Logger *logger.Logger	
	RootPath string
	Routes *chi.Mux
	Session *scs.SessionManager
	DB Database
	config config
	AppConfig AppConfig
}

type config struct {
	port string
	renderer string
	cookie cookieConfig
	sessionType string
	database databaseConfig	
}

type AppConfig struct {
	limiter limiterConfig
}

func (c *Celeritas) New(rootPath string) error {
	/*pathConfig := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "data", "public", "tmp", "logs", "middleware"},
	}
	err := c.Init(pathConfig)
	if err != nil {
		return err
	}*/

	err := c.checkDotEnv(rootPath)
	if err != nil {
		return err
	}

	// read .env
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	// create loggers
	logger := logger.New(os.Stdout, logger.LevelInfo)

	// connect to database
	databaseConfig := databaseConfig{
		dbType: os.Getenv("DATABASE_TYPE"),
		dsn: c.BuildDSN(),
		maxOpenConns: 25,
		maxIdleConns: 25,
		maxIdleTime: "15m",
	}

	if os.Getenv("DATABASE_TYPE") != "" {
		db, err := c.OpenDB(databaseConfig)
		if err != nil {
			logger.PrintError(err,nil)
			os.Exit(1)
		}
		c.DB = Database {
			DataType: os.Getenv("DATABASE_TYPE"),
			Pool:db,
		}
	}


	c.Logger = logger	
	c.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	c.Version = version
	c.RootPath = rootPath
	c.Routes = c.routes().(*chi.Mux)

	c.config = config {
		port: os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
		cookie: cookieConfig{
			name:os.Getenv("COOKIE_NAME"),
			lifetime: os.Getenv("COOKIE_LIFETIME"),
			persist: os.Getenv("COOKIE_PERSISTS"),
			secure: os.Getenv("COOKIE_SECURE"),
			domain: os.Getenv("COOKIE_DOMAIN"),
		},
		sessionType: os.Getenv("SESION_TYPE"),
		database: databaseConfig,
	}


	limiterConfig := limiterConfig{
		enabled: true,
		rps: 2,
		burst:4,
	}

	c.AppConfig = AppConfig{
		limiter:limiterConfig,
	}

	sess := session.Session{
		CookieLifeTime: c.config.cookie.lifetime,
		CookiePersist: c.config.cookie.persist,
		CookieName: c.config.cookie.name,
		SessionType: c.config.sessionType,
		CookieDomain: c.config.cookie.domain,
	}

	c.Session = sess.InitSession()

	return nil
}

func (c *Celeritas) Init(p initPaths) error {
	root := p.rootPath
	for _, path := range p.folderNames {
		err := c.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Celeritas) ListenAndServe() {
	
	srv := &http.Server{
		Addr:fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog: log.New(c.Logger, "", 0),
		Handler:c.Routes,
		IdleTimeout: 30 * time.Second,
		ReadTimeout: 30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	defer c.DB.Pool.Close()

	c.Logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env": os.Getenv("ENV"),
	})
	err := srv.ListenAndServe()
	c.Logger.PrintFatal(err, nil)
}



func (c *Celeritas) checkDotEnv(path string) error {
	err := c.CreateFileIfNotExist(fmt.Sprintf("%s/.env",path))
	if err != nil {
		return err
	}

	return nil
}


func (c *Celeritas) BuildDSN() string{
	
	var dsn string
	switch os.Getenv("DATABASE_TYPE"){
	case "postgres","postgresql": dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",os.Getenv("DATABASE_HOST"),os.Getenv("DATABASE_PORT"),os.Getenv("DATABASE_USER"),os.Getenv("DATABASE_NAME"),os.Getenv("DATABASE_SSL_MODE"))
	if os.Getenv("DATABASE_PASS") != "" {
		dsn = fmt.Sprintf("%s password=%s", dsn, os.Getenv("DATABASE_PASS"))
	}
	default:
	}
	return dsn
}