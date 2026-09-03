package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestStartDeviceFlow(t *testing.T) {
	var gotClientID, gotScope string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("ParseForm() error = %v", err)
		}
		gotClientID = r.FormValue("client_id")
		gotScope = r.FormValue("scope")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"device_code":               "dc-1",
			"user_code":                 "ABCD-EFGH",
			"verification_uri":          "https://auth.example.com/device",
			"verification_uri_complete": "https://auth.example.com/device?user_code=ABCD-EFGH",
			"expires_in":                300,
			"interval":                  5,
		})
	}))
	defer srv.Close()

	c := NewClient(Endpoints{
		DeviceAuthorizationURL: srv.URL + "/device/code",
		TokenURL:               srv.URL + "/token",
		ClientID:               "test-cli",
		Scopes:                 []string{"orders:read"},
	})

	resp, err := c.StartDeviceFlow(context.Background())
	if err != nil {
		t.Fatalf("StartDeviceFlow() error = %v", err)
	}
	if gotClientID != "test-cli" {
		t.Errorf("client_id = %q, want %q", gotClientID, "test-cli")
	}
	if gotScope != "orders:read" {
		t.Errorf("scope = %q, want %q", gotScope, "orders:read")
	}
	if resp.DeviceCode != "dc-1" || resp.UserCode != "ABCD-EFGH" {
		t.Errorf("device flow 响应字段错误: %+v", resp)
	}
	if uri := resp.BestVerificationURI(); uri != "https://auth.example.com/device?user_code=ABCD-EFGH" {
		t.Errorf("BestVerificationURI() = %q", uri)
	}
}

func TestExchangeTokenPending(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error":             "authorization_pending",
			"error_description": "用户尚未完成授权",
		})
	}))
	defer srv.Close()

	c := NewClient(Endpoints{DeviceAuthorizationURL: srv.URL, TokenURL: srv.URL + "/token", ClientID: "t"})
	_, err := c.ExchangeToken(context.Background(), "dc")
	if err == nil {
		t.Fatal("ExchangeToken() 应返回错误")
	}
	te, ok := err.(*TokenError)
	if !ok {
		t.Fatalf("错误类型应为 *TokenError, got %T", err)
	}
	if te.Code != "authorization_pending" {
		t.Errorf("Code = %q, want authorization_pending", te.Code)
	}
}

func TestExchangeTokenSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("ParseForm() error = %v", err)
		}
		if gt := r.FormValue("grant_type"); gt != "urn:ietf:params:oauth:grant-type:device_code" {
			t.Errorf("grant_type = %q", gt)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": "tok-xyz",
			"token_type":   "Bearer",
			"expires_in":   3600,
			"account":      "user@example.com",
		})
	}))
	defer srv.Close()

	c := NewClient(Endpoints{DeviceAuthorizationURL: srv.URL, TokenURL: srv.URL + "/token", ClientID: "t"})
	resp, err := c.ExchangeToken(context.Background(), "dc")
	if err != nil {
		t.Fatalf("ExchangeToken() error = %v", err)
	}

	cred := TokenToCredentials(resp, time.Unix(1000000, 0).UTC())
	if cred.AccessToken != "tok-xyz" {
		t.Errorf("AccessToken = %q", cred.AccessToken)
	}
	if want := time.Unix(1000000+3600, 0).UTC(); !cred.ExpiresAt.Equal(want) {
		t.Errorf("ExpiresAt = %v, want %v", cred.ExpiresAt, want)
	}
	if cred.ExpiresAt.Unix() != 1003600 {
		t.Errorf("ExpiresAt.Unix() = %d, want 1003600", cred.ExpiresAt.Unix())
	}
}

