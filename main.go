package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type TimeInfo struct {
	UTC       string `json:"utc"`
	KST       string `json:"kst"`
	Timestamp int64  `json:"timestamp"`
}

type ServerInfo struct {
	Hostname    string   `json:"hostname"`
	IPAddresses []string `json:"ip_addresses"`
	Time        TimeInfo `json:"time"`
}

func formatTime(t time.Time) TimeInfo {
	// UTC 시간
	utc := t.UTC().Format("2006-01-02 15:04:05 MST")

	// KST 시간
	kst := t.In(time.FixedZone("KST", 9*60*60)).Format("2006-01-02 15:04:05 MST")

	return TimeInfo{
		UTC:       utc,
		KST:       kst,
		Timestamp: t.Unix(),
	}
}

func getServerInfo() (*ServerInfo, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("호스트명을 가져오는 중 오류 발생: %v", err)
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("IP 주소를 가져오는 중 오류 발생: %v", err)
	}

	var ips []string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}

	return &ServerInfo{
		Hostname:    hostname,
		IPAddresses: ips,
		Time:        formatTime(time.Now()),
	}, nil
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	// 메소드 검증
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed")
		return
	}

	// 서버 정보 가져오기
	info, err := getServerInfo()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal server error: %v", err)
		return
	}

	// 응답 헤더 설정
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	// JSON 응답 전송
	if err := json.NewEncoder(w).Encode(info); err != nil {
		log.Printf("JSON 인코딩 중 오류 발생: %v", err)
	}
}

func main() {
	// /data 핸들러 등록
	http.HandleFunc("/data", dataHandler)

	// 서버 설정
	port := "8080"
	server := &http.Server{
		Addr: ":" + port,
		// 기본 타임아웃 설정
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 서버 실행
	log.Printf("서버가 포트 %s에서 실행 중입니다...", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("서버 실행 중 오류 발생: %v", err)
	}
}
