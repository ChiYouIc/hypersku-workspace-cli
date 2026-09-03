// Package auth 提供登录态管理与 OAuth2 Device Code Flow（RFC 8628）客户端。
//
// 凭证保存在用户主目录（~/.hypersku-cli/credentials.json），与 CLI 安装目录分离；
// 设备授权中间态保存在 ~/.hypersku-cli/device-pending.json，
// 供 auth status 在外部轮询场景（如 workbuddy 连接器）下逐步推进换取 token。
package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	credentialsFileName   = "credentials.json"
	pendingDeviceFileName = "device-pending.json"
)

// Store 管理凭证与设备授权中间态的磁盘读写。
type Store struct {
	dir string
}

// NewStore 返回指向指定目录的 Store。
func NewStore(dir string) *Store { return &Store{dir: dir} }

// DefaultStore 返回默认 Store（~/.hypersku-cli）。
func DefaultStore() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("无法获取用户主目录: %w", err)
	}
	return NewStore(filepath.Join(home, ".hypersku-cli")), nil
}

// Dir 返回存储目录路径。
func (s *Store) Dir() string { return s.dir }

func (s *Store) credentialsPath() string { return filepath.Join(s.dir, credentialsFileName) }

func (s *Store) pendingDevicePath() string { return filepath.Join(s.dir, pendingDeviceFileName) }

// Credentials 表示已持久化的登录凭证。
type Credentials struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Scope        string    `json:"scope,omitempty"`
	Account      string    `json:"account,omitempty"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"` // 零值表示未知/不过期
	ObtainedAt   time.Time `json:"obtained_at"`
}

// Valid 判断凭证在给定时刻是否可用（留 30 秒时钟偏差余量）。
func (c *Credentials) Valid(now time.Time) bool {
	if c == nil || c.AccessToken == "" {
		return false
	}
	if c.ExpiresAt.IsZero() {
		return true
	}
	return now.Before(c.ExpiresAt.Add(-30 * time.Second))
}

// PendingDevice 表示一次进行中的设备授权（auth login 发起后暂存，
// 由 auth status 或 login 前台轮询推进）。
type PendingDevice struct {
	DeviceCode          string    `json:"device_code"`
	UserCode            string    `json:"user_code,omitempty"`
	VerificationURI     string    `json:"verification_uri,omitempty"`
	VerificationURIComp string    `json:"verification_uri_complete,omitempty"`
	Interval            int       `json:"interval,omitempty"` // 轮询间隔（秒）
	ExpiresAt           time.Time `json:"expires_at"`
}

// LoadCredentials 读取凭证；文件不存在时返回 (nil, nil)。
func (s *Store) LoadCredentials() (*Credentials, error) {
	return loadJSON[Credentials](s.credentialsPath())
}

// SaveCredentials 写入凭证（权限 0600，仅当前用户可读）。
func (s *Store) SaveCredentials(c *Credentials) error {
	if c == nil {
		return fmt.Errorf("凭证为空，拒绝写入")
	}
	if c.ObtainedAt.IsZero() {
		c.ObtainedAt = time.Now()
	}
	return saveJSON(s.credentialsPath(), c, 0600)
}

// ClearCredentials 删除凭证文件；文件不存在时不视为错误。
func (s *Store) ClearCredentials() error {
	return removeIfExists(s.credentialsPath())
}

// LoadPendingDevice 读取设备授权中间态；文件不存在时返回 (nil, nil)。
func (s *Store) LoadPendingDevice() (*PendingDevice, error) {
	return loadJSON[PendingDevice](s.pendingDevicePath())
}

// SavePendingDevice 写入设备授权中间态。
func (s *Store) SavePendingDevice(p *PendingDevice) error {
	if p == nil {
		return fmt.Errorf("设备授权状态为空，拒绝写入")
	}
	return saveJSON(s.pendingDevicePath(), p, 0600)
}

// ClearPendingDevice 删除设备授权中间态；文件不存在时不视为错误。
func (s *Store) ClearPendingDevice() error {
	return removeIfExists(s.pendingDevicePath())
}

// ---------- 通用 JSON 文件读写 ----------

func loadJSON[T any](path string) (*T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("读取文件失败 %s: %w", path, err)
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("解析文件失败 %s: %w", path, err)
	}
	return &v, nil
}

func saveJSON(path string, v any, perm os.FileMode) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化失败: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("无法创建目录 %s: %w", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, data, perm); err != nil {
		return fmt.Errorf("写入文件失败 %s: %w", path, err)
	}
	return nil
}

func removeIfExists(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除文件失败 %s: %w", path, err)
	}
	return nil
}
