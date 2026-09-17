package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestCheckError(t *testing.T) {
	// should not panic on nil error
	checkError(nil)

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected checkError to panic on non-nil error")
		}
	}()
	checkError(os.ErrNotExist)
}

func TestCountLinesInFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "words.txt")
	content := "apple\nbanana\ncherry\n"
	if err := os.WriteFile(tmpFile, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	lines, err := countLinesInFile(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lines != 3 {
		t.Fatalf("expected 3 lines, got %d", lines)
	}
}

func TestCountLinesInFileMissing(t *testing.T) {
	_, err := countLinesInFile(filepath.Join(t.TempDir(), "does-not-exist.txt"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestGetSpecialChar(t *testing.T) {
	allowed := map[string]bool{
		"@": true, "#": true, "$": true, "%": true, "^": true,
		"&": true, "+": true, "_": true, "?": true, "~": true,
	}

	for i := 0; i < 50; i++ {
		c := getSpecialChar()
		if !allowed[c] {
			t.Fatalf("unexpected special char %q", c)
		}
	}
}

func TestGetRandomNumber(t *testing.T) {
	for numPwr := 1; numPwr <= 6; numPwr++ {
		n := getRandomNumber(numPwr)
		digits := len(strconv.Itoa(n))
		if digits != numPwr {
			t.Fatalf("numPwr=%d: expected %d digits, got %d (n=%d)", numPwr, numPwr, digits, n)
		}
	}
}

func TestGetRandomWord(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "words.txt")
	content := "apple\nbanana\ncherry\n"
	if err := os.WriteFile(tmpFile, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	valid := map[string]bool{"apple": true, "banana": true, "cherry": true}
	for i := 0; i < 50; i++ {
		w := getRandomWord(tmpFile)
		if !valid[w] {
			t.Fatalf("unexpected word %q", w)
		}
	}
}

func TestGetRandomWordStripsPunctuation(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "words.txt")
	content := "wo-rd!\n"
	if err := os.WriteFile(tmpFile, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	w := getRandomWord(tmpFile)
	if w != "word" {
		t.Fatalf("expected punctuation stripped to \"word\", got %q", w)
	}
}

func TestGetRandomWordMissingFileFallsBackToHex(t *testing.T) {
	w := getRandomWord(filepath.Join(t.TempDir(), "does-not-exist.txt"))
	if w == "nada" || w == "" {
		t.Fatalf("expected a hex fallback word, got %q", w)
	}
	if strings.ContainsAny(w, "ghijklmnopqrstuvwxyz") {
		t.Fatalf("expected hex string, got %q", w)
	}
}
