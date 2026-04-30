package console

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestError(t *testing.T) {
	output := captureOutput(func() {
		Error("test error")
	})
	if !strings.Contains(output, "test error") {
		t.Errorf("Error() output = %q, want contains %q", output, "test error")
	}
}

func TestSuccess(t *testing.T) {
	output := captureOutput(func() {
		Success("test success")
	})
	if !strings.Contains(output, "test success") {
		t.Errorf("Success() output = %q, want contains %q", output, "test success")
	}
}

func TestWarning(t *testing.T) {
	output := captureOutput(func() {
		Warning("test warning")
	})
	if !strings.Contains(output, "test warning") {
		t.Errorf("Warning() output = %q, want contains %q", output, "test warning")
	}
}

func TestInfo(t *testing.T) {
	output := captureOutput(func() {
		Info("test info")
	})
	if !strings.Contains(output, "test info") {
		t.Errorf("Info() output = %q, want contains %q", output, "test info")
	}
}

func TestMuted(t *testing.T) {
	output := captureOutput(func() {
		Muted("test muted")
	})
	if !strings.Contains(output, "test muted") {
		t.Errorf("Muted() output = %q, want contains %q", output, "test muted")
	}
}

func TestSuccessSameLine(t *testing.T) {
	output := captureOutput(func() {
		SuccessSameLine("test same line")
	})
	if !strings.Contains(output, "test same line") {
		t.Errorf("SuccessSameLine() output = %q, want contains %q", output, "test same line")
	}
}

func TestErrorWithArgs(t *testing.T) {
	output := captureOutput(func() {
		Error("test %s %d", "error", 123)
	})
	expected := fmt.Sprintf("test %s %d", "error", 123)
	if !strings.Contains(output, expected) {
		t.Errorf("Error() with args output = %q, want contains %q", output, expected)
	}
}

func TestSuccessWithArgs(t *testing.T) {
	output := captureOutput(func() {
		Success("test %s %d", "success", 456)
	})
	expected := fmt.Sprintf("test %s %d", "success", 456)
	if !strings.Contains(output, expected) {
		t.Errorf("Success() with args output = %q, want contains %q", output, expected)
	}
}

func TestWarningWithArgs(t *testing.T) {
	output := captureOutput(func() {
		Warning("test %s %d", "warning", 789)
	})
	expected := fmt.Sprintf("test %s %d", "warning", 789)
	if !strings.Contains(output, expected) {
		t.Errorf("Warning() with args output = %q, want contains %q", output, expected)
	}
}

func TestInfoWithArgs(t *testing.T) {
	output := captureOutput(func() {
		Info("test %s %d", "info", 101)
	})
	expected := fmt.Sprintf("test %s %d", "info", 101)
	if !strings.Contains(output, expected) {
		t.Errorf("Info() with args output = %q, want contains %q", output, expected)
	}
}

func TestMutedWithArgs(t *testing.T) {
	output := captureOutput(func() {
		Muted("test %s %d", "muted", 202)
	})
	expected := fmt.Sprintf("test %s %d", "muted", 202)
	if !strings.Contains(output, expected) {
		t.Errorf("Muted() with args output = %q, want contains %q", output, expected)
	}
}

func TestSuccessSameLineWithArgs(t *testing.T) {
	output := captureOutput(func() {
		SuccessSameLine("test %s %d", "same line", 303)
	})
	expected := fmt.Sprintf("test %s %d", "same line", 303)
	if !strings.Contains(output, expected) {
		t.Errorf("SuccessSameLine() with args output = %q, want contains %q", output, expected)
	}
}

func TestErrorNewline(t *testing.T) {
	output := captureOutput(func() {
		Error("test")
	})
	if !strings.HasSuffix(output, "\n") {
		t.Errorf("Error() should end with newline, got %q", output)
	}
}

func TestSuccessNewline(t *testing.T) {
	output := captureOutput(func() {
		Success("test")
	})
	if !strings.HasSuffix(output, "\n") {
		t.Errorf("Success() should end with newline, got %q", output)
	}
}

func TestWarningNewline(t *testing.T) {
	output := captureOutput(func() {
		Warning("test")
	})
	if !strings.HasSuffix(output, "\n") {
		t.Errorf("Warning() should end with newline, got %q", output)
	}
}

func TestInfoNewline(t *testing.T) {
	output := captureOutput(func() {
		Info("test")
	})
	if !strings.HasSuffix(output, "\n") {
		t.Errorf("Info() should end with newline, got %q", output)
	}
}

func TestMutedNewline(t *testing.T) {
	output := captureOutput(func() {
		Muted("test")
	})
	if !strings.HasSuffix(output, "\n") {
		t.Errorf("Muted() should end with newline, got %q", output)
	}
}

func TestSuccessSameLineNoNewline(t *testing.T) {
	output := captureOutput(func() {
		SuccessSameLine("test")
	})
	if strings.HasSuffix(output, "\n") {
		t.Errorf("SuccessSameLine() should not end with newline, got %q", output)
	}
}

func TestEmptyString(t *testing.T) {
	output := captureOutput(func() {
		Error("")
	})
	if output == "" {
		t.Error("Error() with empty string should still output styled content")
	}
}

func TestSpecialCharacters(t *testing.T) {
	output := captureOutput(func() {
		Error("special: %s", "&\"<>")
	})
	if !strings.Contains(output, "special:") {
		t.Errorf("Error() with special chars output = %q, want contains %q", output, "special:")
	}
}

func TestMultipleArgs(t *testing.T) {
	output := captureOutput(func() {
		Success("a=%d b=%s c=%v", 1, "two", 3.0)
	})
	if !strings.Contains(output, "a=1") {
		t.Errorf("Success() with multiple args output = %q, want contains %q", output, "a=1")
	}
}
