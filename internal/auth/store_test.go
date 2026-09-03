package auth

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCredentialsValid(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name string
		cred *Credentials
		want bool
	}{
		{"nil 凭证", nil, false},
		{"空 token", &Credentials{}, false},
		{"无过期时间", &Credentials{AccessToken: "t"}, true},
		{
			"未过期",
			&Credentials{AccessToken: "t", ExpiresAt: now.Add(time.Hour)},
			true,
		},
		{
			"已过期",
			&Credentials{AccessToken: "t", ExpiresAt: now.Add(-time.Hour)},
			false,
		},
		{
			"临近过期（30 秒余量内）",
			&Credentials{AccessToken: "t", ExpiresAt: now.Add(10 * time.Second)},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cred.Valid(now); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreCredentialsRoundTrip(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir)

	// 未写入时 Load 返回 (nil, nil)
	cred, err := s.LoadCredentials()
	if err != nil {
		t.Fatalf("LoadCredentials() error = %v", err)
	}
	if cred != nil {
		t.Fatalf("LoadCredentials() = %+v, want nil", cred)
	}

	want := &Credentials{
		AccessToken: "tok-123",
		TokenType:   "Bearer",
		Account:     "user@example.com",
		ExpiresAt:   time.Now().Add(time.Hour).UTC().Truncate(time.Second),
	}
	if err := s.SaveCredentials(want); err != nil {
		t.Fatalf("SaveCredentials() error = %v", err)
	}

	got, err := s.LoadCredentials()
	if err != nil {
		t.Fatalf("LoadCredentials() error = %v", err)
	}
	if got.AccessToken != want.AccessToken || got.Account != want.Account {
		t.Errorf("round-trip 不一致: got %+v, want %+v", got, want)
	}
	if !got.ExpiresAt.Equal(want.ExpiresAt) {
		t.Errorf("ExpiresAt 不一致: got %v, want %v", got.ExpiresAt, want.ExpiresAt)
	}

	// 文件权限应为 0600（Windows 上跳过权限断言）
	info, err := os.Stat(s.credentialsPath())
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if runtimeIsPOSIX() && info.Mode().Perm() != 0600 {
		t.Errorf("凭证文件权限 = %v, want 0600", info.Mode().Perm())
	}

	// Clear 后 Load 应返回 nil，再次 Clear 不报错
	if err := s.ClearCredentials(); err != nil {
		t.Fatalf("ClearCredentials() error = %v", err)
	}
	if err := s.ClearCredentials(); err != nil {
		t.Fatalf("重复 ClearCredentials() error = %v", err)
	}
	cred, err = s.LoadCredentials()
	if err != nil || cred != nil {
		t.Fatalf("Clear 后 LoadCredentials() = (%v, %v), want (nil, nil)", cred, err)
	}
}

func TestStorePendingDeviceRoundTrip(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir)

	p, err := s.LoadPendingDevice()
	if err != nil || p != nil {
		t.Fatalf("LoadPendingDevice() = (%v, %v), want (nil, nil)", p, err)
	}

	want := &PendingDevice{
		DeviceCode:      "dc-abc",
		UserCode:        "ABCD-EFGH",
		VerificationURI: "https://auth.hypersku.com/device",
		ExpiresAt:       time.Now().Add(15 * time.Minute).UTC().Truncate(time.Second),
	}
	if err := s.SavePendingDevice(want); err != nil {
		t.Fatalf("SavePendingDevice() error = %v", err)
	}

	got, err := s.LoadPendingDevice()
	if err != nil {
		t.Fatalf("LoadPendingDevice() error = %v", err)
	}
	if got.DeviceCode != want.DeviceCode || got.UserCode != want.UserCode {
		t.Errorf("round-trip 不一致: got %+v, want %+v", got, want)
	}

	if err := s.ClearPendingDevice(); err != nil {
		t.Fatalf("ClearPendingDevice() error = %v", err)
	}
	if err := s.ClearPendingDevice(); err != nil {
		t.Fatalf("重复 ClearPendingDevice() error = %v", err)
	}
}

func TestSaveCredentialsRejectsNil(t *testing.T) {
	s := NewStore(filepath.Join(t.TempDir(), "sub"))
	if err := s.SaveCredentials(nil); err == nil {
		t.Error("SaveCredentials(nil) 应返回错误")
	}
	if err := s.SavePendingDevice(nil); err == nil {
		t.Error("SavePendingDevice(nil) 应返回错误")
	}
}

// runtimeIsPOSIX 报告当前平台是否为类 Unix（用于权限断言）。
func runtimeIsPOSIX() bool {
	return os.PathSeparator == '/' && os.Getenv("OS") != "Windows_NT"
}

// 确保 Credentials / PendingDevice JSON 字段名符合序列化约定。
func TestJSONFieldNames(t *testing.T) {
	c := Credentials{AccessToken: "a"}
	b, _ := json.Marshal(c)
	if string(b) == "" || !json.Valid(b) {
		t.Fatalf("Credentials 序列化结果非法: %s", b)
	}
}
