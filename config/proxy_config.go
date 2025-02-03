package config

type ProxyConfig struct {
    EnableProxy     bool   `json:"enable_proxy"`
    ProxyURL       string `json:"proxy_url"`
    ProxyTimeoutMS int    `json:"proxy_timeout_ms"`
}

func (c *Config) WithProxy() *Config {
    c.Proxy = ProxyConfig{
        EnableProxy:     true,
        ProxyURL:       "http://your-proxy-server:8080",
        ProxyTimeoutMS: 5000,
    }
    return c
} 