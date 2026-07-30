package version

import (
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	info := Get()

	if info.Version != "0.1.0" {
		t.Errorf("期望 Version = \"0.1.0\"，实际得到 %q", info.Version)
	}

	if info.Commit == "" {
		t.Error("Commit 不应为空")
	}

	if info.Date == "" {
		t.Error("Date 不应为空")
	}

	if info.GoOS == "" {
		t.Error("GoOS 不应为空")
	}

	if info.GoArch == "" {
		t.Error("GoArch 不应为空")
	}
}

func TestInfoString(t *testing.T) {
	info := Info{
		Version: "1.0.0",
		Commit:  "abc123",
		Date:    "2026-07-29",
		GoOS:    "linux",
		GoArch:  "amd64",
	}

	str := info.String()

	expectedParts := []string{"hypersku-cli", "v1.0.0", "abc123", "2026-07-29", "linux", "amd64"}
	for _, part := range expectedParts {
		if !strings.Contains(str, part) {
			t.Errorf("期望字符串包含 %q，实际得到 %q", part, str)
		}
	}
}

func TestInfoStringWithUnknown(t *testing.T) {
	info := Info{
		Version: "0.1.0",
		Commit:  "unknown",
		Date:    "unknown",
		GoOS:    "windows",
		GoArch:  "amd64",
	}

	str := info.String()

	if !strings.Contains(str, "unknown") {
		t.Errorf("期望字符串包含 \"unknown\"，实际得到 %q", str)
	}

	if !strings.Contains(str, "windows/amd64") {
		t.Errorf("期望字符串包含 \"windows/amd64\"，实际得到 %q", str)
	}
}
