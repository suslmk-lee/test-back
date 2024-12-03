package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func dataHandler(w http.ResponseWriter, r *http.Request) {
	// 메소드 검증
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed")
		return
	}

	// 환경변수 읽기
	multi := os.Getenv("multi")
	if multi == "" {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal server error: environment variable 'multi' not set")
		return
	}

	// 응답 헤더 설정
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	// 응답 데이터 작성
	fmt.Fprintf(w, "Value: %s", multi)
}

func main() {
	// 환경변수 확인
	multi := os.Getenv("multi")
	if multi == "" {
		log.Fatal("환경변수 'multi'가 설정되지 않았습니다.")
	}

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
