package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

type Config struct {
    Host string
    Port int
    User string
    Password string
    Dbname string
    Sslmode string
}

var DB *sql.DB

func NewConfig() *Config {
    port, _ := strconv.Atoi(getEnvOrDefault("DB_PORT", "5433"))
    return &Config{
        Host: getEnvOrDefault("DB_HOST", "localhost"),
        Port: port,
        User: getEnvOrDefault("DB_USER", "postgres"),
        Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
        Dbname: getEnvOrDefault("DB_NAME", "messaging_db"),
        Sslmode: getEnvOrDefault("DB_SSLMODE", "disable"),
    }
}

func (c *Config) ConnectionString() string {
    return fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        c.Host, c.Port, c.User,c.Password,c.Dbname, c.Sslmode,
    )
}

func Connect() error {
    log.Println("Connecting to database...")
    var err error
    DB, err = sql.Open("postgres", NewConfig().ConnectionString())
    if err != nil {
        return fmt.Errorf("failed to connect to database: %v", err)
    }

    err = DB.Ping()
    if err != nil {
        return fmt.Errorf("failed to ping database: %v", err)
    }
    
    return InitializeSchema()
}

func InitializeSchema() error {
    initSQL, err := os.ReadFile("init.sql")
    if err != nil {
        return fmt.Errorf("could not read init.sql: %v", err)
    }

    _, err = DB.Exec(string(initSQL))
    if err != nil {
        return fmt.Errorf("schema initialization failed: %v", err)
    }

    log.Println("Database connected and schema initialized")
    return nil
}

func getEnvOrDefault(key string, defaultValue string) string {
    value, exists := os.LookupEnv(key)
    if !exists {
        return defaultValue
    }
    return value
}

func Exec(query string, args ...interface{}) (sql.Result, error) {
    return DB.Exec(query, args...)
}

// returns a single row
func QueryRow(query string, args ...interface{}) *sql.Row {
    return DB.QueryRow(query, args...) 
}

// returns multiple rows
func Query(query string, args ...interface{}) (*sql.Rows, error) {
    return DB.Query(query, args...) 
}