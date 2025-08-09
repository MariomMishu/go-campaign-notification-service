package config

import "time"

type AppConfig struct {
	Name            string
	Port            string
	NumberOfWorkers int
}
type DbConfig struct {
	Host            string
	Port            string
	User            string
	Pass            string
	Schema          string
	MaxIdleConn     int
	MaxOpenConn     int
	MaxConnLifetime time.Duration
	Debug           bool
}
type RedisConfig struct {
	Host               string
	Port               string
	Pass               string
	Db                 int
	MandatoryPrefix    string
	AccessUuidPrefix   string
	RefreshUuidPrefix  string
	UserPrefix         string
	PermissionPrefix   string
	UserCacheTTL       time.Duration
	PermissionCacheTTL time.Duration
}
type AsynqConfig struct {
	RedisAddr                   string
	DB                          int
	Pass                        string
	Concurrency                 int
	Queue                       string
	Retention                   time.Duration
	RetryCount                  int
	Delay                       time.Duration
	EmailSendTaskDelay          time.Duration
	EmailSendTaskRetryCount     int
	EmailSendTaskRetryDelay     time.Duration
	ReminderTaskRetryCount      int
	ReminderTaskRetryDelay      time.Duration
	ReminderEmailTaskRetryCount int
	ReminderEmailTaskRetryDelay time.Duration
	ReminderEmailTaskDelay      time.Duration
}
type JwtConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}
type EmailConfig struct {
	//Url     string
	//Timeout time.Duration
	Host     string        // e.g., smtp.gmail.com
	Port     string        // e.g., 587
	Username string        // e.g., your email address
	Password string        // e.g., app password or real password
	Timeout  time.Duration // optional
}
type Config struct {
	App   *AppConfig
	DB    *DbConfig
	Redis *RedisConfig
	Jwt   *JwtConfig
	Email *EmailConfig
	Asynq *AsynqConfig
}

var config Config

func GetAll() Config {
	return config
}
func App() *AppConfig {
	return config.App
}
func Db() *DbConfig {
	return config.DB
}
func Redis() *RedisConfig {
	return config.Redis
}
func Jwt() *JwtConfig {
	return config.Jwt
}
func Email() *EmailConfig { return config.Email }
func Asynq() *AsynqConfig {
	return config.Asynq
}

func LoadConfig() {
	setDefaultConfig()
}
func setDefaultConfig() {
	config.App = &AppConfig{
		Name:            "EMS",
		Port:            "8081",
		NumberOfWorkers: 5,
	}
	config.DB = &DbConfig{
		Host:            "127.0.0.1",
		Port:            "3306",
		User:            "mariom",
		Pass:            "password",
		Schema:          "campaign_notification_service",
		MaxIdleConn:     1,
		MaxOpenConn:     2,
		MaxConnLifetime: 30,
		Debug:           true,
	}
	config.Redis = &RedisConfig{
		Host:            "127.0.0.1",
		Port:            "6379",
		Pass:            "",
		Db:              0,
		MandatoryPrefix: "campaign_notification_",
	}
	config.Asynq = &AsynqConfig{
		RedisAddr:   "127.0.0.1:6379",
		DB:          0,
		Pass:        "",
		Concurrency: 10,
		Queue:       "app",
		Retention:   168,
		RetryCount:  5,
		Delay:       120,
	}
	// Initialize JwtConfig here:
	config.Jwt = &JwtConfig{
		AccessTokenSecret:  "your-access-secret",
		RefreshTokenSecret: "your-refresh-secret",
		// expiry in seconds or whatever unit you prefer
		AccessTokenExpiry:  time.Hour,      // e.g. 1 hour
		RefreshTokenExpiry: 24 * time.Hour, // e.g. 24 hours
	}
	config.Email = &EmailConfig{
		Host:     "smtp.gmail.com",
		Port:     "587",
		Username: "mishu.cste08@gmail.com",
		Password: "guld gosd rxnx pkjb",
		Timeout:  10 * time.Second,
	}
}
