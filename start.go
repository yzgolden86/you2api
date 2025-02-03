package main

import (
	"fmt"
	"log"
	"net/http"

	api "you2api/api" // 请替换为您的实际项目名
	config "you2api/config"
	proxy "you2api/proxy"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("运行错误: %v", err)
	}
}

func run() error {
	// 加载配置
	config, err := config.Load()
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 如果启用代理
	if config.Proxy.EnableProxy {
		proxy, err := proxy.NewProxy(config.Proxy.ProxyURL, config.Proxy.ProxyTimeoutMS)
		if err != nil {
			return fmt.Errorf("初始化代理失败: %w", err)
		}
		
		// 注册代理处理器
		http.Handle("/proxy/", http.StripPrefix("/proxy", proxy))
	}

	// 注册普通API处理器
	http.HandleFunc("/api/", api.Handler)

	port := fmt.Sprintf(":%d", config.Port)
	fmt.Printf("Server is running on http://localhost%s\n", port)

	// 启动服务器
	if err := http.ListenAndServe(port, nil); err != nil {
		return fmt.Errorf("启动服务器失败: %w", err)
	}
	return nil
}
