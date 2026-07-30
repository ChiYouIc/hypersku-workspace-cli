package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config 表示 HyperSKU CLI 的配置文件结构
type Config struct {
	APIBaseURL string `json:"api_base_url,omitempty"`
	APITimeout int    `json:"api_timeout,omitempty"`
	APIToken   string `json:"api_token,omitempty"`
}

// DefaultConfigDir 返回默认的配置目录: ~/.hypersku-cli
func DefaultConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("无法获取用户主目录: %w", err)
	}
	return filepath.Join(homeDir, ".hypersku-cli"), nil
}

// DefaultConfigPath 返回默认的配置文件路径
func DefaultConfigPath() (string, error) {
	dir, err := DefaultConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// EnsureConfigDir 确保配置目录存在，如果不存在则创建
func EnsureConfigDir() (string, error) {
	dir, err := DefaultConfigDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("无法创建配置目录 %s: %w", dir, err)
	}
	return dir, nil
}

// Load 从指定路径加载配置文件，如果文件不存在则返回默认配置
func Load(path string) (*Config, error) {
	cfg := &Config{
		APITimeout: 30, // 默认超时 30 秒
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// 配置文件不存在，返回默认配置
			return cfg, nil
		}
		return nil, fmt.Errorf("读取配置文件失败 %s: %w", path, err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败 %s: %w", path, err)
	}

	// 设置默认超时
	if cfg.APITimeout <= 0 {
		cfg.APITimeout = 30
	}

	return cfg, nil
}

// Save 将配置保存到指定路径
func Save(path string, cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("无法创建目录 %s: %w", dir, err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败 %s: %w", path, err)
	}

	return nil
}

// LogDir 返回日志目录
func LogDir() (string, error) {
	dir, err := DefaultConfigDir()
	if err != nil {
		return "", err
	}
	logDir := filepath.Join(dir, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return "", fmt.Errorf("无法创建日志目录 %s: %w", logDir, err)
	}
	return logDir, nil
}

// DataDir 返回数据存储目录
func DataDir() (string, error) {
	dir, err := DefaultConfigDir()
	if err != nil {
		return "", err
	}
	dataDir := filepath.Join(dir, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return "", fmt.Errorf("无法创建数据目录 %s: %w", dataDir, err)
	}
	return dataDir, nil
}
