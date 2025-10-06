package config

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host string
	User string
	Pass string
	Port string
	Name string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	Database string
}

type JwtConfig struct {
	Secret string
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JwtConfig
}

var (
	instance *Config
	once     sync.Once
)

func EnvModule() *Config {
	once.Do(func() {
		_ = loadEnvFile(".env")

		instance = &Config{
			Server: ServerConfig{
				Port: getEnv("APP_PORT", "3000"),
			},
			Database: DatabaseConfig{
				Host: getEnv("DB_HOST", "localhost"),
				Port: getEnv("DB_PORT", "3309"),
				User: getEnv("DB_USER", "root"),
				Pass: getEnv("DB_PASS", ""),
				Name: getEnv("DB_NAME", "test"),
			},
			Redis: RedisConfig{
				Host:     getEnv("REDIS_HOST", "localhost"),
				Port:     getEnv("REDIS_PORT", "6379"),
				Password: getEnv("REDIS_PASSWORD", ""),
				Database: getEnv("REDIS_DATABASE", "0"),
			},
			JWT: JwtConfig{
				Secret: getEnv("JWT_SECRET", "ggwp1234"),
			},
		}
	})

	return instance
}

func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)

		os.Setenv(key, value)
	}

	return scanner.Err()
}

func getEnv(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
