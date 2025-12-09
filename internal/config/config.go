package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
}

type ServerConfig struct {
	Port string
	Host string
	Mode string // debug, release, test
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	DSN      string
}

type LogConfig struct {
	Level     string // debug, info, warn, error
	Format    string // json, text
	Output    string // stdout, file, both
	LogFile   string // 请求日志文件路径
	AuditFile string // 审计日志文件路径
}

func LoadConfig() (*Config, error) {
	// 加载.env文件（如果存在）
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "testdb"),
		},
		Log: LogConfig{
			Level:     getEnv("LOG_LEVEL", "info"),
			Format:    getEnv("LOG_FORMAT", "text"),
			Output:    getEnv("LOG_OUTPUT", "stdout"),             // stdout, file, both
			LogFile:   getEnv("APP_LOG_FILE", "logs/app.log"),     // 请求日志文件路径
			AuditFile: getEnv("AUDIT_LOG_FILE", "logs/audit.log"), // 审计日志文件路径
		},
	}

	// 构建DSN
	config.Database.DSN = buildDSN(config.Database)

	return config, nil
}

func buildDSN(db DatabaseConfig) string {
	return db.User + ":" + db.Password + "@tcp(" + db.Host + ":" + db.Port + ")/" + db.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
