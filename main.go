package main

import (
	"embed"
	"fmt"
	"log"
	"minichat/config"
	"minichat/conversation"
	"minichat/server"
	"net/http"
	"strings"
)

//go:embed static
var DirStatic embed.FS

//go:embed templates/*
var DirTemplate embed.FS

// CORS middleware
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowed := false

		// 提取请求来源的主机名（去除协议和端口）
		host := extractHost(origin)

		for _, allowedOrigin := range config.GlobalConfig.CorsOrigins {
			log.Printf("Checking CORS origin: %s against allowed: %s\n", origin, allowedOrigin)
			if allowedOrigin == "*" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				allowed = true
				break
			}
			if allowedOrigin == origin || allowedOrigin == host {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				allowed = true
				break
			}
			// 支持 *.domain.com 通配符
			if strings.HasPrefix(allowedOrigin, "*.") {
				domain := strings.TrimPrefix(allowedOrigin, "*.")
				if strings.HasSuffix(host, "."+domain) {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					allowed = true
					break
				}
			}
		}

		if !allowed {
			log.Printf("Blocked CORS request from origin: %s\n", origin)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// extractHost 从完整的 URL 中提取主机名（去除协议和端口）
func extractHost(rawURL string) string {
	// 去除协议部分
	url := strings.TrimPrefix(rawURL, "http://")
	url = strings.TrimPrefix(url, "https://")

	// 截取主机名部分
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		url = parts[0]
	}

	// 去除端口部分
	parts = strings.Split(url, ":")
	if len(parts) > 0 {
		return parts[0]
	}

	return url
}

func main() {

	http.HandleFunc("/precheck", enableCORS(server.PreCheck))
	http.HandleFunc("/ws", enableCORS(server.HandleWs))
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		server.HandleFiles(writer, request, DirTemplate)
	})
	fs := http.FileServer(http.FS(DirStatic))
	http.Handle("/static/", fs)

	go conversation.Manager.Start()

	configVal := config.ParseConfig("config.yaml")

	log.Printf("\n\n********************************\nChat server is running at %d !\n********************************\n\n", configVal.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", configVal.Port), nil)
	if err != nil {
		fmt.Printf("Server start fail, error is: [ %+v ]", err)
		return
	}
}
