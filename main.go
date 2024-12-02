package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// 환경변수 읽기
	multi := os.Getenv("multi")
	if multi == "" {
		log.Fatal("환경변수 'multi'가 설정되지 않았습니다.")
	}

	// /data 핸들러 정의
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Value: %s", multi)
	})

	// 서버 실행
	port := "8080"
	log.Printf("서버가 포트 %s에서 실행 중입니다...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("서버 실행 중 오류 발생: %v", err)
	}
}