func TestAwaitTokenPollsUntilAuthorized(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		w.Header().Set("Content-Type", "application/json")
		if n < 3 {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "authorization_pending"})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": "tok-final",
			"expires_in":   600,
		})
	}))
	defer srv.Close()

	c := NewClient(Endpoints{DeviceAuthorizationURL: srv.URL, TokenURL: srv.URL + "/token", ClientID: "t"})
	d := &DeviceAuthorizationResponse{DeviceCode: "dc", VerificationURI: "https://x/device", ExpiresIn: 60, Interval: 1}

	cred, err := c.AwaitToken(context.Background(), d, 10*time.Millisecond)
	if err != nil {
		t.Fatalf("AwaitToken() error = %v", err)
	}
	if cred.AccessToken != "tok-final" {
		t.Errorf("AccessToken = %q, want tok-final", cred.AccessToken)
	}
	if atomic.LoadInt32(&calls) != 3 {
		t.Errorf("token 端点调用次数 = %d, want 3", calls)
	}
}

func TestAwaitTokenDenied(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "access_denied"})
	}))
	defer srv.Close()

	c := NewClient(Endpoints{DeviceAuthorizationURL: srv.URL, TokenURL: srv.URL + "/token", ClientID: "t"})
	d := &DeviceAuthorizationResponse{DeviceCode: "dc", VerificationURI: "https://x/device", ExpiresIn: 60}

	if _, err := c.AwaitToken(context.Background(), d, 5*time.Millisecond); err == nil {
		t.Fatal("access_denied 时 AwaitToken() 应返回错误")
	}
}

func TestTryAdvancePendingDevice(t *testing.T) {
	t.Run("pending 中间态保存并推进为已登录", func(t *testing.T) {
		var calls int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if atomic.AddInt32(&calls, 1) < 2 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "authorization_pending"})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "tok-adv", "expires_in": 300})
		}))
		defer srv.Close()

		s := NewStore(t.TempDir())
		c := NewClient(Endpoints{DeviceAuthorizationURL: srv.URL, TokenURL: srv.URL + "/token", ClientID: "t"})

		if err := s.SavePendingDevice(&PendingDevice{
			DeviceCode:      "dc",
			VerificationURI: "https://x/device",
			ExpiresAt:       time.Now().Add(time.Minute),
		}); err != nil {
			t.Fatalf("SavePendingDevice() error = %v", err)
		}

		// 第一次：仍在 pending
		if TryAdvancePendingDevice(context.Background(), c, s, time.Now()) {
			t.Error("第一次推进不应成功")
		}
		if _, err := s.LoadCredentials(); err != nil {
			t.Fatalf("LoadCredentials() error = %v", err)
		}

		// 第二次：服务端发放 token
		if !TryAdvancePendingDevice(context.Background(), c, s, time.Now()) {
			t.Fatal("第二次推进应成功")
		}
		cred, err := s.LoadCredentials()
		if err != nil || cred == nil || cred.AccessToken != "tok-adv" {
			t.Fatalf("推进后凭证错误: (%+v, %v)", cred, err)
		}
		if p, _ := s.LoadPendingDevice(); p != nil {
			t.Error("推进成功后应清理 pending 中间态")
		}
	})

	t.Run("无中间态时为幂等 no-op", func(t *testing.T) {
		s := NewStore(t.TempDir())
		c := NewClient(DefaultEndpoints())
		if TryAdvancePendingDevice(context.Background(), c, s, time.Now()) {
			t.Error("无中间态时不应推进成功")
		}
	})

	t.Run("过期中间态被清理", func(t *testing.T) {
		s := NewStore(t.TempDir())
		c := NewClient(DefaultEndpoints())
		if err := s.SavePendingDevice(&PendingDevice{
			DeviceCode: "dc",
			ExpiresAt:  time.Now().Add(-time.Minute),
		}); err != nil {
			t.Fatalf("SavePendingDevice() error = %v", err)
		}
		if TryAdvancePendingDevice(context.Background(), c, s, time.Now()) {
			t.Error("过期中间态不应推进成功")
		}
		if p, _ := s.LoadPendingDevice(); p != nil {
			t.Error("过期中间态应被清理")
		}
	})
}
