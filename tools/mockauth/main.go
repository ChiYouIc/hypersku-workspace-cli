package main
// mockauth 是开发调试用的 Device Code Flow 模拟认证服务。
//
// 用途：在真实认证服务就绪前，验证 hypersku-cli auth login/status 的完整流程。
//
// 端点：
//   POST /oauth/device/code  申请设备码（RFC 8628 §3.1）
//   GET  /verify             模拟用户在浏览器完成授权
//   POST /oauth/token        设备码换 token（授权前返回 authorization_pending）
//
// 用法：
//
//	go run ./tools/mockauth            # 监听 127.0.0.1:9999
//	go run ./tools/mockauth -addr :8080
//
// 然后配置 ~/.hypersku-cli/config.json:
//
//	{ "auth": { "base_url": "http://127.0.0.1:9999" } }
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type deviceSession struct {
	code     string
	approved bool
}

var (
	mu       sync.Mutex
	sessions = map[string]*deviceSession{} // device_code → session
	seq      int
)

func main() {
	addr := flag.String("addr", "127.0.0.1:9999", "监听地址")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/oauth/device/code", handleDeviceCode)
	mux.HandleFunc("/verify", handleVerify)
	mux.HandleFunc("/oauth/token", handleToken)

	fmt.Printf("mockauth 监听 http://%s （Ctrl+C 退出）\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

func handleDeviceCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	mu.Lock()
	seq++
	dc := fmt.Sprintf("mock-device-%d-%d", seq, time.Now().UnixNano())
	sessions[dc] = &deviceSession{code: fmt.Sprintf("MOCK-%04d", seq)}
	userCode := sessions[dc].code
	mu.Unlock()

	writeJSON(w, map[string]any{
		"device_code":              dc,
		"user_code":                userCode,
		"verification_uri":         fmt.Sprintf("http://%s/verify?dc=%s", r.Host, dc),
		"verification_uri_complete": fmt.Sprintf("http://%s/verify?dc=%s&uc=%s", r.Host, dc, userCode),
		"expires_in":               600,
		"interval":                 2,
	})
}

func handleVerify(w http.ResponseWriter, r *http.Request) {
	dc := r.URL.Query().Get("dc")
	mu.Lock()
	s, ok := sessions[dc]
	if ok {
		s.approved = true
	}
	mu.Unlock()

	if !ok {
		http.Error(w, "未知 device_code", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "<html><body><h2>授权成功（mock）</h2><p>可关闭此页面，回到 CLI 查看 auth status。</p></body></html>")
}

func handleToken(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	dc := r.FormValue("device_code")

	mu.Lock()
	s, ok := sessions[dc]
	approved := ok && s.approved
	mu.Unlock()

	if !ok {
		writeJSONStatus(w, http.StatusBadRequest, map[string]string{
			"error": "expired_token", "error_description": "未知设备码",
		})
		return
	}
	if !approved {
		writeJSONStatus(w, http.StatusBadRequest, map[string]string{
			"error": "authorization_pending", "error_description": "等待用户授权",
		})
		return
	}

	writeJSON(w, map[string]any{
		"access_token": "mock-token-" + dc,
		"token_type":   "Bearer",
		"expires_in":   3600,
		"scope":        "",
		"account":      "mock-user@hypersku.com",
	})
}

func writeJSON(w http.ResponseWriter, v any) { writeJSONStatus(w, http.StatusOK, v) }

func writeJSONStatus(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
