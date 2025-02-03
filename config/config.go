package config

import (
	"os"
	"strconv"
)

type Config struct {
    Port     int         `json:"port"`
    LogLevel string      `json:"log_level"`
    Proxy    ProxyConfig `json:"proxy"`
    // 其他配置项...
}

func Load() (*Config, error) {
    config := &Config{
        Port:     8080,
        LogLevel: "info",
        Proxy: ProxyConfig{
            EnableProxy:     getEnvBool("ENABLE_PROXY", false),
            ProxyURL:       getEnv("PROXY_URL", ""),
            ProxyTimeoutMS: getEnvInt("PROXY_TIMEOUT_MS", 5000),
        },
    }
    return config, nil
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
    if value, exists := os.LookupEnv(key); exists {
        return value == "true"
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value, exists := os.LookupEnv(key); exists {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
} 