package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	dbHost     = "DB_HOST"
	dbPort     = "DB_PORT"
	dbUser     = "DB_USER"
	dbPassword = "DB_PASSWORD"
	dbName     = "DB_NAME"

	csrfKey = "CSRF_KEY"
	jwtKey  = "JWT_KEY"

	varIsEmpty = "%s is empty"
)

// Config represents config.
type Config struct {
	DB   DB
	CSRF CSRF
	JWT  JWT
}

// CSRF represents CSRF structure.
type CSRF struct {
	Key []byte
}

// JWT represents JWT structure.
type JWT struct {
	Key []byte
}

// DB represents DB config.
type DB struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

// New creates new config.
func New() (Config, error) {
	config := config()
	if err := config.validate(); err != nil {
		return Config{}, err
	}
	return config, nil
}

func (c *Config) validate() error {
	for _, validator := range []interface{ validate() error }{
		c.DB,
		c.CSRF,
		c.JWT,
	} {
		if err := validator.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (jwt JWT) validate() error {
	if len(jwt.Key) == 0 {
		return fmt.Errorf(varIsEmpty, jwtKey)
	}
	return nil
}

func (csrf CSRF) validate() error {
	if len(csrf.Key) == 0 {
		return fmt.Errorf(varIsEmpty, csrfKey)
	}
	return nil
}

func (db DB) validate() error {
	if db.Host == "" {
		return fmt.Errorf(varIsEmpty, dbHost)
	}
	if db.Port == 0 {
		return fmt.Errorf(varIsEmpty, dbPort)
	}
	if db.User == "" {
		return fmt.Errorf(varIsEmpty, dbUser)
	}
	if db.Password == "" {
		return fmt.Errorf(varIsEmpty, dbPassword)
	}
	if db.Name == "" {
		return fmt.Errorf(varIsEmpty, dbName)
	}
	return nil
}

func config() Config {
	return Config{
		DB: DB{
			Host:     os.Getenv(dbHost),
			Port:     getInt(os.Getenv(dbPort)),
			User:     os.Getenv(dbUser),
			Password: os.Getenv(dbPassword),
			Name:     os.Getenv(dbName),
		},
		CSRF: CSRF{
			Key: []byte(os.Getenv(csrfKey)),
		},
		JWT: JWT{
			[]byte(os.Getenv(jwtKey)),
		},
	}
}

func getInt(val string) int {
	res, err := strconv.Atoi(val)
	if err != nil {
		res = 0
	}
	return res
}
