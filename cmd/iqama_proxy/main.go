package main

import (
	"fmt"
	"io"
	"net/http"

	v1 "github.com/ccil-kbw/robot/pkg/iqama/v1"
)

func main() {
	http.HandleFunc("/today", today)

	fmt.Println("watchtower test: Running iqama-proxy Go server on port :3333")
	_ = http.ListenAndServe(":3333", nil)
}

func today(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("request: %s %s ", r.Method, r.URL)
	io.WriteString(w, string(v1.GetRAW()))
}
