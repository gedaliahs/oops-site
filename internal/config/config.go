package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	RetentionDays int   `json:"retention_days"`
	MaxTrashBytes int64 `json:"max_trash_bytes"`
	RiskWarning   bool  `json:"risk_warning"`
}

var Default = Config{
	RetentionDays: 7,
	MaxTrashBytes: 5 * 1024 * 1024 * 1024, // 5GB
	RiskWarning:   true,
}

func OopsDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".oops")
}

func TrashDir() string {
	return filepath.Join(OopsDir(), "trash")
}

func JournalPath() string {
	return filepath.Join(OopsDir(), "journal.jsonl")
}

func ConfigPath() string {
	return filepath.Join(OopsDir(), "config.json")
}

func LastCleanupPath() string {
	return filepath.Join(OopsDir(), ".last_cleanup")
}

func Load() Config {
	cfg := Default
	data, err := os.ReadFile(ConfigPath())
	if err != nil {
		return cfg
	}
	_ = json.Unmarshal(data, &cfg)
	if cfg.RetentionDays <= 0 {
		cfg.RetentionDays = Default.RetentionDays
	}
	if cfg.MaxTrashBytes <= 0 {
		cfg.MaxTrashBytes = Default.MaxTrashBytes
	}
	return cfg
}

func Save(cfg Config) error {
	if err := os.MkdirAll(OopsDir(), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigPath(), data, 0o644)
}

func Get(key string) string {
	cfg := Load()
	switch key {
	case "retention_days":
		return formatInt(cfg.RetentionDays)
	case "max_trash_bytes":
		return formatInt64(cfg.MaxTrashBytes)
	case "risk_warning":
		if cfg.RiskWarning {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

func Set(key, value string) error {
	cfg := Load()
	switch key {
	case "retention_days":
		n, err := parseInt(value)
		if err != nil {
			return err
		}
		cfg.RetentionDays = n
	case "max_trash_bytes":
		n, err := parseInt64(value)
		if err != nil {
			return err
		}
		cfg.MaxTrashBytes = n
	case "risk_warning":
		cfg.RiskWarning = value == "true" || value == "1"
	}
	return Save(cfg)
}

func EnsureDir() error {
	if err := os.MkdirAll(OopsDir(), 0o755); err != nil {
		return err
	}
	return os.MkdirAll(TrashDir(), 0o755)
}

func ShouldCleanup() bool {
	data, err := os.ReadFile(LastCleanupPath())
	if err != nil {
		return true
	}
	t, err := time.Parse(time.RFC3339, string(data))
	if err != nil {
		return true
	}
	return time.Since(t) > time.Hour
}

func MarkCleanup() {
	_ = os.WriteFile(LastCleanupPath(), []byte(time.Now().Format(time.RFC3339)), 0o644)
}

func formatInt(n int) string {
	return json.Number(itoa(n)).String()
}

func formatInt64(n int64) string {
	return json.Number(i64toa(n)).String()
}

func parseInt(s string) (int, error) {
	var n int
	err := json.Unmarshal([]byte(s), &n)
	return n, err
}

func parseInt64(s string) (int64, error) {
	var n int64
	err := json.Unmarshal([]byte(s), &n)
	return n, err
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 20)
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	for n > 0 {
		buf = append(buf, byte('0'+n%10))
		n /= 10
	}
	if neg {
		buf = append(buf, '-')
	}
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}

func i64toa(n int64) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 20)
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	for n > 0 {
		buf = append(buf, byte('0'+n%10))
		n /= 10
	}
	if neg {
		buf = append(buf, '-')
	}
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}
