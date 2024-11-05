package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

var natEnabled bool

// 获取客户端的 IP 地址
func getClientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return strings.Split(forwarded, ",")[0]
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// 获取外网地址
func getExternalIP() (string, error) {
	cmd := exec.Command("curl", "-4", "ifconfig.me")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// NAT 穿透功能
func natPunchThrough(port int) {

	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		log.Fatalf("Failed to listen on UDP: %v", err)
	}
	defer conn.Close()

	// 获取并打印外网 IP 地址
	externalIP, err := getExternalIP()
	if err != nil {
		log.Printf("Failed to get external IP: %v", err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Printf("NAT punch-through enabled on local %s:%d, external IP: %s\n", localAddr.IP.String(), localAddr.Port, externalIP)

	// 发送一个空的数据包，保持连接
	_, err = conn.WriteToUDP([]byte("ping"), localAddr)
	if err != nil {
		log.Printf("Failed to send UDP packet: %v", err)
	}

	// 等待接收数据包
	buffer := make([]byte, 1024)
	for {
		_, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			break
		}
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := getClientIP(r)
	buffer := make([]byte, 1024*1024) // 1MB buffer
	totalBytes := 0

	start := time.Now() // 开始计时

	for {
		n, err := r.Body.Read(buffer)
		if n > 0 {
			totalBytes += n
			duration := time.Since(start).Seconds()
			speed := float64(totalBytes) / (1024 * 1024) / duration
			fmt.Printf("Client IP: %s - Received %d MB... Speed: %.2f MB/s\n", clientIP, totalBytes/(1024*1024), speed)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
	}

	totalMB := totalBytes / (1024 * 1024)
	fmt.Printf("Client IP: %s - Total received: %d MB in %.2f seconds. Average speed: %.2f MB/s\n",
		clientIP, totalMB, time.Since(start).Seconds(), float64(totalBytes)/(1024*1024)/time.Since(start).Seconds())

	w.Write([]byte("Upload complete"))
}

func main() {
	port := flag.Int("port", 18080, "Port to listen on")
	networkType := flag.String("t", "", "Network type (e.g., nat1 for NAT traversal)")

	flag.Parse()

	if *networkType == "nat1" {
		natEnabled = true
		go natPunchThrough(*port) // 启用 NAT 穿透
	}

	http.HandleFunc("/upload", uploadHandler) // 路由处理
	fmt.Printf("Server is listening on port %d...\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)) // 启动服务器
}
