package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setenv temporarily sets an environment variable and restores the original
// value (or unsets it) when the test ends.
func setenv(t *testing.T, key, value string) {
	t.Helper()
	original, existed := os.LookupEnv(key)
	if value == "" {
		os.Unsetenv(key)
	} else {
		os.Setenv(key, value)
	}
	t.Cleanup(func() {
		if existed {
			os.Setenv(key, original)
		} else {
			os.Unsetenv(key)
		}
	})
}

// --- XDG_CONFIG_HOME ---

func TestDir_XDGConfigHome_TakesPriority(t *testing.T) {
	setenv(t, "XDG_CONFIG_HOME", "/custom/config")
	setenv(t, "HOME", "/home/user")

	got := Dir()
	want := "/custom/config/dayplanner"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestDir_XDGConfigHome_EmptyString_Ignored(t *testing.T) {
	setenv(t, "XDG_CONFIG_HOME", "")
	setenv(t, "HOME", "/home/user")

	got := Dir()
	if strings.Contains(got, "XDG") {
		t.Errorf("empty XDG_CONFIG_HOME should be ignored, got %q", got)
	}
	// Should fall through to HOME-based path.
	if !strings.Contains(got, ".config") {
		t.Errorf("expected HOME-based .config path, got %q", got)
	}
}

// --- HOME fallback ---

func TestDir_HomeFallback_UsesConfigSubdir(t *testing.T) {
	setenv(t, "XDG_CONFIG_HOME", "")
	setenv(t, "HOME", "/home/testuser")

	got := Dir()
	want := "/home/testuser/.config/dayplanner"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

// --- App name ---

func TestDir_AlwaysEndsWithAppName(t *testing.T) {
	cases := []struct {
		xdg  string
		home string
	}{
		{"/xdg", "/home/user"},
		{"", "/home/user"},
	}
	for _, c := range cases {
		setenv(t, "XDG_CONFIG_HOME", c.xdg)
		setenv(t, "HOME", c.home)

		got := Dir()
		if filepath.Base(got) != "dayplanner" {
			t.Errorf("expected path to end with 'dayplanner', got %q", got)
		}
	}
}

func TestDir_ReturnsAbsolutePath(t *testing.T) {
	setenv(t, "XDG_CONFIG_HOME", "")
	setenv(t, "HOME", "/home/user")

	got := Dir()
	if !filepath.IsAbs(got) {
		t.Errorf("expected absolute path, got %q", got)
	}
}

func TestDir_XDGPath_IsAbsolute(t *testing.T) {
	setenv(t, "XDG_CONFIG_HOME", "/custom/xdg")
	setenv(t, "HOME", "")

	got := Dir()
	if !filepath.IsAbs(got) {
		t.Errorf("expected absolute path from XDG, got %q", got)
	}
}

